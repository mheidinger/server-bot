package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mheidinger/server-bot/services"
)

// ConfigFile is the location of the config file
const ConfigFile = "data/config.json"

// GeneralConfig represents the general config needed for the bot
type GeneralConfig struct {
	TelegramToken string `json:"telegram_token,omitempty"`
	BotSecret     string `json:"bot_secret,omitempty"`
}

// CompleteConfig represents the complete config
type CompleteConfig struct {
	General  *GeneralConfig      `json:"general,omitempty"`
	Services []*services.Service `json:"services,omitempty"`
}

var config CompleteConfig

// LoadConfig loads the config, returns the general config and sets the services
func LoadConfig() *CompleteConfig {
	configFile, err := os.Open(ConfigFile)
	defer configFile.Close()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	return &config
}
