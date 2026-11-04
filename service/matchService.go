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
	errMatchOpen           = errors.New("a match is currently open")
	errNoMatchOpen         = errors.New("no match is currently open")
	errMatchNotClosed      = errors.New("the match is not closed")
	errNegativeMatchID     = errors.New("match id cannot be negative")
	errNegativeRoundNumber = errors.New("round number cannot be negative")
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
func (m *MatchService) GetCLosedMatchByID(matchID int) (*dto.ClosedMatchDTO, error) {
	var (
		match       *dto.MatchDTO
		closedMatch *dto.ClosedMatchDTO
		err         error
	)

	if matchID < 0 {
		return nil, errNegativeMatchID
	}

	if match, err = m.repo.FindMatchByID(matchID); err != nil {
		return nil, err
	}

	if !match.IsClosed() {
		return nil, errMatchNotClosed
	}

	var roundCount int
	if roundCount, err = m.repo.FindMatchRoundCount(matchID); err != nil {
		return nil, err
	}

	closedMatch = dto.NewClosedMatchDTO(
		match.GetID(),
		match.GetTitle(),
		match.GetDate(),
		match.GetOpponentName(),
		match.IsClosed(),
		roundCount,
	)

	return closedMatch, nil
}

func (m *MatchService) GetAllClosedMatchesFront() ([]*dto.ClosedMatchDTO, error) {
	var (
		matches      []*dto.MatchDTO
		matchesFront []*dto.ClosedMatchDTO
		err          error
	)

	if matches, err = m.repo.FindAllClosedMatches(); err != nil {
		return nil, err
	}

	var roundCount int
	for _, match := range matches {
		if roundCount, err = m.repo.FindMatchRoundCount(match.GetID()); err != nil {
			return nil, err
		}

		matchFront := dto.NewClosedMatchDTO(
			match.GetID(),
			match.GetTitle(),
			match.GetDate(),
			match.GetOpponentName(),
			match.IsClosed(),
			roundCount,
		)
		matchesFront = append(matchesFront, matchFront)
	}

	return matchesFront, nil
}

func (m *MatchService) GetMatchRound(matchID, roundNumber int) (*dto.RoundDTO, error) {
	var (
		round *dto.RoundDTO
		err   error
	)

	if matchID < 0 {
		return nil, errNegativeMatchID
	}
	if roundNumber < 0 {
		return nil, errNegativeRoundNumber
	}

	if round, err = m.repo.FindMatchRound(matchID, roundNumber); err != nil {
		return nil, err
	}

	return round, nil
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
