package main

import (
	"database/sql"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TeleBot) gameStatus(u tgbotapi.Update) error {

	db := dbConn("football_bot")
	gameDetails := ""
	f := ""

	stmt, err := db.Prepare("select event_organizer, event_time, event_team_players from events where bot_id=? and event_status=1")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(botID).Scan(&t.event.eventUserID, &t.event.eventTime, &t.event.eventTeamPlayers)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "\u26A0 " + `No Games here
To create a new game run /new_game command`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg))
		}
		t.HandleError(err)
		return err
	}

	db.Close()

	td, _ := time.Parse(sqlDateTimeForm, t.event.eventTime)
	allPlayers, err := t.dbPlayers(u)
	eventTeamPlayers, _ := strconv.Atoi(t.event.eventTeamPlayers)

	if len(allPlayers) >= eventTeamPlayers*2 && len(teamA) == 0 {
		allPlayersIDs := make([]int64, 0)

		for k, v := range allPlayers {
			if v[4] != "0" {
				continue
			}
			allPlayersIDs = append(allPlayersIDs, k)
		}

		err = createTeam(allPlayersIDs, eventTeamPlayers)
		if err != nil {
			return err
		}
	}

	switch {
	case len(allPlayers) < eventTeamPlayers*2:

		f = ""
		i := 0
		for _, v := range allPlayers {
			i++
			f += strconv.Itoa(i) + ". " + v[1] + " " + v[2] + "\n"
		}

		gameDetails = `Game on: ` + td.Weekday().String() + ", " + strconv.Itoa(td.Day()) + ` of ` + td.Month().String() + ` at ` + strconv.Itoa(td.Hour()) + `:00
Current Players: `

	case len(allPlayers) == eventTeamPlayers*2:

		f = "<u>Team A:</u>\n"

		for i, s := range teamA {
			i++
			f += strconv.Itoa(i) + ". " + allPlayers[s][1] + " " + allPlayers[s][2] + "\n"
		}

		f += "\n<u>Team B:</u>\n"
		for i, s := range teamB {
			i++
			f += strconv.Itoa(i) + ". " + allPlayers[s][1] + " " + allPlayers[s][2] + "\n"
		}

		gameDetails = `Teams for: ` + td.Weekday().String() + ", " + strconv.Itoa(td.Day()) + ` of ` + td.Month().String() + ` at ` + strconv.Itoa(td.Hour()) + `:00`

	case len(allPlayers) > eventTeamPlayers*2:

		f = "<u>Team A:</u>\n"
		for i, s := range teamA {
			i++
			f += strconv.Itoa(i) + ". " + allPlayers[s][1] + " " + allPlayers[s][2] + "\n"
		}

		f += "\n<u>Team B:</u>\n"
		for i, s := range teamB {
			i++
			f += strconv.Itoa(i) + ". " + allPlayers[s][1] + " " + allPlayers[s][2] + "\n"
		}

		for _, v := range allPlayers {
			if v[4] != "0" {
				f += "\n" + v[1] + " " + v[2] + " (Reserve: " + v[4] + ")"
				continue
			}
		}

		gameDetails = `Teams for: ` + td.Weekday().String() + ", " + strconv.Itoa(td.Day()) + ` of ` + td.Month().String() + ` at ` + strconv.Itoa(td.Hour()) + `:00`

	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "")
	msg.Text = gameDetails + "\n" + f + "\n<i>Team Size: " + t.event.eventTeamPlayers + "x" + t.event.eventTeamPlayers + "</i>"
	msg.ParseMode = tgbotapi.ModeHTML
	t.botAPI.Send(msg)
	return nil
}
