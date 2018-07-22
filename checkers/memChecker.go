package checkers

import (
	"time"

	"github.com/mheidinger/server-bot/services"
	"github.com/shirou/gopsutil/mem"
)

// MemChecker represents a checker that checks the system memory
type MemChecker struct {
}

// NewMemChecker returns a new instance of the checker
func NewMemChecker() *MemChecker {
	return &MemChecker{}
}

// RunTest runs the system memory check
func (checker *MemChecker) RunTest(service *services.Service) *CheckResult {
	var maxMemUsedPercent = 85.0
	if maxMemUsedPercentInt, ok := service.Config["max_mem_used_percentage"]; ok {
		maxMemUsedPercent, _ = maxMemUsedPercentInt.(float64)
	}

	memory, err := mem.VirtualMemory()

	var resVals = make(map[string]interface{})
	var res = &CheckResult{Service: service, TimeStamp: time.Now()}
	if err != nil {
		res.Success = false
		resVals["error"] = err.Error()
	} else {
		if memory.UsedPercent >= maxMemUsedPercent {
			res.Success = false
		} else {
			res.Success = true
		}
		resVals["used_percentage"] = memory.UsedPercent
		resVals["total"] = memory.Total
		resVals["available"] = memory.Available
		resVals["used"] = memory.Used
	}

	res.Values = resVals
	return res
}

// NeedsNotification returns whether the result needs to be notified depending on lastResult
func (checker *MemChecker) NeedsNotification(checkResult *CheckResult) bool {
	if checkResult.LastResult != nil && checkResult.Success != checkResult.LastResult.Success {
		return true
	} else if checkResult.LastResult == nil && checkResult.Success == false {
		return true
	}

	return false
}
