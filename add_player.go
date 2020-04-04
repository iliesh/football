package main

import (
	"database/sql"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TeleBot) addPlayer(u tgbotapi.Update) bool {

	var eventID string
	var eventTime string
	var eventTeamPlayers int
	var userIDDB string
	var userTypeDB string
	var reserveNr int

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
			msg := `No such user here
Try to run command /start firstly and then try again`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 "+msg))
			return false
		}
		t.HandleError(err)
	}

	stmt, err = db.Prepare("SELECT id, event_time, event_team_players FROM events WHERE bot_id=? AND event_status=1")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(botID).Scan(&eventID, &eventTime, &eventTeamPlayers)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := `There are no Games currently created
Try to run command /new_game firstly and then try again`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 "+msg))
			return false
		}
		t.HandleError(err)
	}

	// fmt.Printf("Event Team Players: %d\n", eventTeamPlayers)

	var countDBPlayer int
	stmt, err = db.Prepare("SELECT (SELECT count(*) FROM players WHERE event_id=? and user_id=?) as count, (SELECT COALESCE(max(reserve_nr), '0') FROM players WHERE event_id=?) as reserve_nr")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(eventID, u.Message.Chat.ID, eventID).Scan(&countDBPlayer, &reserveNr)
	if err != nil {
		if err != sql.ErrNoRows {
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
			return false
		}
		t.HandleError(err)
	}

	// fmt.Printf("dbPlayerID: %d, ChatID: %d\n", countDBPlayer, u.Message.Chat.ID)

	if countDBPlayer > 0 {
		msg := "You are already subscribed to this game"
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg))
		return true
	}

	allPlayers, err := t.dbPlayers(u)
	allPlayersIDs := make([]int64, 0)

	if len(allPlayers) >= eventTeamPlayers*2 {
		reserveNr++

		msg := `The Team is full already!
At the moment there are ` + strconv.Itoa(len(allPlayers)) + ` players available to play
Putting you in Reserve, your Reserve Number: ` + strconv.Itoa(reserveNr)
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u261D "+msg))
	}

	query := `REPLACE INTO players (user_id,event_id,reserve_nr) VALUES (?,?,?);`
	stmt, err = db.Prepare(query)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}

	_, err = stmt.Exec(u.Message.Chat.ID, eventID, reserveNr)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}

	msg := "Success!"
	t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26BD "+msg))

	db.Close()

	switch {
	case len(allPlayers)+1 < eventTeamPlayers*2 || len(allPlayers)+1 > eventTeamPlayers*2:
		err := t.bcastMessage(`Player `+u.Message.From.FirstName+" "+u.Message.From.LastName+` Joined in!
Current Players: `+strconv.Itoa(len(allPlayers)+1)+`/`+strconv.Itoa(eventTeamPlayers*2), u)

		// fmt.Printf("Len Map %d\n", len(allPlayers))

		if err != nil {
			return false
		}
		return true
	case len(allPlayers)+1 == eventTeamPlayers*2:
		allPlayers, err = t.dbPlayers(u)
		if err != nil {
			return false
		}
		for k := range allPlayers {
			allPlayersIDs = append(allPlayersIDs, k)
		}

		err = createTeam(allPlayersIDs, eventTeamPlayers)
		if err != nil {
			return false
		}

		msg = `Team Complete! Game is on!`

		for _, playerID := range allPlayersIDs {
			if playerID != u.Message.Chat.ID && playerID > 1000 {
				t.botAPI.Send(tgbotapi.NewMessage(playerID, "\u26BD "+msg))
			}
		}

		t.gameStatus(u)

		// fmt.Printf("Count Players: %d, Event Team Players: %d, All Players ID: %v\n", len(allPlayers), eventTeamPlayers, allPlayersIDs)

		return true
	}

	return true
}
