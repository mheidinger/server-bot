package main

import (
	"sync"
	"time"

	"github.com/mheidinger/server-bot/checkers"
	"github.com/mheidinger/server-bot/services"
)

func runResultCollector(results map[string]*checkers.CheckResult, mutex *sync.Mutex, notificationChannel chan *checkers.CheckResult) {
	go func() {
		for true {
			for _, service := range services.Services {
				mutex.Lock()
				exisRes, exisResOk := results[service.Name]

				if exisResOk {
					newTestTime := exisRes.TimeStamp.Add(time.Second * time.Duration(service.Interval))
					if time.Now().Before(newTestTime) {
						mutex.Unlock()
						continue
					}
				}
				mutex.Unlock()

				var result *checkers.CheckResult
				checker, checkerOK := checkers.Checkers[service.CheckerName]
				if checkerOK {
					result = checker.RunTest(service)
				} else {
					checkers.CheckerNotFoundRes.TimeStamp = time.Now()
					checkers.CheckerNotFoundRes.Service = service
					result = checkers.CheckerNotFoundRes
				}

				mutex.Lock()
				if exisResOk {
					result.LastResult = exisRes
				}

				if checkerOK && checker.NeedsNotification(result) {
					notificationChannel <- result
				}

				results[service.Name] = result
				mutex.Unlock()
			}

			time.Sleep(time.Second)
		}
	}()
}
