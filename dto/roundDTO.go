package dto

type RoundDTO struct {
	roundNumber int
	sensors     []SensorDTO
}

func NewRoundDTO(roundNumber int, sensors []SensorDTO) *RoundDTO {
	return &RoundDTO{
		roundNumber: roundNumber,
		sensors:     sensors,
	}
}

func (r *RoundDTO) GetNumber() int {
	return r.roundNumber
}

func (r *RoundDTO) GetSensors() []SensorDTO {
	return r.sensors
}
