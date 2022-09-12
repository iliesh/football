package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	app "github.com/iliesh/football/internal/app"
	log "github.com/iliesh/go-templates/logger"
)

const (
	ctxKey logReq = "logid"
)

// HandlerBot is called everytime telegram sends us a webhook event on the specific path
func HandlerBot(res http.ResponseWriter, req *http.Request, app *app.FootballApp) {
	// Generating new request id every time this function is called
	log.ReqID = log.RandomString(8)
	ctx := context.Background()
	logReq := context.WithValue(ctx, ctxKey, log.ReqID)

	log.DebugX(log.ReqID, "New HTTP Request: %s", req.Method)

	if req.Method != "POST" {
		log.ErrorX(log.ReqID, "<Method: <%s> not supported")
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.ErrorX(log.ReqID, "Error reading request body: %s", err.Error())
		http.Error(res, err.Error(), 500)
		return
	}

	log.DebugX(log.ReqID, "Got Request: %s", string(b))

	w := webHookReqBodyT{}

	// First, decode the JSON response body
	err = json.Unmarshal(b, &w)
	if err != nil {
		log.ErrorX(log.ReqID, "Error decoding request body, error: %s", err.Error())
		http.Error(res, err.Error(), 500)
		return
	}
	log.DebugX(log.ReqID, "Request has been successfully decoded in Update ID: <%d>, Message ID: <%d>, CallBack Query: <%v>", w.UpdateID, w.Message.MessageID, w.CallBackQuery)

	tgBot := BotAPI{
		WebHookReq: w,
	}

	if !tgBot.reqInit(logReq) {
		log.ErrorX(log.ReqID, "Request didn't pass Initial Checks")
		return
	}

	log.DebugX(log.ReqID, "Processing message ID: <%d>, text: <%s>", w.Message.MessageID, w.Message.Text)
	if err := tgBot.botCommands(logReq, tgBot.WebHookReq.Message.Text); err != nil {
		return
	}
}

func (bot *BotAPI) reqInit(ctx context.Context) bool {
	log.DebugX(ctx.Value(ctxKey).(string), "Initial Checks")

	log.DebugX(ctx.Value(ctxKey).(string), "Checking Message ID Value: <%d>", bot.WebHookReq.Message.MessageID)
	if bot.WebHookReq.Message.MessageID == 0 {
		log.ErrorX(ctx.Value(ctxKey).(string), "Unable to continue without Message ID Value")
		return false
	}

	log.DebugX(ctx.Value(ctxKey).(string), "Check if Message came from another bot")
	if bot.WebHookReq.Message.From.IsBot {
		log.ErrorX(ctx.Value(ctxKey).(string), "Bots are not allowed here")
		return false
	}

	log.DebugX(ctx.Value(ctxKey).(string), "Check if this is a bot command")
	if len(bot.WebHookReq.Message.Entities) == 0 {
		log.ErrorX(ctx.Value(ctxKey).(string), "Missing Entities, Ignoring Request")
		return false
	}

	log.DebugX(log.ReqID, "Got <%d> Entities, First Type: <%s>", len(bot.WebHookReq.Message.Entities), bot.WebHookReq.Message.Entities[0].Type)
	if bot.WebHookReq.Message.Entities[0].Type == "bot_command" {
		return true
	}

	return true
}

func (bot *BotAPI) botCommands(ctx context.Context, cmd string) error {
	log.TraceX(ctx.Value(ctxKey).(string), "Processing Bot Command")
	switch cmd {
	case "/start":
		log.InfoX(ctx.Value(ctxKey).(string), "Processing <%s> command", cmd)

		// if err := a.cmdStart(ctx); err != nil {
		// 	return errors.New("error processing help command")
		// }
		return nil
	case "/help":
		log.InfoX(ctx.Value(ctxKey).(string), "Processing <%s> command", cmd)
		// if err := a.checkUserPerms(ctx); err != nil {
		// 	log.ErrorX(ctx.Value(ctxKey).(string), "Unable to get user permissions, error: <%v>", err)
		// 	return err
		// }

		// if err := a.cmdHelp(ctx); err != nil {
		// 	return errors.New("error processing help command")
		// }
		return nil
	case "/new_game":
		log.InfoX(ctx.Value(ctxKey).(string), "Processing <%s> command", cmd)
		// if err := a.checkUserPerms(ctx); err != nil {
		// 	log.ErrorX(ctx.Value(ctxKey).(string), "Unable to get user permissions, error: <%v>", err)
		// 	return err
		// }

		return nil
	default:
		log.ErrorX(ctx.Value(ctxKey).(string), "Unknown command: <%s>", cmd)
		return errors.New("unknown Command")
	}
}

// func sendMessage() error {
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
// 	return nil
// }
