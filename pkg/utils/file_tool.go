package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func DirectoryCheck(dirName string) {
	_, err := os.Stat(dirName)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			log.Error(err)
		}
	}
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
