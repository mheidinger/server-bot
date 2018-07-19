package main

import (
	"fmt"

	"github.com/mheidinger/server-bot/checkers"
	"github.com/mheidinger/server-bot/services"
)

func main() {
	websiteServiceConfig := map[string]interface{}{"URL": "max-heidinger.de", "expectedResponse": 200}
	websiteService := &services.Service{Name: "Own Website", CheckerName: "HTTPGetChecker", Config: websiteServiceConfig}
	checkServices := []*services.Service{websiteService}

	httpGetChecker := checkers.NewHTTPGetChecker()

	results := map[string]*checkers.CheckResult{}

	for _, service := range checkServices {
		var checkRes *checkers.CheckResult
		switch service.CheckerName {
		case "HTTPGetChecker":
			checkRes = httpGetChecker.RunTest(service)
		}

		if exisRes, ok := results[service.Name]; ok {
			checkRes.LastResult = exisRes
		}
		results[service.Name] = checkRes
	}

	for _, res := range results {
		fmt.Println(res)
	}
}
