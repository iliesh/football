package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TeleBot) eventDel(u tgbotapi.Update) bool {

	var userIDDB string
	var userTypeDB string

	db := dbConn("football_bot")

	stmt, err := db.Prepare("select id, user_type from users where user_name=?")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(u.Message.From.UserName).Scan(&userIDDB, &userTypeDB)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "\u26A0 " + `No such user here
Try to run the command /start firstly and then try again `
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg+t.HandleError(err)))
		}
		return false
	}

	if userTypeDB != "organizer" {
		msg := "\u26A0 Forbidden! Only Organizers can Cancel Games! "
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg+t.HandleError(err)))
		return false
	}

	t.event.eventUserID = userIDDB

	stmt, err = db.Prepare("select event_time from events where event_organizer=?")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.event.eventUserID).Scan(&t.event.eventTime)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "\u26A0 " + `No Games here
To create a new game run /new_game command`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg))
		}
		return false
	}

	t.dbUpdateEventDel(u)
	return true
}

func (t *TeleBot) dbUpdateEventDel(u tgbotapi.Update) {
	fmt.Printf("Updating DB| User ID :%s\n", t.event.eventUserID)
	allPlayers, err := t.dbPlayers(u)
	allPlayersIDs := make([]int64, 0)
	if err != nil {
		return
	}

	db := dbConn("football_bot")
	query := `DELETE FROM events where event_organizer=? and bot_id=?;`
	stmt, err := db.Prepare(query)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}

	fmt.Printf("DB Querry: DELETE FROM events where event_organizer=%v;", t.event.eventUserID)
	_, err = stmt.Exec(t.event.eventUserID, botID)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}
	db.Close()

	td, _ := time.Parse(sqlDateTimeForm, t.event.eventTime)
	t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "Game on "+strconv.Itoa(td.Day())+" of "+td.Month().String()+" at "+strconv.Itoa(td.Hour())+":00 has been cancelled"))

	for k := range allPlayers {
		allPlayersIDs = append(allPlayersIDs, k)
	}
	msg := `Game on ` + td.Weekday().String() + ", " + strconv.Itoa(td.Day()) + ` of ` + td.Month().String() + ` at ` + strconv.Itoa(td.Hour()) + `:00
Has been Cancelled`

	for _, playerID := range allPlayersIDs {
		if playerID != u.Message.Chat.ID && playerID > 1000 {
			t.botAPI.Send(tgbotapi.NewMessage(playerID, "\u26BD "+msg))
		}
	}
}
