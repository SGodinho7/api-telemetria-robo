package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"api-telemetria-robo/dto"
	"api-telemetria-robo/logs"
	"api-telemetria-robo/service"
)

var pkgName logs.PackageName = "Controller"

type MatchController struct {
	matchService *service.MatchService
}

func NewMatchController(service *service.MatchService) *MatchController {
	return &MatchController{matchService: service}
}

func (m *MatchController) LoadRoutes(smux *http.ServeMux) {
	smux.HandleFunc("POST /post-match", m.postNewMatch)
	smux.HandleFunc("GET /get-current-match", m.getCurrentMatch)
}

func (m *MatchController) postNewMatch(w http.ResponseWriter, r *http.Request) {
	var (
		newMatch matchCreateJSON
		err      error
	)

	if err = json.NewDecoder(r.Body).Decode(&newMatch); err != nil {
		logs.Errorf(pkgName, "Could not decode JSON: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	dateLayout := "2006-01-02"
	date, err := time.Parse(dateLayout, newMatch.Date)
	if err != nil {
		logs.Errorf(pkgName, "Could not parse date: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	err = m.matchService.OpenNewMatch(newMatch.Title, date, newMatch.OpponentName)
	if err != nil {
		logs.Errorf(pkgName, "Could not open new match: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	serveJSON(w, `{"status": "success"}`)
}

func (m *MatchController) getCurrentMatch(w http.ResponseWriter, r *http.Request) {
	var (
		match    matchJSON
		curMatch dto.MatchDTO
		err      error
	)

	curMatch, err = m.matchService.GetCurrentMatch()
	if err != nil {
		logs.Errorf(pkgName, "Could not get current match: %s", err.Error())
		serveError(w, err.Error())
	}

	match = matchJSON{
		ID:           curMatch.GetID(),
		Title:        curMatch.GetTitle(),
		Date:         curMatch.GetDate().Format("2006-01-02"),
		OpponentName: curMatch.GetOpponentName(),
		Closed:       curMatch.IsClosed(),
	}

	serveJSON(w, match)
}
