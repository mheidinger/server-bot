package checkers

import (
	"github.com/mheidinger/server-bot/services"
)

// Checker declares all functions needed for a checker to be used
type Checker interface {
	RunTest(service *services.Service) *CheckResult
	NeedsNotification(checkResult *CheckResult) bool
}

// Default Error Results
var wrongConfigVals = map[string]interface{}{"error": "Wrong configuration for HTTPGetChecker"}
var checkerNotFoundVals = map[string]interface{}{"error": "Checker with given name not found"}

// WrongConfigRes represents that the checker did not get correct config vals
var WrongConfigRes = &CheckResult{Success: false, Values: wrongConfigVals}

// CheckerNotFoundRes represents that a checker with the given name was not found
var CheckerNotFoundRes = &CheckResult{Success: false, Values: checkerNotFoundVals}

// Checkers stores all instances of checkers
var Checkers map[string]Checker

// Init initializes all checkers
func Init() {
	Checkers = make(map[string]Checker, 1)
	Checkers["HTTPGetChecker"] = NewHTTPGetChecker()
	Checkers["HTTPPostChecker"] = NewHTTPPostChecker()
	Checkers["TCPDialChecker"] = NewTCPDialChecker()
}
