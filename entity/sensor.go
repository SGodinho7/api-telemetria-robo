package entity

type Reading struct {
	timestamp string
	value     int
}

func NewReading(timestamp string, value int) Reading {
	return Reading{
		timestamp: timestamp,
		value:     value,
	}
}

func (r *Reading) GetTimestamp() string {
	return r.timestamp
}

func (r *Reading) GetValue() int {
	return r.value
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
		timestamp: timestamp,
		value:     value,
	})
}
