package entity

type Reading struct {
	Timestamp string
	Value     int
}

type Sensor struct {
	name     string
	readings []Reading
}

func NewSensor(name string) *Sensor {
	return &Sensor{
		name: name,
	}
}

func (s *Sensor) GetName() string {
	return s.name
}

func (s *Sensor) GetReadings() []Reading {
	return s.readings
}

func (s *Sensor) AppendReading(timestamp string, value int) {
	s.readings = append(s.readings, Reading{
		Timestamp: timestamp,
		Value:     value,
	})
}
