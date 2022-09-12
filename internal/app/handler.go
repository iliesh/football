package app

import (
	"context"
	"database/sql"

	log "github.com/iliesh/go-templates/logger"
)

const (
	ctxKey logReq = "logid"
)

func (a *FootballApp) checkUserPerms(ctx context.Context) error {
	// Check User Permission
	log.TraceX(ctx.Value(ctxKey).(string), "Check User Type from our Database, user ID: <%d>", a.Player.ID)

	var p sql.NullString
	sqlQuery := "select user_type from telegram_users where id=?"
	log.TraceX("Executing SQL Query: <%s>, with Value: <%d>", sqlQuery, a.Player.ID)
	err := a.DB.QueryRow(sqlQuery, a.Player.ID).Scan(&p)

	if err != nil && err != sql.ErrNoRows {
		log.ErrorX(ctx.Value(ctxKey).(string), "DB Query Error: <%v>", err)
		return err
	}

	if err == sql.ErrNoRows {
		log.ErrorX(ctx.Value(ctxKey).(string), "There is no user with such ID <%d> in our DB", a.Player.ID)
		return err
	}

	a.Player.Type = p.String
	return nil
}

// if len(body.Message.Entities) > 0 && body.Message.Entities[0].Type == "bot_command" {
// 	switch body.Message.Text {
// 	case "/help":
// 		log.Info("show help text")
// 		showHelp(body.Message.From.ID)
// 		return
// 	case "/new_game":
// 		log.Info("creating a new game")
// 		log.Debug("Check if User: %s %s (%s) is allowed to create a new game", body.Message.From.FirstName, body.Message.From.LastName, body.Message.From.Username)
// 		userPerm, err := newGamePermission(body.Message.From.ID)
// 		if err != nil {
// 			log.Error("Error: %s", err.Error())
// 			m := sendMessageReqBodyT{ChatID: body.Message.From.ID, Text: "\u26A0 Internal Error"}
// 			err = sendMessage(m)
// 			if err != nil {
// 				log.Error("Error: %s", err.Error())
// 			}
// 			return
// 		}
// 		if !userPerm {
// 			log.Warning("User ID: <%d> is not allowed to create a new game", body.Message.From.ID)
// 			m := sendMessageReqBodyT{ChatID: body.Message.From.ID, Text: "\u26A0 Sorry, only admins can create or cancel games"}
// 			err = sendMessage(m)
// 			if err != nil {
// 				log.Error("Error: %s", err.Error())
// 			}
// 			return
// 		}
// 		log.Debug("Select the Game Date")
// 		gameDate, err := selectDate(body.Message.From.ID)
// 		if err != nil {
// 			log.Error("Error: %s", err.Error())
// 			return
// 		}
// 		log.Info("Selected Game Date: %s", gameDate)
// 		return
// 	default:
// 		log.Warning("unknown command")
// 		showHelp(body.Message.From.ID)
// 		return
// 	}
// }

// if body.CallBackQuery.ID != "" {
// 	log.Info("Processing Call Back Query ID: %s", body.CallBackQuery.ID)

// 	// Declining Messages and Commands from another bots
// 	if body.CallBackQuery.From.IsBot {
// 		log.Warning("bots are not accepted here")
// 		return
// 	}
// 	log.Debug("selected date: %s", body.CallBackQuery.Data)

// 	if body.CallBackQuery.Data == "month" {
// 		log.Warning("unable to select month as a value")
// 		m := answerCallbackQueryT{CallBackQueryID: body.CallBackQuery.ID, Text: "Cannot select month, please select the date", ShowAlert: true}
// 		err = answerCallBack(m)
// 		if err != nil {
// 			log.Error("Error: %s", err.Error())
// 		}
// 	}
// 	m := answerCallbackQueryT{CallBackQueryID: body.CallBackQuery.ID, Text: "Selected Date: " + body.CallBackQuery.Data, ShowAlert: false}
// 	err = answerCallBack(m)
// 	if err != nil {
// 		log.Error("Error: %s", err.Error())
// 	}

// 	d := editMessageTextT{MessageID: body.CallBackQuery.Message.MessageID, Text: "date: 01.01.01", ReplyMarkup: inlineKeyboardMarkupT{}}
// 	err = editMessage(d)
// 	if err != nil {
// 		log.Error("Error: %s", err.Error())
// 	}
// }

// log.Warning("unknown command")
