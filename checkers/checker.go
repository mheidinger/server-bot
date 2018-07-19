package checkers

import "github.com/mheidinger/server-bot/services"

// Checker declares all functions needed for a checker to be used
type Checker interface {
	RunTest(service *services.Service) *CheckResult
}
