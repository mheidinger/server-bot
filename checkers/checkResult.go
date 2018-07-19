package checkers

import (
	"time"
)

// CheckResult represents the result of one of the checkers
type CheckResult struct {
	Success    bool
	TimeStamp  time.Time
	Values     map[string]interface{}
	LastResult *CheckResult
}
