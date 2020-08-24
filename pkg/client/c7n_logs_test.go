package client

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/utils"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func initKubeClient() *kubernetes.Clientset {
	var defaultKubeconfigPath string
	if home := homeDir(); home != "" {
		defaultKubeconfigPath = filepath.Join(home, ".kube", "config")
	}

	client, err := GetKubeClient(defaultKubeconfigPath)
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func TestSaveAndGetTask(t *testing.T) {
	InitC7nLogs(initKubeClient(), "ydq-test")

	c7nLogsTest := []struct {
		TaskInfo,
		Result TaskInfo
	}{
		{
			TaskInfo: TaskInfo{
				Name:      "task-test01",
				Namespace: "test",
				RefName:   "haha",
				Type:      consts.StaticReleaseKey,
				Status:    consts.SucceedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-01",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
			Result: TaskInfo{
				Name:      "task-test01",
				Namespace: "test",
				RefName:   "haha-test",
				Type:      consts.StaticReleaseKey,
				Status:    consts.SucceedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-01",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
		},
		{
			TaskInfo: TaskInfo{
				Name:      "task-test01",
				Namespace: "test",
				RefName:   "haha",
				Type:      consts.StaticReleaseKey,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-01",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
			Result: TaskInfo{
				Name:      "task-test01",
				Namespace: "test",
				RefName:   "haha-test",
				Type:      consts.StaticReleaseKey,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-01",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
		},
		{
			TaskInfo: TaskInfo{
				Name:      "task-test02",
				Namespace: "test",
				RefName:   "haha",
				Type:      consts.StaticReleaseKey,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-02",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
			Result: TaskInfo{
				Name:      "task-test02",
				Namespace: "test",
				RefName:   "haha-test",
				Type:      consts.StaticTaskKey,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-02",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
		},
		{
			TaskInfo: TaskInfo{
				Name:      "task-test03",
				Namespace: "test",
				RefName:   "haha",
				Type:      consts.StaticPersistentKey,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-03",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
			Result: TaskInfo{
				Name:      "task-test03",
				Namespace: "test",
				RefName:   "haha-test",
				Type:      consts.StaticPersistentKey,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-03",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
		},
		{
			TaskInfo: TaskInfo{
				Name:      "task-test03",
				Namespace: "test",
				RefName:   "haha",
				Type:      consts.SqlTask,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-03",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
			Result: TaskInfo{
				Name:      "task-test03",
				Namespace: "test",
				RefName:   "haha-test",
				Type:      consts.SqlTask,
				Status:    consts.FailedStatus,
				Reason:    "haha",
				Date:      time.Time{},
				Values: []ChartValue{
					{
						Name:  "values-03",
						Value: "hah",
						Input: utils.Input{
							Tip:      "haha",
							Enabled:  false,
							Password: false,
							Include: []utils.KV{
								{
									Name:  "asdfa",
									Value: "asdf",
								},
							},
						},
						Case:  "haha",
						Check: "haha",
					},
				},
				Resource: config.Resource{},
				TaskType: consts.SqlTask,
				Version:  "0.22",
				Prefix:   "asdf",
			},
		},
	}
	for _, c := range c7nLogsTest {
		_, err := SaveTask(c.TaskInfo)
		if err != nil {
			t.Error(err)
		}
		// update
		c.TaskInfo.RefName += "-test"
		task, err := GetTask(c.TaskInfo.Name)
		_, err = SaveTask(c.TaskInfo)
		if err != nil {
			t.Error(err)
		}

		if err != nil {
			t.Error(err)
		}
		if reflect.DeepEqual(c.Result, task) {
			t.Errorf("Taskinfo no equal %+v", task)
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
