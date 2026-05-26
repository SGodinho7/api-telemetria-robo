package dto

import "api-telemetria-robo/entity"

type SensorDTO struct {
	name     string
	readings []entity.Reading
}

func NewSensorDTO(name string, readings []entity.Reading) *SensorDTO {
	return &SensorDTO{
		name:     name,
		readings: readings,
	}
}

func (s *SensorDTO) GetName() string {
	return s.name
}

func (s *SensorDTO) GetReadings() []entity.Reading {
	return s.readings
}
