package checkers

import (
	"net/http"
	"time"

	"github.com/mheidinger/server-bot/services"
)

// HTTPGetChecker represents a checker that checks http requests for wanted response codes
type HTTPGetChecker struct {
	httpClient http.Client
}

var wrongConfigVals = map[string]interface{}{"error": "Wrong configuration for HTTPGetChecker"}
var wrongConfigRes = &CheckResult{Success: false, TimeStamp: time.Now(), Values: wrongConfigVals}

// NewHTTPGetChecker returns a new instance of the checker
func NewHTTPGetChecker() *HTTPGetChecker {
	return &HTTPGetChecker{}
}

// RunTest runs the http get test against the given service
func (checker *HTTPGetChecker) RunTest(service *services.Service) *CheckResult {
	var url string
	urlInt, ok := service.Config["URL"]
	if ok {
		url, ok = urlInt.(string)
		if !ok {
			return wrongConfigRes
		}
	}

	var expRes = 200
	expRespInt, ok := service.Config["expectedResp"]
	if ok {
		expRes, _ = expRespInt.(int)
	}

	t1 := time.Now()
	response, err := checker.httpClient.Get(url)
	latency := time.Now().Sub(t1).Seconds()

	var resVals = make(map[string]interface{})
	var res = &CheckResult{TimeStamp: time.Now()}
	if err != nil {
		res.Success = false
		resVals["error"] = err.Error()
	} else if response.StatusCode != expRes {
		res.Success = false
		resVals["error"] = "Got unexpected response code"
		resVals["respCode"] = response.StatusCode
		resVals["expRespCode"] = expRes
		resVals["latency"] = latency
	} else {
		res.Success = true
		resVals["respCode"] = response.StatusCode
		resVals["latency"] = latency
	}

	res.Values = resVals
	return res
}
