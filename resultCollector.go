package main

import (
	"sync"
	"time"

	"github.com/mheidinger/server-bot/checkers"
	"github.com/mheidinger/server-bot/services"
)

func runResultCollector(mutex *sync.Mutex, results map[string]*checkers.CheckResult) {
	go func() {
		for true {
			for _, service := range services.GetServices() {
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
				if checker, ok := checkers.Checkers[service.CheckerName]; ok {
					result = checker.RunTest(service)
				} else {
					checkers.CheckerNotFoundRes.TimeStamp = time.Now()
					result = checkers.CheckerNotFoundRes
				}

				mutex.Lock()
				if exisResOk {
					result.LastResult = exisRes
				}

				results[service.Name] = result
				mutex.Unlock()
			}

			time.Sleep(time.Second)
		}
	}()
}
