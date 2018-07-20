package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mheidinger/server-bot/services"

	"github.com/mheidinger/server-bot/checkers"
	"github.com/mheidinger/server-bot/config"
)

func main() {
	completeConfig := config.LoadConfig()

	checkers.Init()
	services.Services = completeConfig.Services

	results := map[string]*checkers.CheckResult{}
	resultsMutex := &sync.Mutex{}

	runResultCollector(results, resultsMutex)

	go func() {
		for true {
			resultsMutex.Lock()
			for _, res := range results {
				fmt.Println(res)
			}
			resultsMutex.Unlock()

			time.Sleep(time.Second * 15)
		}
	}()

	StartBot(completeConfig.General.TelegramToken, completeConfig.General.BotSecret, results, resultsMutex)
}
