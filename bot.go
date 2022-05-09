package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/iliesh/go-templates/logger"
)

func showHelp(id int64) error {
	h := sendMessageReqBodyT{}
	h.ChatID = id
	h.ParseMode = "html"
	h.Text = htmlHelpText

	reqBytes, err := json.Marshal(h)
	if err != nil {
		log.Error("Unable to Encode the Request <%v>, Error: %s", h, err.Error())
	}

	// Send a post request with your token
	resp, err := http.Post("https://api.telegram.org/"+token_football_t00002_bot+"/sendmessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	log.Debug("Response: %v", resp)

	if resp.StatusCode != 200 {
		log.Error("bad http status: %d", resp.StatusCode)
		return errors.New("bad http status code")
	}
	return nil
}

func newGame(body *webHookReqBodyT) error {

	log.Debug("Check if User: %s %s (%s) is allowed to create a new game", body.Message.From.FirstName, body.Message.From.LastName, body.Message.From.Username)
	userPerm, err := newGamePermission(body.Message.From.ID)
	if err != nil {
		log.Error("Error: %s", err.Error())
		m := sendMessageReqBodyT{ChatID: body.Message.From.ID, Text: "\u26A0 Internal Error"}
		err = sendMessage(m)
		if err != nil {
			log.Error("Error: %s", err.Error())
		}
		return errors.New("permission error")
	}
	if !userPerm {
		log.Warning("User ID: <%d> is not allowed to create a new game", body.Message.From.ID)
		m := sendMessageReqBodyT{ChatID: body.Message.From.ID, Text: "\u26A0 Sorry, only admins can create or cancel games"}
		err = sendMessage(m)
		if err != nil {
			log.Error("Error: %s", err.Error())
		}
		return nil
	}
	log.Debug("Select the Game Date")
	err = selectDate(body.Message.From.ID)
	if err != nil {
		log.Error("Error: %s", err.Error())
		return err
	}
	return nil
}
