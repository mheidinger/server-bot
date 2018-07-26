package main

import (
	"sync"
	"time"

	"github.com/mheidinger/server-bot/checkers"
	"github.com/mheidinger/server-bot/services"
)

func runResultCollector(results map[string]*checkers.CheckResult, mutex *sync.Mutex, notificationChannel chan *checkers.CheckResult) {
	for true {
		for _, service := range services.Services {
			newTestTime := service.LastStarted.Add(time.Second * time.Duration(service.Interval))
			if time.Now().Before(newTestTime) {
				continue
			}
			service.LastStarted = time.Now()

			go func(s *services.Service) {
				var result *checkers.CheckResult
				checker, checkerOK := checkers.Checkers[s.CheckerName]
				if checkerOK {
					result = checker.RunTest(s)
				} else {
					checkers.CheckerNotFoundRes.TimeStamp = time.Now()
					checkers.CheckerNotFoundRes.Service = s
					result = checkers.CheckerNotFoundRes
				}

				mutex.Lock()
				exisRes, exisResOk := results[s.Name]

				if exisResOk {
					result.LastResult = exisRes
				}

				if checkerOK && checker.NeedsNotification(result) {
					notificationChannel <- result
				}

				results[s.Name] = result
				mutex.Unlock()
			}(service)
		}

		time.Sleep(time.Second)
	}
}
