package controller

type matchJSON struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Date         string `json:"date"`
	OpponentName string `json:"opponentName"`
	Closed       bool   `json:"closed"`
}

type matchCreateJSON struct {
	Title        string `json:"title"`
	Date         string `json:"date"`
	OpponentName string `json:"opponentName"`
}
