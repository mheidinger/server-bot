package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mheidinger/server-bot/checkers"
)

func main() {
	checkers.Init()

	results := map[string]*checkers.CheckResult{}
	resultsMutex := &sync.Mutex{}

	runResultCollector(resultsMutex, results)

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

	time.Sleep(time.Minute * 5)
}
