package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TeleBot) start(u tgbotapi.Update) bool {

	db := dbConn("football_bot")
	userType := "player"

	var countDBPlayer int
	stmt, err := db.Prepare("select count(user_id) from users where bot_id=?")
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	defer stmt.Close()
	err = stmt.QueryRow(botID, u.Message.Chat.ID).Scan(&countDBPlayer)
	if err != nil {
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		t.HandleError(err)
	}

	if countDBPlayer == 0 {
		userType = "organizer"
	}

	query := `INSERT INTO users(bot_id,user_id,first_name,last_name,user_name,user_lang,user_type) VALUES (?,?,?,?,?,?,?)
	ON DUPLICATE KEY UPDATE first_name='` + u.Message.From.FirstName + `',last_name='` + u.Message.From.LastName + `',
	user_lang='` + u.Message.From.LanguageCode + `'`
	stmt, err = db.Prepare(query)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}

	_, err = stmt.Exec(botID, u.Message.Chat.ID, u.Message.From.FirstName, u.Message.From.LastName, u.Message.From.UserName, u.Message.From.LanguageCode, userType)
	if err != nil {
		db.Close()
		t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
		return false
	}
	db.Close()

	htmlT := `<b></b><strong>Hi ` + u.Message.From.FirstName + ` ` + u.Message.From.LastName + `!</strong>
<em>We are delighted to have you among us!</em>
` + htmlHelpText

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, htmlT)
	msg.ParseMode = tgbotapi.ModeHTML

	t.botAPI.Send(msg)
	return true
}
