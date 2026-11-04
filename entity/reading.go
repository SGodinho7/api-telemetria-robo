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
