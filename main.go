package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/iliesh/go-templates/logger"
	"github.com/spf13/viper"
)

const (
	htmlHelpText = `
<i>- Create a new Game:</i>
<code>/new_game</code>
<i>- Cancel current Game:</i>
<code>/stop_game</code>
<i>- Adding you to the Game:</i>
<code>/add</code>
<i>- Removing you from the Game:</i>
<code>/cancel</code>
<i>- Show Game Status:</i>
<code>/status</code>
<i>- Show Help Message:</i>
<code>/help</code>
`
)

var (
	token_football_t00002_bot = "telegram_token"
)

type Config struct {
	URLPath        string `mapstructure:"URL_PATH"`
	BotToken       string `mapstructure:"TG_BOT_TOKEN"`
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBDatabase     string `mapstructure:"DB_DATABASE"`
	SSLCertificate string `mapstructure:"SSL_CERTIFICATE"`
	SSLPrivateKey  string `mapstructure:"SSL_PRIVATE_KEY"`
}

// userT Struct
type userT struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

// chatT struct
type chatT struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// voiceT struct
type voiceT struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

// messageT Struct
type messageT struct {
	MessageID            int64  `json:"message_id"`
	From                 userT  `json:"from"`
	ForwardFrom          userT  `json:"forward_from"`
	Chat                 chatT  `json:"chat"`
	SenderChat           chatT  `json:"sender_chat"`
	ForwardFromChat      chatT  `json:"forward_from_chat"`
	ForwardFromMessageID int64  `json:"forward_from_message_id"`
	ForwardSignature     string `json:"forward_signature"`
	ForwardSenderName    string `json:"forward_sender_name"`
	ForwardDate          int64  `json:"forward_date"`
	Date                 int64  `json:"date"`
	Text                 string `json:"text"`
	Entities             []struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Type   string `json:"type"`
	} `json:"entities"`
	ViaBot          userT  `json:"via_bot"`
	EditDate        int64  `json:"edit_date"`
	MediaGroupID    string `json:"media_group_id"`
	AuthorSignature string `json:"author_signature"`
	Voice           voiceT `json:"voice"`
	Caption         string `json:"caption"`
}

type editMessageTextT struct {
	ChatID          string                `json:"chat_id,omitempty"`
	MessageID       int64                 `json:"message_id,omitempty"`
	InlineMessageID string                `json:"inline_message_id,omitempty"`
	Text            string                `json:"text,omitempty"`
	ParseMode       string                `json:"parse_mode,omitempty"`
	ReplyMarkup     inlineKeyboardMarkupT `json:"reply_markup,omitempty"`
}

// callBackQueryT Struct
type callBackQueryT struct {
	ID              string   `json:"id"`
	From            userT    `json:"from"`
	Message         messageT `json:"message,omitempty"`
	InlineMessageID string   `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
}

type answerCallbackQueryT struct {
	CallBackQueryID string `json:"callback_query_id,omitempty"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int64  `json:"cache_time,omitempty"`
}

type webHookReqBodyT struct {
	UpdateID      int            `json:"update_id"`
	Message       messageT       `json:"message,omitempty"`
	CallBackQuery callBackQueryT `json:"callback_query,omitempty"`
}

// sendMessageReqBodyT Create a struct to conform to the JSON body
// of the send message request
// https://core.telegram.org/bots/api#sendmessage
type sendMessageReqBodyT struct {
	ChatID      int64                 `json:"chat_id,omitempty"`
	Text        string                `json:"text,omitempty"`
	ParseMode   string                `json:"parse_mode,omitempty"`
	ReplyMarkup inlineKeyboardMarkupT `json:"reply_markup,omitempty"`
}

type inlineKeyboardMarkupT struct {
	InlineKeyboard [][]inlineKeyboardButtonT `json:"inline_keyboard,omitempty"`
}

type inlineKeyboardButtonT struct {
	Text         string `json:"text,omitempty"`
	URL          string `json:"url,omitempty"`
	CallBackData string `json:"callback_data,omitempty"`
}

func main() {
	log.Version = "4.0.1"
	log.Info("Start Application")
	log.AppName = "Football Manager"

	config, err := LoadConfig(".")
	if err != nil {
		log.Error("cannot load config:", err)
		return
	}

	token_football_t00002_bot = config.BotToken

	log.Debug("Config URL Path: %s", config.URLPath)

	http.HandleFunc(config.URLPath, HandlerBot)
	http.HandleFunc("/", HandlerRoot)

	err = http.ListenAndServeTLS(":90", config.SSLCertificate, config.SSLPrivateKey, nil)
	if err != nil {
		log.Error("ListenAndServe Error: %s", err.Error())
		os.Exit(1)
	}
}

