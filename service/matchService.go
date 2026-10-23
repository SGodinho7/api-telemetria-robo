package service

import (
	"errors"
	"time"

	"api-telemetria-robo/dto"
	"api-telemetria-robo/logs"
	"api-telemetria-robo/repository"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var pkgName logs.PackageName = "Service"

var (
	errMatchOpen       = errors.New("a match is currently open")
	errNoMatchOpen     = errors.New("no match is currently open")
	errNegativeMatchID = errors.New("match id cannot be negative")
)

type MatchService struct {
	repo repository.MatchReposiroty
}

// NewMatchService creates a new MatchService instance and returns its pointer.
func NewMatchService(repo repository.MatchReposiroty) *MatchService {
	return &MatchService{
		repo: repo,
	}
}

// OpenNewMatch creates a new match, checks if there isn't an open match in the repository and saves it.
func (m *MatchService) OpenNewMatch(title string, date time.Time, opponentName string) error {
	var (
		openMatch *dto.MatchDTO
		err       error
	)

	openMatch, err = m.GetCurrentMatch()
	if err != nil && err != errNoMatchOpen {
		return err
	}
	if openMatch != nil {
		return errMatchOpen
	}

	err = m.repo.CreateMatch(title, date, opponentName)
	if err != nil {
		return err
	}

	return nil
}

// GetMatchByID retrives a match from the repository, identifing it by its ID.
func (m *MatchService) GetMatchByID(matchID int) (*dto.MatchDTO, error) {
	var (
		match *dto.MatchDTO
		err   error
	)

	if matchID < 0 {
		return nil, errNegativeMatchID
	}

	if match, err = m.repo.FindMatchByID(matchID); err != nil {
		return nil, err
	}

	return match, nil
}

// GetCurrentMatch retrives the currently open match from the repository.
// If there's none, it returns an error.
func (m *MatchService) GetCurrentMatch() (*dto.MatchDTO, error) {
	var (
		curMatch *dto.MatchDTO
		err      error
	)

	curMatch, err = m.repo.FindOpenMatch()
	if err == mongo.ErrNoDocuments {
		return nil, errNoMatchOpen
	} else if err != nil {
		return nil, err
	}

	return curMatch, nil
}

// CLoseCurrentMatch closes the currently open match in the repository.
// If there's none, it returns an error.
func (m *MatchService) CloseCurrentMatch() error {
	var (
		curMatch *dto.MatchDTO
		err      error
	)

	curMatch, err = m.GetCurrentMatch()
	if err != nil {
		return err
	}
	if curMatch.GetID() == 0 {
		return errNoMatchOpen
	}

	err = m.repo.CloseMatch(curMatch.GetID())
	if err != nil {
		return err
	}

	return nil
}
