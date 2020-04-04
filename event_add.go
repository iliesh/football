package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TeleBot) eventAdd(u tgbotapi.Update) bool {

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
			msg := `No such user here
Try to run command /start firstly and then try again`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 "+msg))
			return false
		}
		t.HandleError(err)
	}

	if userTypeDB != "organizer" {
		msg := "Forbidden! Only Organizers can create Games!"
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 "+msg+" "))
	} else {
		t.event.eventUserID = userIDDB
		t.event.eventOrganizer = u.Message.From.UserName
	}
	return true
}

func (t *TeleBot) addCalendar(u tgbotapi.Update) {

	var nextMonth time.Month
	var selYear int

	year, month, date := time.Now().Date()
	currDate := date
	td := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

	if month == 12 {
		nextMonth = 1
		selYear = td.Year() + 1
	} else {
		nextMonth = td.Month() + 1
		selYear = year
	}

	varWeekDate := weekDayToInt[time.Now().Weekday().String()]

	if varWeekDate <= 0 || varWeekDate > 7 {
		varWeekDate = 1
	}

	firstDayCurrentMonth := time.Date(td.Year(), td.Month(), 1, 0, 0, 0, 0, time.Local)
	lastDayCurrentMonth := firstDayCurrentMonth.AddDate(0, 1, 0).Add(time.Nanosecond * -1)

	sliceWeek := []string{}
	sliceMonth := [][]string{}

	for i := 0; i < varWeekDate-1; i++ {
		sliceWeek = append(sliceWeek, "0")
	}

	for n := 0; n <= 14; n++ {

		if n >= varWeekDate {
			sliceWeek = append(sliceWeek, strconv.Itoa(date))
			date++
		}

		if date == lastDayCurrentMonth.Day() {
			date = 1
			sliceWeek = append(sliceWeek, strconv.Itoa(lastDayCurrentMonth.Day()))
		}

		if len(sliceWeek) == 7 {
			sliceMonth = append(sliceMonth, sliceWeek)
			sliceWeek = []string{}
		}
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Select date:")
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var rowMonth []tgbotapi.InlineKeyboardButton
	var rowDays []tgbotapi.InlineKeyboardButton
	var rowDatesWeek1 []tgbotapi.InlineKeyboardButton
	var rowDatesWeek2 []tgbotapi.InlineKeyboardButton

	days := []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}

	selMonth := month.String()

	currMonth := tgbotapi.NewInlineKeyboardButtonData(month.String()+" "+strconv.Itoa(year), "eventDate#\u0001"+month.String()+" "+strconv.Itoa(year))

	rowMonth = append(rowMonth, currMonth)

	for _, d := range days {
		btDay := tgbotapi.NewInlineKeyboardButtonData(d, "eventDate#\u0002#"+d+"#"+strconv.Itoa(varWeekDate))
		rowDays = append(rowDays, btDay)
	}

	for i := 0; i < 7; i++ {
		if sliceMonth[0][i] == "0" {
			sliceMonth[0][i] = "\u0000"
			selMonth = "0"
		} else {
			selMonth = month.String()
		}

		fmt.Printf("DEBUG MONTH: %v, Date: %v\n", sliceMonth, currDate)
		if sliceMonth[0][i] > "0" && sliceMonth[0][i] < strconv.Itoa(currDate) {
			selMonth = nextMonth.String()
			fmt.Printf("Next month: %v, date: %v\n", selMonth, currDate)
		}
		fmt.Printf("Select Month: %v\n", selMonth)

		btDate := tgbotapi.NewInlineKeyboardButtonData(sliceMonth[0][i], "eventDate#"+sliceMonth[0][i]+" "+selMonth+" "+strconv.Itoa(selYear))
		rowDatesWeek1 = append(rowDatesWeek1, btDate)
	}

	for i := 0; i < 7; i++ {
		btDate := tgbotapi.NewInlineKeyboardButtonData(sliceMonth[1][i], "eventDate#"+sliceMonth[1][i]+" "+selMonth+" "+strconv.Itoa(selYear))
		rowDatesWeek2 = append(rowDatesWeek2, btDate)
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowMonth)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDays)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDatesWeek1)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rowDatesWeek2)

	msg.ReplyMarkup = keyboard

	msg.Text = "Select Date:"
	t.botAPI.Send(msg)
}

func (t *TeleBot) eventDate(u tgbotapi.Update) {
	resCallBackQuery := strings.SplitN(u.CallbackQuery.Data, "#", -1)
	fmt.Printf("Selected Date: %v\n", resCallBackQuery[1])

	t.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Selected Date: "+resCallBackQuery[1]))
	t.botAPI.Send(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Selected Date: "+resCallBackQuery[1]))

	t.event.eventDate = resCallBackQuery[1]
	t.addTime(u)
}

func (t *TeleBot) addTime(u tgbotapi.Update) {
	del := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
	config := tgbotapi.NewCallback(u.CallbackQuery.ID, strconv.Itoa(u.CallbackQuery.Message.MessageID))
	go t.botAPI.AnswerCallbackQuery(config)
	go t.botAPI.Send(del)
	go t.creatInlineKeyboardTimeMarkup(u.CallbackQuery.Message.Chat.ID)
	go t.emptyAnswer(u.CallbackQuery.ID)
}

func (t *TeleBot) eventTime(u tgbotapi.Update) {
	resCallBackQuery := strings.SplitN(u.CallbackQuery.Data, "#", -1)
	fmt.Printf("Selected Time: %v\n", resCallBackQuery[1])

	t.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Selected Time: "+resCallBackQuery[1]))
	t.botAPI.Send(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Selected Time: "+resCallBackQuery[1]))

	t.event.eventTime = resCallBackQuery[1]
	t.teamSize(u)
}

func (t *TeleBot) teamSize(u tgbotapi.Update) {
	del := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
	config := tgbotapi.NewCallback(u.CallbackQuery.ID, strconv.Itoa(u.CallbackQuery.Message.MessageID))
	go t.botAPI.AnswerCallbackQuery(config)
	go t.botAPI.Send(del)
	go t.creatInlineKeyboardTeamSizeMarkup(u.CallbackQuery.Message.Chat.ID)
	go t.emptyAnswer(u.CallbackQuery.ID)
}

func (t *TeleBot) creatInlineKeyboardTeamSizeMarkup(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Pitch Size?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5x5", "pitchSize#5"),
			tgbotapi.NewInlineKeyboardButtonData("7x7", "pitchSize#7"),
			tgbotapi.NewInlineKeyboardButtonData("9x9", "pitchSize#9"),
			tgbotapi.NewInlineKeyboardButtonData("11x11", "pitchSize#11"),
		),
	)
	t.botAPI.Send(msg)
}

