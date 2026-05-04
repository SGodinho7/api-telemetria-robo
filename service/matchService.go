package service

import (
	"errors"
	"time"

	"api-telemetria-robo/dto"
	"api-telemetria-robo/repository"
)

var (
	errMatchOpen  = errors.New("a match is currently open")
	errNoCurMatch = errors.New("no match is currently open")
)

type MatchService struct {
	repo repository.MatchReposiroty
}

func NewMatchService(repo repository.MatchReposiroty) *MatchService {
	return &MatchService{repo: repo}
}

func (m *MatchService) OpenNewMatch(title string, date time.Time, opponentName string) error {
	var (
		openMatch dto.MatchDTO
		err       error
	)

	openMatch, err = m.repo.GetOpenMatch()
	if err != nil {
		return err
	}
	if openMatch.GetID() != 0 {
		return errMatchOpen
	}

	err = m.repo.CreateMatch(title, date, opponentName)
	if err != nil {
		return err
	}

	return nil
}

func (m *MatchService) GetCurrentMatch() (dto.MatchDTO, error) {
	var (
		curMatch dto.MatchDTO
		err      error
	)

	curMatch, err = m.repo.GetOpenMatch()
	if err != nil {
		return dto.MatchDTO{}, err
	}
	if curMatch.GetID() == 0 {
		return dto.MatchDTO{}, errNoCurMatch
	}

	return curMatch, nil
}

func (m *MatchService) CloseCurrentMatch() error {
	var (
		curMatch dto.MatchDTO
		err      error
	)

	curMatch, err = m.repo.GetOpenMatch()
	if err != nil {
		return err
	}
	if curMatch.GetID() == 0 {
		return errNoCurMatch
	}

	err = m.repo.CloseMatch(curMatch.GetID())
	if err != nil {
		return err
	}

	return nil
}
