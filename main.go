package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	debug bool
	teamA []int64
	teamB []int64
)

var weekDayToInt = map[string]int{
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
	"Sunday":    7,
}

const (
	version         string = "3.1"
	sqlDateTimeForm        = "2006-01-02 15:04:05"
	botID                  = "football_t00001_bot"
	htmlHelpText           = `
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

// TeleBot Struct
type TeleBot struct {
	botAPI  *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	event   Event
}

// Event Struct
type Event struct {
	eventUserID      string
	eventOrganizer   string
	eventDate        string
	eventTime        string
	eventTeamPlayers string
}

func dbConn(d string) (db *sql.DB) {
	dbDriver := "mysql"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbUser := "root"
	dbPassword := "password"
	dbDatabase := d

	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbDatabase)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func (t *TeleBot) sendAnswerCallbackQuery() {
	for update := range t.updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.CallbackQuery != nil {
			resCallBackQuery := strings.SplitN(update.CallbackQuery.Data, "#", -1)

			if len(resCallBackQuery) == 1 {
				log.Printf("Missing Function %v\n", resCallBackQuery)
				t.botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
				continue
			}

			switch resCallBackQuery[0] {
			case "callback_data":
				del := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				config := tgbotapi.NewCallback(update.CallbackQuery.ID, strconv.Itoa(update.CallbackQuery.Message.MessageID))
				go t.botAPI.AnswerCallbackQuery(config)
				go t.botAPI.Send(del)
			case "eventDate":
				if resCallBackQuery[1][0] < 5 {
					del := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					config := tgbotapi.NewCallback(update.CallbackQuery.ID, strconv.Itoa(update.CallbackQuery.Message.MessageID))
					go t.botAPI.AnswerCallbackQuery(config)
					go t.botAPI.Send(del)

					switch resCallBackQuery[1][0] {
					case 0:
						t.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Cannot select date in the Past, Please try again"))
						t.botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Cannot select date in the Past, Please try again"))
					case 1:
						t.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Cannot select Month as date, Please try again"))
						t.botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Cannot select Month as date, Please try again"))
					case 2:
						t.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Cannot select Day as date, Please try again"))
						t.botAPI.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Cannot select Day as date, Please try again"))
					}
					break
				}
				fmt.Printf("GOT CALL BACK QUERRY! %v\n", resCallBackQuery)
				t.eventDate(update)
			case "eventTime":
				fmt.Printf("GOT CALL BACK QUERRY! %v\n", resCallBackQuery)
				t.eventTime(update)
			case "pitchSize":
				fmt.Printf("GOT CALL BACK QUERRY! %v\n", resCallBackQuery)
				t.pitchSize(update)
				t.dbUpdateEventAdd(update)
			}
		}

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			var err error

			switch update.Message.Command() {
			case "start":
				if !t.start(update) {
					continue
				}
			case "new_game":
				if t.eventAdd(update) {
					t.addCalendar(update)
				}
			case "stop_game":
				if !t.eventDel(update) {
					continue
				}
			case "add":
				if !t.addPlayer(update) {
					continue
				}
			case "cancel":
				if !t.cancelPlayer(update) {
					continue
				}
			case "status":
				err := t.gameStatus(update)
				if err != nil {
					fmt.Printf("ERROR!, %+v\n", err)
				}
			case "help":
				msg.Text = htmlHelpText
				msg.ParseMode = tgbotapi.ModeHTML
				if _, err := t.botAPI.Send(msg); err != nil {
					log.Panic(err)
				}
			default:
				msg.Text = htmlHelpText
				msg.ParseMode = tgbotapi.ModeHTML
				if _, err := t.botAPI.Send(msg); err != nil {
					log.Panic(err)
				}
			}
			if err != nil {
				fmt.Printf("ERROR!, %+v\n", err)
			}

			if msg.Text == "" {
				msg.Text = "Something goes wrong"
			}

		}
	}
}

// HandleError Function
func (t *TeleBot) HandleError(err error) (msg string) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		fmt.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
		return "(" + strconv.Itoa(line) + ")"
	}
	return "\u26A0 Unknown Error"
}

func main() {
	bot, err := tgbotapi.NewBotAPI("BOT-API")
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	teleBot := TeleBot{
		botAPI: bot,
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	teleBot.updates, err = bot.GetUpdatesChan(u)
	teleBot.sendAnswerCallbackQuery()
}
