package main

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
	General       *GeneralConfig      `json:"general,omitempty"`
	TelegramUsers []int               `json:"telegram_users,omitempty"`
	Services      []*services.Service `json:"services,omitempty"`
}

var loadedConfig CompleteConfig

// LoadConfig loads the config, returns the general config and sets the services
func LoadConfig() {
	configFile, err := os.Open(ConfigFile)
	defer configFile.Close()
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&loadedConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	services.Services = loadedConfig.Services
	TelegramUsers = loadedConfig.TelegramUsers
}

// WriteConfig writes the currently used config into the config file
func WriteConfig() {
	loadedConfig.Services = services.Services
	loadedConfig.TelegramUsers = TelegramUsers

	configFile, err := os.Create(ConfigFile)
	defer configFile.Close()
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}

	jsonWriter := json.NewEncoder(configFile)
	jsonWriter.SetIndent("", "\t")
	err = jsonWriter.Encode(&loadedConfig)
	if err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}
