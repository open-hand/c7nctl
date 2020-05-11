package utils

import (
	"github.com/vinkdong/gox/log"
	"io/ioutil"
	"net/http"
	"os"
)

func GetRemoteResource(resourceUrl string) (data []byte) {
	log.Infof("getting resource %s", resourceUrl)

	resp, err := http.Get(resourceUrl)
	if err != nil {
		log.Error(err)
		os.Exit(127)
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Get resource %s failed", resourceUrl)
		log.Error(err)
		os.Exit(127)
	}
	return data
}
