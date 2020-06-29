package utils

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func GetRemoteResource(resourceUrl string) (data []byte, err error) {
	log.WithField("url", resourceUrl).Infof("getting resource")

	resp, err := http.Get(resourceUrl)
	if err != nil {
		return nil, errors.WithMessage(err, "Unknown error when get remote resource")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Network error when get remote resource.")
	}

	return ioutil.ReadAll(resp.Body)
}
