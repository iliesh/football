package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/iliesh/go-templates/logger"
)

func sendMessage(m sendMessageReqBodyT) error {

	if m.ChatID == 0 {
		log.Error("got: %v", m)
		return errors.New("missing chat id")
	}

	reqBytes, err := json.Marshal(m)
	if err != nil {
		log.Error("Unable to Encode the Request <%v>, Error: %s", m, err.Error())
	}

	log.Debug("post url: %s", string(reqBytes))
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

func answerCallBack(m answerCallbackQueryT) error {

	if m.CallBackQueryID == "" {
		log.Error("got: %v", m)
		return errors.New("missing call back query id")
	}

	reqBytes, err := json.Marshal(m)
	if err != nil {
		log.Error("Unable to Encode the Request <%v>, Error: %s", m, err.Error())
	}

	log.Debug("post url: %s", string(reqBytes))
	// Send a post request with your token
	resp, err := http.Post("https://api.telegram.org/"+token_football_t00002_bot+"/answerCallbackQuery", "application/json", bytes.NewBuffer(reqBytes))
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

func editMessage(m editMessageTextT) error {

	log.Debug("Edit Message: %v", m)

	reqBytes, err := json.Marshal(m)
	if err != nil {
		log.Error("Unable to Encode the Request <%v>, Error: %s", m, err.Error())
	}

	log.Debug("post url: %s", string(reqBytes))
	// Send a post request with your token
	resp, err := http.Post("https://api.telegram.org/"+token_football_t00002_bot+"/editMessageText", "application/json", bytes.NewBuffer(reqBytes))
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
