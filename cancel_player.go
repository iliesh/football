package main

import (
	"database/sql"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TeleBot) cancelPlayer(u tgbotapi.Update) bool {

	var eventTeamPlayers int

	eventID := ""
	allPlayers, err := t.dbPlayers(u)
	db := dbConn("football_bot")

	stmt, err := db.Prepare("SELECT id, event_team_players FROM events WHERE bot_id=? AND event_status=1")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(botID).Scan(&eventID, &eventTeamPlayers)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := `There are no Games currently created
Try to run command /new_game firstly and then try again`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 "+msg))
			return false
		}
		t.HandleError(err)
	}

	var countDBPlayer int
	stmt, err = db.Prepare("select count(*) from players where event_id=? and user_id=?")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(eventID, u.Message.Chat.ID).Scan(&countDBPlayer)
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		t.HandleError(err)
	}

	fmt.Printf("Count DB player: %d\n", countDBPlayer)
	if countDBPlayer < 1 {
		msg := "You are not registered to this game"
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg))
		return true
	}

	query := `DELETE FROM players where user_id=? and event_id=?;`
	stmt, err = db.Prepare(query)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}

	_, err = stmt.Exec(u.Message.Chat.ID, eventID)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}

	if allPlayers[u.Message.Chat.ID][4] == "0" {
		teamA = nil
		teamB = nil
		query = `UPDATE players SET reserve_nr=reserve_nr-1 WHERE reserve_nr > 0 ORDER BY id`
	} else {
		query = `UPDATE players SET reserve_nr=reserve_nr-1 WHERE reserve_nr > 1 ORDER BY id`
	}

	stmt, err = db.Prepare(query)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	_, err = stmt.Exec()
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}

	msg := "You have been removed from the upcoming game"
	t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26BD "+msg))

	db.Close()

	err = t.bcastMessage(`Player `+u.Message.From.FirstName+" "+u.Message.From.LastName+` has left the game!
Current Players: `+strconv.Itoa(len(allPlayers)-1)+`/`+strconv.Itoa(eventTeamPlayers*2), u)

	if err != nil {
		fmt.Printf("Error: %v", err)
		return false
	}

	return true
}
