package checkers

import (
	"time"

	"gopkg.in/clog.v1"

	"github.com/mheidinger/server-bot/services"
	"github.com/shirou/gopsutil/host"
)

const defaultMaxTemp = 90.0

// TempChecker represents a checker that checks the system temperature
type TempChecker struct {
}

// NewTempChecker returns a new instance of the checker
func NewTempChecker() *TempChecker {
	return &TempChecker{}
}

// RunTest runs the system temperature check
func (checker *TempChecker) RunTest(service *services.Service) *CheckResult {
	var maxTemp = defaultMaxTemp
	if maxTempInt, ok := service.Config["max_temperature"]; ok {
		if maxTemp, ok = maxTempInt.(float64); !ok {
			maxTemp = defaultMaxTemp
		}
	}

	clog.Trace("Run TempChecker with: maxTemp: %v", maxTemp)

	temperatures, err := host.SensorsTemperatures()

	var resVals = make(map[string]interface{})
	var res = &CheckResult{Service: service, TimeStamp: time.Now()}
	if err != nil {
		res.Success = false
		resVals["error"] = err.Error()
	} else {
		res.Success = true
		for _, temp := range temperatures {
			if temp.Temperature >= maxTemp {
				res.Success = false
			}
			resVals[temp.SensorKey] = temp.Temperature
		}
	}

	res.Values = resVals
	return res
}

// NeedsNotification returns whether the result needs to be notified depending on lastResult
func (checker *TempChecker) NeedsNotification(checkResult *CheckResult) bool {
	return defaultNeedsNotification(checkResult)
}
