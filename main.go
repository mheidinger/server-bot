package main

import (
	"fmt"

	"github.com/mheidinger/server-bot/checkers"
	"github.com/mheidinger/server-bot/services"
)

func main() {
	websiteServiceConfig := map[string]interface{}{"URL": "https://max-heidinger.de", "expectedResponse": 200}
	websiteService := &services.Service{Name: "Own Website", CheckerName: "HTTPGetChecker", Config: websiteServiceConfig}

	httpGetChecker := checkers.NewHTTPGetChecker()

	checkRes := httpGetChecker.RunTest(websiteService)

	fmt.Println(checkRes)
}
