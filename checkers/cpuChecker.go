package checkers

import (
	"time"

	"github.com/mheidinger/server-bot/services"
	"github.com/shirou/gopsutil/cpu"
)

// CPUChecker represents a checker that checks the system cpu
type CPUChecker struct {
}

// NewCPUChecker returns a new instance of the checker
func NewCPUChecker() *CPUChecker {
	return &CPUChecker{}
}

// RunTest runs the system cpu check
func (checker *CPUChecker) RunTest(service *services.Service) *CheckResult {
	var maxCPUUsedPercent = 85.0
	if maxCPUUsedPercentInt, ok := service.Config["max_cpu_used_percentage"]; ok {
		maxCPUUsedPercent, _ = maxCPUUsedPercentInt.(float64)
	}

	var measureInterval = 2
	if measureIntervalInt, ok := service.Config["measure_interval"]; ok {
		measureInterval, _ = measureIntervalInt.(int)
	}

	cpuPercent, err := cpu.Percent(time.Second*time.Duration(measureInterval), false)

	var resVals = make(map[string]interface{})
	var res = &CheckResult{Service: service, TimeStamp: time.Now()}
	if err != nil {
		res.Success = false
		resVals["error"] = err.Error()
	} else {
		if cpuPercent[0] >= maxCPUUsedPercent {
			res.Success = false
		} else {
			res.Success = true
		}
		resVals["used_percentage"] = cpuPercent[0]
	}

	res.Values = resVals
	return res
}

// NeedsNotification returns whether the result needs to be notified depending on lastResult
func (checker *CPUChecker) NeedsNotification(checkResult *CheckResult) bool {
	if checkResult.LastResult != nil && checkResult.Success != checkResult.LastResult.Success {
		return true
	} else if checkResult.LastResult == nil && checkResult.Success == false {
		return true
	}

	return false
}
