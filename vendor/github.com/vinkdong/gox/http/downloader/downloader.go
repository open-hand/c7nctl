package downloader

import (
	"github.com/vinkdong/gox/log"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func GetFileContent(url string) (data []byte, statusCode int, err error) {
	log.Infof("getting resource %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		os.Exit(127)
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return data, resp.StatusCode, err
}

func Download(url, path string) (statusCode int, err error) {
	dir := filepath.Dir(path)
	_, err = os.Stat(dir)

	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	data, statusCode, err := GetFileContent(url)
	if err != nil {
		return statusCode, err
	}

	if err := ioutil.WriteFile(path, data, 0755); err == nil {
		log.Infof("downloaded %s", url)
	}
	return statusCode, err
}
