package checkers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mheidinger/server-bot/services"
)

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

	var contentType = "application/json"
	if contentTypeInt, ok := service.Config["content_type"]; ok {
		contentType, _ = contentTypeInt.(string)
	}

	var expRes = 200
	if expRespInt, ok := service.Config["expected_resp"]; ok {
		expRes, _ = expRespInt.(int)
	}

	var expBody = ""
	if expBodyInt, ok := service.Config["expected_body"]; ok {
		expBody, _ = expBodyInt.(string)
	}

	t1 := time.Now()
	response, err := checker.httpClient.Post(url, contentType, bytes.NewBufferString(body))
	latency := time.Now().Sub(t1).Seconds()
	defer response.Body.Close()

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
		resVals["resp_body"] = bodyString
		resVals["exp_resp_body"] = expBody
		resVals["latency"] = latency
	} else {
		res.Success = true
		resVals["resp_code"] = response.StatusCode
		resVals["latency"] = latency
	}

	res.Values = resVals
	return res
}

func (checker *HTTPPostChecker) sanitizeURL(url string) string {
	if !strings.Contains(url, "://") {
		return "http://" + url
	}

	return url
}
