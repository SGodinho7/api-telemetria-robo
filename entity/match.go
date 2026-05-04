package entity

import "time"

type Match struct {
	id           int
	title        string
	date         time.Time
	opponentName string
	closed       bool
}

func NewMatch(id int, title string, date time.Time, opponentName string) *Match {
	return &Match{
		id:           id,
		title:        title,
		date:         date,
		opponentName: opponentName,
	}
}

func (m *Match) GetID() int {
	return m.id
}

func (m *Match) GetTitle() string {
	return m.title
}

func (m *Match) UpdateTitle(title string) {
	m.title = title
}

func (m *Match) GetDate() time.Time {
	return m.date
}

func (m *Match) GetOpponentName() string {
	return m.opponentName
}

func (m *Match) UpdateOpponentName(opponentName string) {
	m.opponentName = opponentName
}

func (m *Match) IsClosed() bool {
	return m.closed
}

func (m *Match) CloseMatch() {
	m.closed = true
}
