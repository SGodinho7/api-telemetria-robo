package repository

import (
	"time"

	"api-telemetria-robo/dto"
)

type MatchReposiroty interface {
	CreateMatch(title string, date time.Time, opponentName string) error
	FindMatchByID(matchID int) (*dto.MatchDTO, error)
	FindMatchRound(matchID, roundNumber int) (*dto.RoundDTO, error)
	FindMatchRoundCount(matchID int) (int, error)
	FindAllClosedMatches() ([]*dto.MatchDTO, error)
	FindOpenMatch() (*dto.MatchDTO, error)
	CloseMatch(matchID int) error
	CreateNewRound(sensors []*dto.SensorDTO) error
}
