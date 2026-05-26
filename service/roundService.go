package service

import (
	"strconv"

	"api-telemetria-robo/dto"
	"api-telemetria-robo/entity"
	"api-telemetria-robo/logs"
	"api-telemetria-robo/repository"
)

type RoundService struct {
	repo repository.MatchReposiroty
}

func NewRoundService(repo repository.MatchReposiroty) *RoundService {
	return &RoundService{repo: repo}
}

func (r *RoundService) SaveNewRoundReadings(records [][]string) error {
	var (
		sensors    []*entity.Sensor
		sensorsDTO []*dto.SensorDTO
		err        error
	)

	sensors = r.convertRecordsToSensors(records)
	for _, sensor := range sensors {
		sensorData := dto.NewSensorDTO(sensor.GetName(), sensor.GetReadings())
		sensorsDTO = append(sensorsDTO, sensorData)
	}

	if err = r.repo.CreateNewRound(sensorsDTO); err != nil {
		return err
	}

	return nil
}

func (r *RoundService) convertRecordsToSensors(records [][]string) []*entity.Sensor {
	var sensors []*entity.Sensor

	names, rows := records[0], records[1:]

	for it, name := range names {
		if it == 0 {
			continue
		}

		newSensor := entity.NewSensor(name)
		sensors = append(sensors, newSensor)
	}

	for _, row := range rows {
		timestamp, sensorValues := row[0], row[1:]
		for it, value := range sensorValues {
			val, err := strconv.Atoi(value)
			if err != nil {
				logs.Errorf(pkgName, "Could not convert value string into int: %s", err.Error())
				continue
			}
			sensors[it].AppendReading(timestamp, val)
		}
	}

	return sensors
}
