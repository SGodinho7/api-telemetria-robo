package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"api-telemetria-robo/dto"
	"api-telemetria-robo/logs"
	"api-telemetria-robo/service"
	"api-telemetria-robo/util"

	"github.com/gorilla/mux"
)

var pkgName logs.PackageName = "Controller"

type MatchController struct {
	matchService *service.MatchService
	roundService *service.RoundService
}

func NewMatchController(matchServ *service.MatchService, roundServ *service.RoundService) *MatchController {
	return &MatchController{
		matchService: matchServ,
		roundService: roundServ,
	}
}

func (m *MatchController) LoadRoutes(mux *mux.Router) {
	mux.HandleFunc("/post-match", m.postNewMatch).Methods("POST")
	mux.HandleFunc("/get-match-id/{id:[0-9]+}", m.getMatchByID).Methods("GET")
	mux.HandleFunc("/get-current-match", m.getCurrentMatch).Methods("GET")
	mux.HandleFunc("/close-current-match", m.closeCurrentMatch).Methods("POST")
	mux.HandleFunc("/post-round", m.postRound).Methods("POST")
}

func (m *MatchController) postNewMatch(w http.ResponseWriter, r *http.Request) {
	var (
		newMatch matchCreateJson
		err      error
	)

	if err = json.NewDecoder(r.Body).Decode(&newMatch); err != nil {
		logs.Errorf(pkgName, "Could not decode JSON data: %s", err.Error())
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

	serveSuccess(w)
}

func (m *MatchController) getMatchByID(w http.ResponseWriter, r *http.Request) {
	var (
		matchID, _ = strconv.Atoi(mux.Vars(r)["id"])
		match      matchJson
		matchDTO   *dto.MatchDTO
		err        error
	)

	if matchDTO, err = m.matchService.GetMatchByID(matchID); err != nil {
		logs.Errorf(pkgName, "Could not find match with id %d: %s", matchID, err)
		serveError(w, err.Error())
		return
	}

	match = matchJson{
		ID:           matchDTO.GetID(),
		Title:        matchDTO.GetTitle(),
		Date:         matchDTO.GetDate().Format("2006-01-02"),
		OpponentName: matchDTO.GetOpponentName(),
		Closed:       matchDTO.IsClosed(),
	}

	serveJson(w, match)
}

func (m *MatchController) getCurrentMatch(w http.ResponseWriter, r *http.Request) {
	var (
		match    matchJson
		curMatch *dto.MatchDTO
		err      error
	)

	curMatch, err = m.matchService.GetCurrentMatch()
	if err != nil {
		logs.Errorf(pkgName, "Could not get current match: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	match = matchJson{
		ID:           curMatch.GetID(),
		Title:        curMatch.GetTitle(),
		Date:         curMatch.GetDate().Format("2006-01-02"),
		OpponentName: curMatch.GetOpponentName(),
		Closed:       curMatch.IsClosed(),
	}

	serveJson(w, match)
}

func (m *MatchController) closeCurrentMatch(w http.ResponseWriter, r *http.Request) {
	var err error

	err = m.matchService.CloseCurrentMatch()
	if err != nil {
		logs.Errorf(pkgName, "Could not close current match: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	serveSuccess(w)
}

func (m *MatchController) postRound(w http.ResponseWriter, r *http.Request) {
	var (
		roundJson roundCreateJson
		err       error
	)

	if err = json.NewDecoder(r.Body).Decode(&roundJson); err != nil {
		logs.Errorf(pkgName, "Could not decode JSON data: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	records, err := util.DecodeBase64CSV(roundJson.Sensors)
	if err != nil {
		logs.Errorf(pkgName, "Cound not decode Base64 to CSV records: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	if err = m.roundService.SaveRoundRecords(records); err != nil {
		logs.Errorf(pkgName, "Could not save round: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	serveSuccess(w)
}
