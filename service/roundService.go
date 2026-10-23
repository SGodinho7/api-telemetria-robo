package service

import (
	"strconv"

	"api-telemetria-robo/dto"
	"api-telemetria-robo/entity"
	"api-telemetria-robo/logs"
	"api-telemetria-robo/repository"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RoundService struct {
	repo repository.MatchReposiroty
}

// NewRoundService creates a new RoundService instance and returns its pointer.
func NewRoundService(repo repository.MatchReposiroty) *RoundService {
	return &RoundService{repo: repo}
}

// SaveRoundCSVRecords saves round CSV records data,
// converting the sensor raw data into objects and sending to its repository.
func (r *RoundService) SaveRoundCSVRecords(records [][]string) error {
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

	err = r.repo.CreateNewRound(sensorsDTO)
	if err == mongo.ErrNoDocuments {
		return errNoMatchOpen
	} else if err != nil {
		return err
	}

	logs.Info(pkgName, "Saved new round successfully!")
	return nil
}

// convertRecordsToSensors convert CSV sensor readings into Sensor objects array.
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
			sensors[it].AppendReading(entity.NewReading(timestamp, val))
		}
	}

	return sensors
}
