package app

import (
	"context"
	"database/sql"
	"errors"

	log "github.com/iliesh/go-templates/logger"
)

func (a *FootballApp) cmdStart(ctx context.Context) error {
	log.DebugX(ctx.Value(ctxKey).(string), "Start Message, Check if user id <%d> is not already registered", a.Player.ID)

	var uid sql.NullInt64
	var ustatus sql.NullString
	var utype sql.NullInt16

	sqlQuery := "select user_id, user_status, user_type from telegram_users where user_id=? and bot_id=?"
	log.TraceX("Executing SQL Query: <%s>, with Value: <%d>", sqlQuery, a.Player.ID)
	err := a.DB.QueryRow(sqlQuery, a.Player.ID, a.BotID).Scan(&uid, &ustatus, &utype)

	if err != nil && err != sql.ErrNoRows {
		log.ErrorX(ctx.Value(ctxKey).(string), "DB Query Error: <%v>, SQL Querry: <%s>, UserID: <%d>, BotID: <%s>", err, sqlQuery, a.Player.ID, a.BotID)
		return err
	}

	if err == sql.ErrNoRows {
		log.DebugX(ctx.Value(ctxKey).(string), "There is no user with such ID <%d> in our DB", a.Player.ID)

		sqlQuery := "INSERT INTO telegram_users(user_id,bot_id,first_name,last_name,user_name,user_lang,user_type,user_status) VALUES (?,?,?,?,?,?,?,?)"
		stmt, err := a.DB.Prepare(sqlQuery)
		if err != nil {
			// t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
			return err
		}

		_, err = stmt.Exec(a.Player.ID, a.BotID, a.Player.FirstName, a.Player.LastName, a.Player.Username, a.Player.Language, 0, 1)
		if err != nil {
			// t.botAPI.Send(tgbotapi.NewMessage(u.Message.Chat.ID, "\u26A0 System Error "+t.HandleError(err)))
			return err
		}

		return nil
	}

	log.DebugX(ctx.Value(ctxKey).(string), "There is already an user with such ID <%d> in our DB, User Type: <%d>, Status: <%s>", a.Player.ID, utype.Int16, ustatus.String)
	return errors.New("user already exists")
}

func (a *FootballApp) cmdHelp(ctx context.Context) error {
	log.DebugX(ctx.Value(ctxKey).(string), "Sending Help Message")

	// m := sendMessageReqBodyT{}
	// m.ChatID = a.TGBot.WebHookReq.Message.Chat.ID
	// m.ParseMode = "html"

	// m.Text = playerHelp

	// if a.UserType == "admin" {
	// 	m.Text = adminHelp
	// }

	// req, err := json.Marshal(m)
	// if err != nil {
	// 	log.ErrorX(ctx.Value(ctxKey).(string), "Unable to Encode the Request <%v>, Error: %s", m, err.Error())
	// 	return err
	// }

	// // Send a post request with your token
	// resp, err := http.Post("https://api.telegram.org/"+a.BotToken+"/sendmessage", "application/json", bytes.NewBuffer(req))
	// if err != nil {
	// 	log.ErrorX(ctx.Value(ctxKey).(string), "HTTP Request Error, Error: %s", m, err.Error())
	// 	return err
	// }
	// log.DebugX(ctx.Value(ctxKey).(string), "HTTP Response: %v", resp)

	// if resp.StatusCode != 200 {
	// 	log.ErrorX(ctx.Value(ctxKey).(string), "Bad HTTP Response Code: %d", resp.StatusCode)
	// 	return errors.New("bad http status code")
	// }

	return nil
}
