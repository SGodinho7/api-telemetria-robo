package dto

import (
	"time"

	"api-telemetria-robo/entity"
)

type MatchDTO struct {
	id           int
	title        string
	date         time.Time
	opponentName string
	closed       bool
}

func NewMatchDTO(id int, title string, date time.Time, opponentName string) *MatchDTO {
	return &MatchDTO{
		id:           id,
		title:        title,
		date:         date,
		opponentName: opponentName,
	}
}

func (m *MatchDTO) GetID() int {
	return m.id
}

func (m *MatchDTO) GetTitle() string {
	return m.title
}

func (m *MatchDTO) GetDate() time.Time {
	return m.date
}

func (m *MatchDTO) GetOpponentName() string {
	return m.opponentName
}

func (m *MatchDTO) IsClosed() bool {
	return !m.closed
}

func (m *MatchDTO) FromEntity(match *entity.Match) {
	m.id = match.GetID()
	m.title = match.GetTitle()
	m.date = match.GetDate()
	m.opponentName = match.GetOpponentName()
	m.closed = match.IsClosed()
}
