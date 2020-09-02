package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestGetResource(t *testing.T) {
	testResource := []struct {
		resource string
		result   string
	}{
		{
			"https://www.baidu.com",
			"",
		},
		{
			"/etc/hosts",
			"",
		},
		{
			"http://file.choerodon.com.cn/choerodon-install/0.21/install.yml",
			"",
		},
		{
			"http://test.yidaqiang.com",
			"",
		},
	}

	for _, r := range testResource {
		content, err := GetResource(r.resource)
		if err != nil {
			log.Error(err)
		}
		str := string(content[:])
		fmt.Println(str)
	}
}
