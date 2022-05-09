package main

import (
	"net/http"
	"os"

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
	token_football_t00002_bot = ""
)

type Config struct {
	URLPath        string `mapstructure:"URL_PATH"`
	ListenPort     string `mapstructure:"LISTEN_PORT"`
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

func main() {
	log.Version = "4.0.2"
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

	log.Debug("Listening on port: %s", config.ListenPort)
	err = http.ListenAndServeTLS(":"+config.ListenPort, config.SSLCertificate, config.SSLPrivateKey, nil)
	if err != nil {
		log.Error("ListenAndServe Error: %s", err.Error())
		os.Exit(1)
	}
}
