package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/mheidinger/server-bot/checkers"
	clog "gopkg.in/clog.v1"
)

func main() {
	initLogging()

	LoadConfig()

	checkers.Init()

	results := map[string]*checkers.CheckResult{}
	resultsMutex := &sync.Mutex{}
	notificationChannel := make(chan *checkers.CheckResult, 20)

	runResultCollector(results, resultsMutex, notificationChannel)

	StartBot(loadedConfig.General.TelegramToken, loadedConfig.General.BotSecret, results, resultsMutex, notificationChannel)
}

func initLogging() {
	err := clog.New(clog.CONSOLE, clog.ConsoleConfig{})
	if err != nil {
		fmt.Printf("Fail to create new logger: %v\n", err)
		os.Exit(1)
	}
}
