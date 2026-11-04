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
	mux.HandleFunc("/get-closed-match-id/{id:[0-9]+}", m.getClosedMatchByID).Methods("GET")
	mux.HandleFunc("/get-closed-matches", m.getClosedMatches).Methods("GET")
	mux.HandleFunc("/get-current-match", m.getCurrentMatch).Methods("GET")
	mux.HandleFunc("/close-current-match", m.closeCurrentMatch).Methods("POST")
	mux.HandleFunc("/post-round", m.postRound).Methods("POST")
	mux.HandleFunc("/get-match-round/{id:[0-9]+}/{roundNumber:[0-9]+}", m.getMatchRound).Methods("GET")
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

func (m *MatchController) getClosedMatchByID(w http.ResponseWriter, r *http.Request) {
	var (
		matchID, _     = strconv.Atoi(mux.Vars(r)["id"])
		matchFront     matchClosedJson
		closedMatchDTO *dto.ClosedMatchDTO
		err            error
	)

	if closedMatchDTO, err = m.matchService.GetCLosedMatchByID(matchID); err != nil {
		logs.Errorf(pkgName, "Could not find closed match with id %d: %s", matchID, err)
		serveError(w, err.Error())
		return
	}

	matchFront = matchClosedJson{
		matchJson: matchJson{
			ID:           closedMatchDTO.GetID(),
			Title:        closedMatchDTO.GetTitle(),
			Date:         closedMatchDTO.GetDate().Format("2006-01-02"),
			OpponentName: closedMatchDTO.GetOpponentName(),
			Closed:       closedMatchDTO.IsClosed(),
		},
		RoundCount: closedMatchDTO.GetRoundCount(),
	}

	serveJson(w, matchFront)
}

func (m *MatchController) getClosedMatches(w http.ResponseWriter, r *http.Request) {
	var (
		matches          []*dto.ClosedMatchDTO
		matchesJsonArray []matchClosedJson
		err              error
	)

	if matches, err = m.matchService.GetAllClosedMatchesFront(); err != nil {
		serveError(w, err.Error())
		return
	}

	for _, match := range matches {
		matchFront := matchClosedJson{
			matchJson: matchJson{
				ID:           match.GetID(),
				Title:        match.GetTitle(),
				Date:         match.GetDate().Format("2006-01-02"),
				OpponentName: match.GetOpponentName(),
				Closed:       match.IsClosed(),
			},
			RoundCount: match.GetRoundCount(),
		}
		matchesJsonArray = append(matchesJsonArray, matchFront)
	}

	serveJson(w, matchesJsonArray)
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

	if err = m.roundService.SaveRoundCSVRecords(records); err != nil {
		logs.Errorf(pkgName, "Could not save round: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	serveSuccess(w)
}

func (m *MatchController) getMatchRound(w http.ResponseWriter, r *http.Request) {
	var (
		matchID, _     = strconv.Atoi(mux.Vars(r)["id"])
		roundNumber, _ = strconv.Atoi(mux.Vars(r)["roundNumber"])
		roundDTO       *dto.RoundDTO
		roundJSON      roundJson
		err            error
		sensors        = []sensorJson{}
		readings       = []readingJson{}
	)

	if roundDTO, err = m.matchService.GetMatchRound(matchID, roundNumber); err != nil {
		logs.Errorf(pkgName, "Could not get round: %s", err.Error())
		serveError(w, err.Error())
		return
	}

	for _, sensor := range roundDTO.GetSensors() {
		readings = []readingJson{}
		for _, reading := range sensor.GetReadings() {
			readings = append(readings, readingJson{
				Timestamp: reading.GetTimestamp(),
				Value:     reading.GetValue(),
			})
		}
		sensors = append(sensors, sensorJson{
			Name:     sensor.GetName(),
			Readings: readings,
		})
	}
	roundJSON.RoundNumber = roundDTO.GetNumber()
	roundJSON.Sensors = sensors

	serveJson(w, roundJSON)
}