func (t *TeleBot) creatInlineKeyboardTimeMarkup(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Please select Event Time:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("07:00", "eventTime#07:00"),
			tgbotapi.NewInlineKeyboardButtonData("08:00", "eventTime#08:00"),
			tgbotapi.NewInlineKeyboardButtonData("09:00", "eventTime#09:00"),
			tgbotapi.NewInlineKeyboardButtonData("10:00", "eventTime#10:00"),
			tgbotapi.NewInlineKeyboardButtonData("11:00", "eventTime#11:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("12:00", "eventTime#12:00"),
			tgbotapi.NewInlineKeyboardButtonData("13:00", "eventTime#13:00"),
			tgbotapi.NewInlineKeyboardButtonData("14:00", "eventTime#14:00"),
			tgbotapi.NewInlineKeyboardButtonData("15:00", "eventTime#15:00"),
			tgbotapi.NewInlineKeyboardButtonData("16:00", "eventTime#16:00"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("17:00", "eventTime#17:00"),
			tgbotapi.NewInlineKeyboardButtonData("18:00", "eventTime#18:00"),
			tgbotapi.NewInlineKeyboardButtonData("19:00", "eventTime#19:00"),
			tgbotapi.NewInlineKeyboardButtonData("20:00", "eventTime#20:00"),
			tgbotapi.NewInlineKeyboardButtonData("21:00", "eventTime#21:00"),
		),
	)
	t.botAPI.Send(msg)
}

func (t *TeleBot) emptyAnswer(CallbackQueryID string) {
	configAlert := tgbotapi.NewCallback(CallbackQueryID, "")
	t.botAPI.AnswerCallbackQuery(configAlert)
}

func (t *TeleBot) pitchSize(u tgbotapi.Update) {
	del := tgbotapi.NewDeleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
	config := tgbotapi.NewCallback(u.CallbackQuery.ID, strconv.Itoa(u.CallbackQuery.Message.MessageID))
	go t.botAPI.AnswerCallbackQuery(config)
	go t.botAPI.Send(del)

	resCallBackQuery := strings.SplitN(u.CallbackQuery.Data, "#", -1)
	fmt.Printf("Pitch Size: %v\n", resCallBackQuery[1])

	t.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Pitch Size: "+resCallBackQuery[1]+"x"+resCallBackQuery[1]))
	t.botAPI.Send(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, "Pitch Size: "+resCallBackQuery[1]+"x"+resCallBackQuery[1]))
	t.event.eventTeamPlayers = resCallBackQuery[1]
}

func (t *TeleBot) dbUpdateEventAdd(u tgbotapi.Update) {
	fmt.Printf("Updating DB| Date:%s, Organizer: %s, Team Size: %s\n", t.event.eventDate, t.event.eventOrganizer, t.event.eventTeamPlayers)

	fmtDateSlice := strings.SplitN(t.event.eventDate, " ", -1)
	myDateString := fmtDateSlice[2] + "-" + fmtDateSlice[1] + "-" + fmtDateSlice[0] + " " + t.event.eventTime
	myDate, err := time.Parse("2006-January-2 15:04", myDateString)
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}

	db := dbConn("football_bot")
	query := `REPLACE INTO events(bot_id,event_time,event_team_players,event_organizer,event_info,event_status) VALUES (?,?,?,?,?,?);`
	stmt, err := db.Prepare(query)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}

	_, err = stmt.Exec(botID, myDate.Format("2006-01-02 15:04:00"), t.event.eventTeamPlayers, t.event.eventUserID, "No Data", 1)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}

	var dbPlayerID int64
	allPlayersID := make([]int64, 0)

	rows, err := db.Query("select user_id from users where bot_id=?", botID)
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&dbPlayerID)
		if err != nil {
			t.HandleError(err)
			return
		}
		allPlayersID = append(allPlayersID, dbPlayerID)
	}
	err = rows.Err()
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}

	stmt, err = db.Prepare("select event_time from events where event_organizer=?")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.event.eventUserID).Scan(&t.event.eventTime)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := "\u26A0 " + `No Games here
To create a new game run /new_game command`
			t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, msg))
		}
		return
	}

	db.Close()

	td, _ := time.Parse(sqlDateTimeForm, t.event.eventTime)
	fmt.Printf("t Event Time: %v\n", t.event.eventTime)
	for _, playerID := range allPlayersID {
		msg := `New Game has been Scheduled on ` + strconv.Itoa(td.Day()) + ` of ` + td.Month().String() + ` at ` + strconv.Itoa(td.Hour()) + `:00
To add yourself to the Game - type /add
Type /help for more information`
		t.botAPI.Send(tgbotapi.NewMessage(playerID, "\u26BD "+msg))
	}
}
