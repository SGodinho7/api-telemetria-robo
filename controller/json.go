package controller

type matchJson struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Date         string `json:"date"`
	OpponentName string `json:"opponentName"`
	Closed       bool   `json:"closed"`
}

type matchCreateJson struct {
	Title        string `json:"title"`
	Date         string `json:"date"`
	OpponentName string `json:"opponentName"`
}

type roundCreateJson struct {
	Sensors string `json:"sensors"`
}
