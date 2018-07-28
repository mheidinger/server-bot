package checkers

import (
	"time"

	"gopkg.in/clog.v1"

	"github.com/mheidinger/server-bot/services"
	"github.com/shirou/gopsutil/mem"
)

const defaultMaxMemUsedPercent = 85.0

// MemChecker represents a checker that checks the system memory
type MemChecker struct {
}

// NewMemChecker returns a new instance of the checker
func NewMemChecker() *MemChecker {
	return &MemChecker{}
}

// RunTest runs the system memory check
func (checker *MemChecker) RunTest(service *services.Service) *CheckResult {
	var maxMemUsedPercent = defaultMaxMemUsedPercent
	if maxMemUsedPercentInt, ok := service.Config["max_mem_used_percentage"]; ok {
		if maxMemUsedPercent, ok = maxMemUsedPercentInt.(float64); !ok {
			maxMemUsedPercent = defaultMaxMemUsedPercent
		}
	}

	clog.Trace("Run MemChecker with: maxMemUsedPercent: %v", maxMemUsedPercent)

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
	return defaultNeedsNotification(checkResult)
}
