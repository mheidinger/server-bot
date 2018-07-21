package checkers

import (
	"time"

	"github.com/mheidinger/server-bot/services"
)

// CheckResult represents the result of one of the checkers
type CheckResult struct {
	Service    *services.Service
	Success    bool
	TimeStamp  time.Time
	Values     map[string]interface{}
	LastResult *CheckResult
}
