package telegram

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/iliesh/go-templates/logger"
)

type logReq string

type BotAPI struct {
	DB       *sql.DB
	WebHook  *webHookReqBodyT
	Token    string
	UserType string
}

const (
	ctxKey logReq = "logid"
)

const (
	adminHelp = `
===== G A M E =====
<i>- Create a new Game</i>
<code>/new_game</code>
<i>- Cancel current Game</i>
<code>/stop_game</code>

===== P L A Y E R =====
<i>- Adding yourself to the Game</i>
<code>/add</code>
<i>- Bring one more player with You</i>
<code>/add_other</code>
<i>- Removing yourself from the Game</i>
<code>/cancel</code>
<i>- Removing someone you brought to the Game</i>
<code>/cancel_other</code>
<i>- Show Game Status</i>
<code>/status</code>

===== A P P =====
<i>- Show Help Message:</i>
<code>/help</code>
`

	playerHelp = `
=====  P L A Y E R =====
<i>- Adding yourself to the Game</i>
<code>/add</code>
<i>- Bring one more player with You</i>
<code>/add_other</code>
<i>- Removing yourself from the Game</i>
<code>/cancel</code>
<i>- Removing someone you brought to the Game</i>
<code>/cancel_other</code>
<i>- Show Game Status</i>
<code>/status</code>

===== A P P =====
<i>- Show Help Message:</i>
<code>/help</code>
`
)

// HandlerBot is called everytime telegram sends us a webhook event on the specific path
func HandlerBot(res http.ResponseWriter, req *http.Request, bot *BotAPI) {

	// Generating new request id every time this function is called
	log.ReqID = log.RandomString(8)
	ctx := context.Background()
	logReq := context.WithValue(ctx, ctxKey, log.ReqID)

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.ErrorX(log.ReqID, "Error reading request body: %s", err.Error())
		http.Error(res, err.Error(), 500)
		return
	}

	log.DebugX(log.ReqID, "Got Request: %s", string(b))

	// First, decode the JSON response body
	err = json.Unmarshal(b, &bot.WebHook)
	if err != nil {
		log.ErrorX(log.ReqID, "Error decoding request body, error: %s", err.Error())
		http.Error(res, err.Error(), 500)
		return
	}

	log.TraceX(log.ReqID, "Request has been successfully decoded in Request ID: <%d>, Message: <%v>, CallBack Query: <%v>", bot.WebHook.UpdateID, bot.WebHook.Message, bot.WebHook.CallBackQuery)

	if !bot.ReqInit(logReq) {
		return
	}

	log.DebugX(log.ReqID, "Processing message ID: %d", bot.WebHook.Message.MessageID)
	log.DebugX(log.ReqID, "Processing Text: %s", bot.WebHook.Message.Text)

	if len(bot.WebHook.Message.Entities) > 0 {
		log.TraceX(log.ReqID, "Got <%d> Entities", len(bot.WebHook.Message.Entities))
		if bot.WebHook.Message.Entities[0].Type == "bot_command" {
			log.DebugX(log.ReqID, "Processing Bot Commands")
			if err := bot.BotCommands(logReq); err != nil {
				return
			}
		}
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
}

func (b *BotAPI) ReqInit(ctx context.Context) bool {
	log.DebugX(ctx.Value(ctxKey).(string), "Initial Checks")

	log.TraceX(ctx.Value(ctxKey).(string), "Checking Message ID Value")
	if b.WebHook.Message.MessageID == 0 {
		log.ErrorX(ctx.Value(ctxKey).(string), "Unable to continue without Message ID Value")
		return false
	}

	log.TraceX(ctx.Value(ctxKey).(string), "Check if Message came from another bot")
	if b.WebHook.Message.From.IsBot {
		log.ErrorX(ctx.Value(ctxKey).(string), "Bots are not allowed here")
		return false
	}

	// Check User Permission
	log.TraceX(ctx.Value(ctxKey).(string), "Check User Type from our Database")
	sqlQuery := "select user_type from telegram_users where id=?"
	log.TraceX("Executing SQL Query: <%s>, with Value: <%d>", sqlQuery, b.WebHook.Message.From.ID)
	err := b.DB.QueryRow(sqlQuery, b.WebHook.Message.From.ID).Scan(&b.UserType)

	if err != nil && err != sql.ErrNoRows {
		log.ErrorX(ctx.Value(ctxKey).(string), "DB Query Error: <%v>", err)
		return false
	}

	if err == sql.ErrNoRows {
		log.ErrorX(ctx.Value(ctxKey).(string), "There is no user with such ID <%d> in our DB", b.WebHook.Message.From.ID)
		return false
	}

	return true
}

func (b *BotAPI) BotCommands(ctx context.Context) error {
	log.TraceX(ctx.Value(ctxKey).(string), "Processing Bot Command")
	switch b.WebHook.Message.Text {
	case "/help":
		log.InfoX(ctx.Value(ctxKey).(string), "Processing <%s> command", b.WebHook.Message.Text)
		if err := b.cmdHelp(ctx); err != nil {
			return errors.New("error processing help command")
		}
		return nil
	case "/new_game":
		log.InfoX(ctx.Value(ctxKey).(string), "Processing <%s> command", b.WebHook.Message.Text)
		return nil
	default:
		log.ErrorX(ctx.Value(ctxKey).(string), "Unknown command: <%s>", b.WebHook.Message.Text)
		return errors.New("unknown Command")
	}
}

func (b *BotAPI) cmdHelp(ctx context.Context) error {
	log.TraceX(ctx.Value(ctxKey).(string), "Sending Help Message")

	m := sendMessageReqBodyT{}
	m.ChatID = b.WebHook.Message.Chat.ID
	m.ParseMode = "html"
	m.Text = playerHelp

	req, err := json.Marshal(m)
	if err != nil {
		log.ErrorX(ctx.Value(ctxKey).(string), "Unable to Encode the Request <%v>, Error: %s", m, err.Error())
		return err
	}

	// Send a post request with your token
	resp, err := http.Post("https://api.telegram.org/"+b.Token+"/sendmessage", "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.ErrorX(ctx.Value(ctxKey).(string), "HTTP Request Error, Error: %s", m, err.Error())
		return err
	}
	log.DebugX(ctx.Value(ctxKey).(string), "HTTP Response: %v", resp)

	if resp.StatusCode != 200 {
		log.ErrorX(ctx.Value(ctxKey).(string), "Bad HTTP Response Code: %d", resp.StatusCode)
		return errors.New("bad http status code")
	}

	return nil
}
