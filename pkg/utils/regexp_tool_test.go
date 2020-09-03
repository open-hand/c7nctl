package utils

import (
	"testing"
)

func TestCheckDomain(t *testing.T) {
	domainTest := []struct {
		Domain  string
		Resoult bool
	}{
		{
			"www.baidu.com",
			true,
		},
		{
			"devops.dev.yidaqiang.com",
			true,
		},
		{
			"www.bai_du.com",
			false,
		},
	}
	for _, d := range domainTest {
		if CheckDomain(d.Domain) != d.Resoult {
			t.Error("domain check failed")
		}
	}
}
