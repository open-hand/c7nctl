package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

func CheckErrAndExit(err error, exitCode int) {
	if err != nil {
		log.Error(err)
		os.Exit(exitCode)
	}
}

func CheckErrAndSendMetrics(err error) {

}

func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
