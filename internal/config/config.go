package config

import (
	"github.com/fsnotify/fsnotify"

	"github.com/go-playground/validator/v10"
	log "github.com/iliesh/go-templates/logger"
	"github.com/spf13/viper"
)

// Default Values
const (
	URLPath    string = "bot"
	ListenPort int64  = 8080
	DBHost     string = "127.0.0.1"
	DBUser     string = "root"
	DBName     string = "football"
)

// Config represents an application configuration.
type Config struct {
	// Application LogLevel (DEBUG,INFO,WARNING,ERROR,PANIC,FATAL)
	LogLevel string `mapstructure:"LOG_LEVEL"`
	// Bot URL Path, default "bot"
	URLPath string `mapstructure:"URL_PATH" validate:"required"`
	// Port on which Bot will listen, default 8080
	ListenPort string `mapstructure:"LISTEN_PORT"`
	// Bot Token, default ""
	BotToken string `mapstructure:"TG_BOT_TOKEN" validate:"required"`
	// Database Driver, default: mysql
	DBDriver string `mapstructure:"DB_DRIVER"`
	// Database Hostname, default: 127.0.0.1
	DBHost string `mapstructure:"DB_HOST"`
	// Database Port, default: 3306
	DBPort string `mapstructure:"DB_PORT"`
	// Database User, default: root
	DBUser string `mapstructure:"DB_USER"`
	// Database Password, default: ""
	DBPassword string `mapstructure:"DB_PASSWORD"`
	// Database Name, default: football
	DBName string `mapstructure:"DB_NAME"`
	// Bot SSL Certificate, default ""
	SSLCertificate string `mapstructure:"SSL_CERTIFICATE" validate:"required"`
	// Bot SSL Private Key, default ""
	SSLPrivateKey string `mapstructure:"SSL_PRIVATE_KEY" validate:"required"`
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(path string) (*Config, error) {
	log.Trace("Load App Configuration from: <%s>", path)
	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("URLPath", URLPath)
	viper.SetDefault("ListenPort", ListenPort)
	viper.SetDefault("DBHost", DBHost)
	viper.SetDefault("DBUser", DBUser)
	viper.SetDefault("DBName", DBName)

	c := Config{}

	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Error reading configuration file: <%v>", err)
		return &c, err
	}

	log.Debug("Decode Config File")
	if err := viper.Unmarshal(&c); err != nil {
		log.Error("Unable to decode configuration file, error: <%v>", err)
		return &c, err
	}

	log.Debug("Validating Config File")
	validate := validator.New()
	if err := validate.Struct(&c); err != nil {
		log.Error("Missing required attributes %v", err)
		return &c, err
	}

	log.Trace("Config was successfully initialized <%v>", c)

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	return &c, nil
}
