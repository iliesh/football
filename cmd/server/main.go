package main

import (
	"net/http"
	"os"

	"github.com/iliesh/football/internal/config"
	hc "github.com/iliesh/football/internal/healtcheck"
	db "github.com/iliesh/football/internal/mysql"
	tg "github.com/iliesh/football/internal/telegram"
	log "github.com/iliesh/go-templates/logger"
)

var (
	Version string = "4.1.0"
	AppName string = "Football Bot"
)

func init() {
	// Initializing Logger
	log.AppName = AppName
	log.Version = Version
	log.Env = "dev"
	log.Color = true
	log.LogLevel = "debug"
}

func main() {
	log.Debug("Starting Application")

	// load application configurations
	cfg, err := config.Load("./.env")
	if err != nil {
		os.Exit(-1)
	}

	// Connecting to the DB
	db, err := db.Open(cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword)
	if err != nil {
		os.Exit(-1)
	}

	log.Debug("Connection to the DB was successfully established, %v", db)

	http.HandleFunc("/", HandlerRoot)
	http.HandleFunc("/healthcheck", hc.Handler)

	bot := &tg.BotAPI{DB: db, Token: cfg.BotToken}

	http.HandleFunc(cfg.URLPath, func(w http.ResponseWriter, r *http.Request) {
		tg.HandlerBot(w, r, bot)
	})

	log.Debug("Listening on port: %s", cfg.ListenPort)
	err = http.ListenAndServeTLS(":"+cfg.ListenPort, cfg.SSLCertificate, cfg.SSLPrivateKey, nil)
	if err != nil {
		log.Error("ListenAndServe Error: %s", err.Error())
		os.Exit(-1)
	}
}

// HandlerRoot is called everytime telegram sends us a webhook event
func HandlerRoot(res http.ResponseWriter, req *http.Request) {
	log.Info("request: %v", req.URL)
}
