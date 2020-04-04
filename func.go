package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TeleBot) dbPlayers(u tgbotapi.Update) (map[int64][]string, error) {
	var dbPlayerID int64
	var dbPlayerFirstName string
	var dbPlayerLastName string
	var dbPlayerType string
	var dbPlayerReserveNr string
	var eventID string
	var eventTime string
	var eventTeamPlayers int
	dbAllPlayers := make(map[int64][]string)
	db := dbConn("football_bot")

	stmt, err := db.Prepare("SELECT id, event_time, event_team_players FROM events WHERE bot_id=? AND event_status=1")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(botID).Scan(&eventID, &eventTime, &eventTeamPlayers)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := `There are no Games currently created
Try to run command /new_game firstly and then try again`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 "+msg))
			return nil, err
		}
		t.HandleError(err)
	}

	rows, err := db.Query("SELECT players.`user_id`, first_name, last_name, user_type, players.`reserve_nr` FROM users RIGHT JOIN players ON users.`user_id` = players.`user_id` WHERE users.`bot_id`=? AND players.event_id=? ORDER by players.reserve_nr, players.id", botID, eventID)
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		dbPlayer := []string{""}
		err := rows.Scan(&dbPlayerID, &dbPlayerFirstName, &dbPlayerLastName, &dbPlayerType, &dbPlayerReserveNr)
		if err != nil {
			t.HandleError(err)
			return nil, err
		}

		dbPlayer = append(dbPlayer, dbPlayerFirstName, dbPlayerLastName, dbPlayerType, dbPlayerReserveNr)
		dbAllPlayers[dbPlayerID] = dbPlayer
	}
	err = rows.Err()
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return nil, err
	}

	// fmt.Printf("All Players Slice: %v\n", dbAllPlayers)
	return dbAllPlayers, nil
}

func (t *TeleBot) bcastMessage(m string, u tgbotapi.Update) error {
	allPlayersID, err := t.dbPlayers(u)

	if err != nil {
		return err
	}

	for k := range allPlayersID {
		if k != u.Message.Chat.ID && k > 1000 {
			fmt.Printf("Sending Message to player: %d, %s\n", k, m)
			t.botAPI.Send(tgbotapi.NewMessage(k, "\u26BD "+m))
		}
	}
	return nil
}

func createTeam(allPlayersIDs []int64, teamPlayers int) (err error) {

	teamA, teamB = nil, nil

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allPlayersIDs), func(i, j int) { allPlayersIDs[i], allPlayersIDs[j] = allPlayersIDs[j], allPlayersIDs[i] })

	teamA = allPlayersIDs[:teamPlayers]
	teamB = allPlayersIDs[teamPlayers : teamPlayers*2]
	return nil
}
