package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	sw "github.com/Azure-Samples/openhack-devops-team/apis/trips/tripsgo"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

var (
	webServerPort    = flag.String("webServerPort", getEnv("WEB_PORT", "8080"), "web server port")
	webServerBaseURI = flag.String("webServerBaseURI", getEnv("WEB_SERVER_BASE_URI", "changeme"), "base portion of server uri")
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {

	var debug, present = os.LookupEnv("DEBUG_LOGGING")
	var aikey, present2 = os.LookupEnv("APPINSIGHTS_INSTRUMENTATIONKEY")

	telemetryConfig := appinsights.NewTelemetryConfiguration(aikey)

	// Configure how many items can be sent in one call to the data collector:
	telemetryConfig.MaxBatchSize = 8192

	// Configure the maximum delay before sending queued telemetry:
	telemetryConfig.MaxBatchInterval = 2 * time.Second

	client := appinsights.NewTelemetryClientFromConfig(telemetryConfig)

	if present && debug == "true" {
		sw.InitLogging(os.Stdout, os.Stdout, os.Stdout)
	} else {
		// if debug env is not present or false, do not log debug output to console
		sw.InitLogging(os.Stdout, ioutil.Discard, os.Stdout)
	}

	sw.Info.Println(fmt.Sprintf("%s%s", "Trips Service Server started on port ", *webServerPort))

	router := sw.NewRouter()

	sw.Fatal.Println(http.ListenAndServe(fmt.Sprintf("%s%s", ":", *webServerPort), router))

}
