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

	runResultCollector(results, resultsMutex)

	StartBot(loadedConfig.General.TelegramToken, loadedConfig.General.BotSecret, results, resultsMutex)
}
