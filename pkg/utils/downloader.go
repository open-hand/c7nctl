package utils

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetRemoteResource(resourceUrl string) (data []byte) {
	log.WithField("url", resourceUrl).Infof("getting resource")

	resp, err := http.Get(resourceUrl)
	CheckErrAndExit(err, 127)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		CheckErrAndExit(errors.New("Network error when get remote resource."), 127)
	}
	data, err = ioutil.ReadAll(resp.Body)
	CheckErrAndExit(err, 127)
	return data
}
