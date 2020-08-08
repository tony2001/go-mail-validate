package main

import (
	"flag"
	"fmt"
	syslog "log"
	"net/http"
	"os"

	"go-mail-validate/config"
	"go-mail-validate/log"
)

const _defaultPort = "8080"

var smtpTimeoutMsec int
var debugMode bool

type endPointT struct {
	url         string
	handlerFunc func(w http.ResponseWriter, r *http.Request)
}

type validatorWeightT struct {
	name   string
	weight uint32
}

//supported endpoints list
var supportedEndPoints = []endPointT{
	{"/", handleDefault},
	{"/email/validate", handleEmailValidate},
}

func main() {

	flag.BoolVar(&debugMode, "d", false, "enable debug logging")
	flag.IntVar(&smtpTimeoutMsec, "s", 1000, "SMTP timeout in milliseconds")
	flag.Parse()

	//check env var PORT
	portStr := os.Getenv("PORT")
	if portStr == "" {
		//or use the default value
		portStr = _defaultPort
	}

	var err error
	var logLevel = log.Info
	if debugMode {
		logLevel = log.Debug
	}

	err = config.NewConfig(portStr, config.LogLevel(logLevel), config.SmtpTimeout(smtpTimeoutMsec))
	if err != nil {
		syslog.Fatalf("failed to create config: %s", err)
	}

	log.NewLogger(config.GetLogLevel())

	for _, point := range supportedEndPoints {
		http.HandleFunc(point.url, point.handlerFunc)
	}

	log.Debugf("starting server with config: %s", config.String())

	listenStr := fmt.Sprintf(":%d", config.GetPort())
	log.Fatalf("%s", http.ListenAndServe(listenStr, nil))

}
