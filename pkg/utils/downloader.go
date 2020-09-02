package utils

import (
	"fmt"
	stderrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetRemoteResource(resourceUrl string) (data []byte, err error) {

	log.WithField("url", resourceUrl).Infof("getting resource")

	resp, err := http.Get(resourceUrl)
	if err != nil {
		return nil, stderrors.WithMessage(err, "Unknown error when get remote resource")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, stderrors.New("Network error when get remote resource.")
	}

	return ioutil.ReadAll(resp.Body)
}

// 获取本地文件或者网络文件
func GetResource(resource string) (data []byte, err error) {
	log.Infof("Start retrieving the resource file %s", resource)
	if err, ok := IsFileExist(resource); ok {
		log.Debugf("Read Local file %s", resource)
		data, err = ioutil.ReadFile(resource)
		if err != nil {
			return nil, stderrors.WithMessage(err, fmt.Sprintf("Failed to Read %s", resource))
		}
		return data, err
	} else if err != nil {
		log.Debugf("can't find file %s : %+v", resource, err)
	}
	if _, err := url.Parse(resource); err != nil {
		return nil, stderrors.WithMessage(err, "parse url %s failed: ")
	}

	log.Debugf("Read remote file %s", resource)
	resp, err := http.Get(resource)
	if err != nil {
		return nil, stderrors.WithMessage(err, "Unknown error when get remote resource")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, stderrors.New("Network error when get remote resource.")
	}

	return ioutil.ReadAll(resp.Body)
}
