package authorize

import (
	"github.com/choerodon/c7nctl/pkg/utils"
	"testing"
)

func TestWrite(t *testing.T) {
	auth := Authorization{
		Username:    "vink",
		Token:       "d97448df-573d-47db-8b8b-358493ca0c38",
		ClusterName: "for-test",
		ServerUrl:   "https://api.vk.vu",
		Config:      &utils.Config{},
	}
	auth.Write()
}

func TestDefaultAuthorization(t *testing.T) {
	DefaultAuthorization()
}
