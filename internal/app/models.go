package app

import (
	"database/sql"
)

type logReq string

type FootballApp struct {
	DB       *sql.DB
	BotToken string
	BotID    string
	Player   playerT
}

type playerT struct {
	ID        int64
	Type      string
	Status    int16
	Username  string
	FirstName string
	LastName  string
	Language  string
}
