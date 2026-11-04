package controller

type matchJson struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Date         string `json:"date"`
	OpponentName string `json:"opponentName"`
	Closed       bool   `json:"closed"`
}

type matchClosedJson struct {
	matchJson
	RoundCount int `json:"roundCount"`
}

type matchCreateJson struct {
	Title        string `json:"title"`
	Date         string `json:"date"`
	OpponentName string `json:"opponentName"`
}

type roundJson struct {
	RoundNumber int          `json:"roundNumeber"`
	Sensors     []sensorJson `json:"sensors"`
}

type sensorJson struct {
	Name     string        `json:"name"`
	Readings []readingJson `json:"readings"`
}

type readingJson struct {
	Timestamp string `json:"timestamp"`
	Value     int    `json:"value"`
}

type roundCreateJson struct {
	Sensors string `json:"sensors"`
}
