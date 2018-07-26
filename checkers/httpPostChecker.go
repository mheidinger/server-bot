package checkers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/clog.v1"

	"github.com/mheidinger/server-bot/services"
)

const defaultContentType = "application/json"
const defaultExpBody = ""

// HTTPPostChecker represents a checker that checks http post requests for wanted response codes and reponse body
type HTTPPostChecker struct {
	httpClient http.Client
}

// NewHTTPPostChecker returns a new instance of the checker
func NewHTTPPostChecker() *HTTPPostChecker {
	return &HTTPPostChecker{}
}

// RunTest runs the http post test against the given service
func (checker *HTTPPostChecker) RunTest(service *services.Service) *CheckResult {
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

	var body string
	if bodyInt, ok := service.Config["body"]; ok {
		body, ok = bodyInt.(string)
		if !ok {
			WrongConfigRes.TimeStamp = time.Now()
			WrongConfigRes.Service = service
			return WrongConfigRes
		}
	}

	var contentType = defaultContentType
	if contentTypeInt, ok := service.Config["content_type"]; ok {
		if contentType, ok = contentTypeInt.(string); !ok {
			contentType = defaultContentType
		}
	}

	var expRes = defaultExpRes
	if expRespInt, ok := service.Config["expected_resp"]; ok {
		if expRes, ok = expRespInt.(int); !ok {
			expRes = defaultExpRes
		}
	}

	var expBody = defaultExpBody
	if expBodyInt, ok := service.Config["expected_body"]; ok {
		if expBody, ok = expBodyInt.(string); !ok {
			expBody = defaultExpBody
		}
	}

	clog.Trace("Run HTTPPostChecker against %s with %s: expRes: %v", url, contentType, expRes)

	t1 := time.Now()
	response, err := checker.httpClient.Post(url, contentType, bytes.NewBufferString(body))
	latency := time.Now().Sub(t1).Seconds()
	if response != nil {
		defer response.Body.Close()
	}

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	bodyString := string(bodyBytes)

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
	} else if expBody != "" && bodyString != expBody {
		res.Success = false
		resVals["error"] = "Got unexpected response body"
		resVals["_resp_body"] = bodyString
		resVals["_exp_resp_body"] = expBody
		resVals["latency"] = latency
	} else {
		res.Success = true
		resVals["resp_code"] = response.StatusCode
		resVals["latency"] = latency
	}

	res.Values = resVals
	return res
}

// NeedsNotification returns whether the result needs to be notified depending on lastResult
func (checker *HTTPPostChecker) NeedsNotification(checkResult *CheckResult) bool {
	if checkResult.LastResult != nil && checkResult.Success != checkResult.LastResult.Success {
		return true
	} else if checkResult.LastResult == nil && checkResult.Success == false {
		return true
	}

	return false
}

func (checker *HTTPPostChecker) sanitizeURL(url string) string {
	if !strings.Contains(url, "://") {
		return "http://" + url
	}

	return url
}
