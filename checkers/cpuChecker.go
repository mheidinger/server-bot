package checkers

import (
	"time"

	"gopkg.in/clog.v1"

	"github.com/mheidinger/server-bot/services"
	"github.com/shirou/gopsutil/cpu"
)

const defaultMaxCPUUsedPercent = 85.0
const defaultMeasureDuration = 60

// CPUChecker represents a checker that checks the system cpu
type CPUChecker struct {
}

// NewCPUChecker returns a new instance of the checker
func NewCPUChecker() *CPUChecker {
	return &CPUChecker{}
}

// RunTest runs the system cpu check
func (checker *CPUChecker) RunTest(service *services.Service) *CheckResult {
	var maxCPUUsedPercent = defaultMaxCPUUsedPercent
	if maxCPUUsedPercentInt, ok := service.Config["max_cpu_used_percentage"]; ok {
		if maxCPUUsedPercent, ok = maxCPUUsedPercentInt.(float64); !ok {
			maxCPUUsedPercent = defaultMaxCPUUsedPercent
		}
	}

	var measureDuration = defaultMeasureDuration
	if measureDurationInt, ok := service.Config["measure_duration"]; ok {
		if measureDuration, ok = measureDurationInt.(int); !ok {
			measureDuration = defaultMeasureDuration
		}
	}

	clog.Trace("Run CPUChecker with: maxCPUUsedPercentage: %v; measureDuration: %v", maxCPUUsedPercent, measureDuration)

	cpuPercent, err := cpu.Percent(time.Second*time.Duration(measureDuration), false)

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
