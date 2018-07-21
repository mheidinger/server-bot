package checkers

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mheidinger/server-bot/services"
)

// TCPDialChecker represents a checker that checks for succesfull tcp dials
type TCPDialChecker struct {
	httpClient http.Client
}

// NewTCPDialChecker returns a new instance of the checker
func NewTCPDialChecker() *TCPDialChecker {
	return &TCPDialChecker{}
}

// RunTest runs the tcp dial check against a service
func (checker *TCPDialChecker) RunTest(service *services.Service) *CheckResult {
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

	var timeout = 5
	if timeoutInt, ok := service.Config["timeout"]; ok {
		timeout, _ = timeoutInt.(int)
	}

	t1 := time.Now()
	conn, err := net.DialTimeout("tcp", url, time.Duration(timeout)*time.Second)
	latency := time.Now().Sub(t1).Seconds()
	defer conn.Close()

	var resVals = make(map[string]interface{})
	var res = &CheckResult{Service: service, TimeStamp: time.Now()}
	if err != nil {
		res.Success = false
		resVals["error"] = err.Error()
	} else {
		res.Success = true
		resVals["latency"] = latency
	}

	res.Values = resVals
	return res
}

func (checker *TCPDialChecker) sanitizeURL(url string) string {
	if !strings.Contains(url, "://") {
		return "http://" + url
	}

	return url
}