// HandlerRoot is called everytime telegram sends us a webhook event
func HandlerRoot(res http.ResponseWriter, req *http.Request) {
	log.Info("request: %v", req.URL)
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
				log.Debug("Check if User: %s %s (%s) is allowed to create a new game", body.Message.From.FirstName, body.Message.From.LastName, body.Message.From.Username)
				userPerm, err := newGamePermission(body.Message.From.ID)
				if err != nil {
					log.Error("Error: %s", err.Error())
					m := sendMessageReqBodyT{ChatID: body.Message.From.ID, Text: "\u26A0 Internal Error"}
					err = sendMessage(m)
					if err != nil {
						log.Error("Error: %s", err.Error())
					}
					return
				}
				if !userPerm {
					log.Warning("User ID: <%d> is not allowed to create a new game", body.Message.From.ID)
					m := sendMessageReqBodyT{ChatID: body.Message.From.ID, Text: "\u26A0 Sorry, only admins can create or cancel games"}
					err = sendMessage(m)
					if err != nil {
						log.Error("Error: %s", err.Error())
					}
					return
				}
				log.Debug("Select the Game Date")
				gameDate, err := selectDate(body.Message.From.ID)
				if err != nil {
					log.Error("Error: %s", err.Error())
					return
				}
				log.Info("Selected Game Date: %s", gameDate)
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
		log.Debug("selected date: %s", body.CallBackQuery.Data)

		if body.CallBackQuery.Data == "month" {
			log.Warning("unable to select month as a value")
			m := answerCallbackQueryT{CallBackQueryID: body.CallBackQuery.ID, Text: "Cannot select month, please select the date", ShowAlert: true}
			err = answerCallBack(m)
			if err != nil {
				log.Error("Error: %s", err.Error())
			}
		}
		m := answerCallbackQueryT{CallBackQueryID: body.CallBackQuery.ID, Text: "Selected Date: " + body.CallBackQuery.Data, ShowAlert: false}
		err = answerCallBack(m)
		if err != nil {
			log.Error("Error: %s", err.Error())
		}

		d := editMessageTextT{MessageID: body.CallBackQuery.Message.MessageID, Text: "date: 01.01.01", ReplyMarkup: inlineKeyboardMarkupT{}}
		err = editMessage(d)
		if err != nil {
			log.Error("Error: %s", err.Error())
		}
	}

	log.Warning("unknown command")
}

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

func newGamePermission(uid int64) (bool, error) {
	return true, nil
}

func selectDate(id int64) (string, error) {
	monthName := inlineKeyboardButtonT{Text: "April 2021", CallBackData: "month"}
	monthButton := []inlineKeyboardButtonT{monthName}

	weekDayMo := inlineKeyboardButtonT{Text: "Mo", CallBackData: "monday"}
	weekDayTu := inlineKeyboardButtonT{Text: "Tu", CallBackData: "tuesday"}
	weekDayWe := inlineKeyboardButtonT{Text: "We", CallBackData: "wednesday"}
	weekDayTh := inlineKeyboardButtonT{Text: "Th", CallBackData: "thursday"}
	weekDayFr := inlineKeyboardButtonT{Text: "Fr", CallBackData: "friday"}
	weekDaySa := inlineKeyboardButtonT{Text: "Sa", CallBackData: "saturday"}
	weekDaySu := inlineKeyboardButtonT{Text: "Su", CallBackData: "sunday"}
	weekDaysButton := []inlineKeyboardButtonT{weekDayMo, weekDayTu, weekDayWe, weekDayTh, weekDayFr, weekDaySa, weekDaySu}
	firstLineDate := []inlineKeyboardButtonT{}
	secondLineDate := []inlineKeyboardButtonT{}

	currTime := time.Now()
	weekDayNumber := currTime.Weekday()
	if weekDayNumber == 0 {
		weekDayNumber = 7
	}
	for i := 1; i < 15; i++ {
		if i < int(weekDayNumber) {
			firstLineDate = append(firstLineDate, inlineKeyboardButtonT{Text: "\u2718", CallBackData: "null"})
			continue
		}
		callBackDate := currTime.AddDate(0, 0, i-int(weekDayNumber))
		_, _, date := callBackDate.Date()
		if i < 8 {
			firstLineDate = append(firstLineDate, inlineKeyboardButtonT{Text: strconv.Itoa(date), CallBackData: callBackDate.Format("2006-01-02")})
			continue
		}
		secondLineDate = append(secondLineDate, inlineKeyboardButtonT{Text: strconv.Itoa(date), CallBackData: callBackDate.Format("2006-01-02")})
	}

	b := [][]inlineKeyboardButtonT{monthButton, weekDaysButton, firstLineDate, secondLineDate}
	c := inlineKeyboardMarkupT{b}
	m := sendMessageReqBodyT{ChatID: id, Text: "Please select the date:", ReplyMarkup: c}
	log.Info("Sending Message: %v", m)
	err := sendMessage(m)
	if err != nil {
		return "", err
	}
	return "", nil
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
