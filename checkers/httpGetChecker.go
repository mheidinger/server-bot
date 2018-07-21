package checkers

import (
	"net/http"
	"strings"
	"time"

	"github.com/mheidinger/server-bot/services"
)

// HTTPGetChecker represents a checker that checks http requests for wanted response codes
type HTTPGetChecker struct {
	httpClient http.Client
}

// NewHTTPGetChecker returns a new instance of the checker
func NewHTTPGetChecker() *HTTPGetChecker {
	return &HTTPGetChecker{}
}

// RunTest runs the http get test against the given service
func (checker *HTTPGetChecker) RunTest(service *services.Service) *CheckResult {
	var url string
	if urlInt, ok := service.Config["url"]; ok {
		url, ok = urlInt.(string)
		if !ok {
			WrongConfigRes.TimeStamp = time.Now()
			WrongConfigRes.Service = service
			return WrongConfigRes
		}
	}
	url = checker.sanitizeURL(url)

	var expRes = 200
	if expRespInt, ok := service.Config["expected_resp"]; ok {
		expRes, _ = expRespInt.(int)
	}

	t1 := time.Now()
	response, err := checker.httpClient.Get(url)
	latency := time.Now().Sub(t1).Seconds()
	defer response.Body.Close()

	var resVals = make(map[string]interface{})
	var res = &CheckResult{Service: service, TimeStamp: time.Now()}
	if err != nil {
		res.Success = false
		resVals["error"] = err.Error()
	} else if response.StatusCode != expRes {
		res.Success = false
		resVals["error"] = "Got unexpected response code"
		resVals["resp_code"] = response.StatusCode
		resVals["exp_resp_code"] = expRes
		resVals["latency"] = latency
	} else {
		res.Success = true
		resVals["resp_code"] = response.StatusCode
		resVals["latency"] = latency
	}

	res.Values = resVals
	return res
}

func (checker *HTTPGetChecker) sanitizeURL(url string) string {
	if !strings.Contains(url, "://") {
		return "http://" + url
	}

	return url
}
