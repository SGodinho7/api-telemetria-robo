package dto

import (
	"time"
)

type ClosedMatchDTO struct {
	id           int
	title        string
	date         time.Time
	opponentName string
	closed       bool
	roundCount   int
}

func NewClosedMatchDTO(id int, title string, date time.Time, opponentName string, closed bool, roundCount int) *ClosedMatchDTO {
	return &ClosedMatchDTO{
		id:           id,
		title:        title,
		date:         date,
		opponentName: opponentName,
		closed:       closed,
		roundCount:   roundCount,
	}
}

func (m *ClosedMatchDTO) GetID() int {
	return m.id
}

func (m *ClosedMatchDTO) GetTitle() string {
	return m.title
}

func (m *ClosedMatchDTO) GetDate() time.Time {
	return m.date
}

func (m *ClosedMatchDTO) GetOpponentName() string {
	return m.opponentName
}

func (m *ClosedMatchDTO) IsClosed() bool {
	return m.closed
}

func (m *ClosedMatchDTO) GetRoundCount() int {
	return m.roundCount
}
