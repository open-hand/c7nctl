package client

import (
	"github.com/choerodon/c7nctl/pkg/common/consts"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"testing"
	"time"
)

func TestK8sClient_SaveTaskInfoToCM(t *testing.T) {
	home, _ := homedir.Dir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	client, _ := GetKubeClient(kubeconfig)

	localClient := NewK8sClient(client)

	task := TaskInfo{
		Name:      "test",
		Namespace: "default",
		RefName:   "qwer",
		Type:      consts.TaskType,
		Status:    consts.SucceedStatus,
		Reason:    "",
		Date:      time.Time{},
		Values:    nil,
		Resource:  config.Resource{},
		TaskType:  consts.PvType,
		Version:   "",
		Prefix:    "",
	}
	_ = localClient.SaveTaskInfoToCM("test", task)
}
