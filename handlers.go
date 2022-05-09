package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/iliesh/go-templates/logger"
)

// HandlerRoot is called everytime telegram sends us a webhook event
func HandlerRoot(res http.ResponseWriter, req *http.Request) {
	log.Info("root request: %v, %s", req.URL, req.RemoteAddr)
}

// HandlerBot is called everytime telegram sends us a webhook event
func HandlerBot(res http.ResponseWriter, req *http.Request) {

	// Generating new request id every time this function is called
	log.ReqID = log.RandomString(8)
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Warning("Unable to get request body: %s", err.Error())
		http.Error(res, err.Error(), 500)
		return
	}

	log.Debug("Got Request: %s", b)

	// First, decode the JSON response body
	body := &webHookReqBodyT{}
	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Error("could not decode request body, error: %s", err.Error())
		return
	}

	// fmt.Printf("%#v\n", body)
	log.Debug("updateid: %d", body.UpdateID)

	if body.Message.MessageID != 0 {
		log.Info("Processing message ID: %d", body.Message.MessageID)

		// Declining Messages and Commands from another bots
		if body.Message.From.IsBot {
			log.Warning("bots are not accepted here")
			return
		}

		log.Info("Processing Text: %s", body.Message.Text)

		if len(body.Message.Entities) > 0 && body.Message.Entities[0].Type == "bot_command" {
			switch body.Message.Text {
			case "/help":
				log.Info("show help text")
				showHelp(body.Message.From.ID)
				return
			case "/new_game":
				log.Info("creating a new game")
				err = newGame(body)
				if err != nil {
					log.Error("Unable to create a new game, error: %s", err.Error())
				}
				return
			default:
				log.Warning("unknown command")
				showHelp(body.Message.From.ID)
				return
			}
		}
	}

	if body.CallBackQuery.ID != "" {
		log.Info("Processing Call Back Query ID: %s", body.CallBackQuery.ID)

		// Declining Messages and Commands from another bots
		if body.CallBackQuery.From.IsBot {
			log.Warning("bots are not accepted here")
			return
		}
		log.Info("selected date: %s", body.CallBackQuery.Data)

		if body.CallBackQuery.Data == "month" {
			log.Warning("unable to select month as a value")
			m := answerCallbackQueryT{CallBackQueryID: body.CallBackQuery.ID, Text: "Cannot select month, please select the date", ShowAlert: true}
			err = answerCallBack(m)
			if err != nil {
				log.Error("Error: %s", err.Error())
			}
			return
		}

		if body.CallBackQuery.Data == "Invalid Date" {
			log.Warning("cannot use the date in the past")
			m := answerCallbackQueryT{CallBackQueryID: body.CallBackQuery.ID, Text: "Cannot select the Day in the past, please select another date", ShowAlert: true}
			err = answerCallBack(m)
			if err != nil {
				log.Error("Error: %s", err.Error())
			}
			return
		}

		m := answerCallbackQueryT{CallBackQueryID: body.CallBackQuery.ID, Text: "Selected Date: " + body.CallBackQuery.Data, ShowAlert: false}
		err = answerCallBack(m)
		if err != nil {
			log.Error("Error: %s", err.Error())
			return
		}

		d := editMessageTextT{ChatID: body.CallBackQuery.Message.Chat.ID, MessageID: body.CallBackQuery.Message.MessageID, Text: "Selected Date: " + body.CallBackQuery.Data}
		err = editMessage(d)
		if err != nil {
			log.Error("Error: %s", err.Error())
			return
		}
		return
	}

	log.Warning("unknown command")
}
