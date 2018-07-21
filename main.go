package main

import (
	"sync"

	"github.com/mheidinger/server-bot/checkers"
)

func main() {
	LoadConfig()

	checkers.Init()

	results := map[string]*checkers.CheckResult{}
	resultsMutex := &sync.Mutex{}
	notificationChannel := make(chan *checkers.CheckResult, 20)

	runResultCollector(results, resultsMutex, notificationChannel)

	StartBot(loadedConfig.General.TelegramToken, loadedConfig.General.BotSecret, results, resultsMutex, notificationChannel)
}
