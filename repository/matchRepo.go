package repository

import (
	"time"

	"api-telemetria-robo/dto"
)

type MatchReposiroty interface {
	CreateMatch(title string, date time.Time, opponentName string) error
	FindMatchByID(matchID int) (*dto.MatchDTO, error)
	GetOpenMatch() (*dto.MatchDTO, error)
	CloseMatch(matchID int) error
	CreateNewRound(sensors []*dto.SensorDTO) error
}
