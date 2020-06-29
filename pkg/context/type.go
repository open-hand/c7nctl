package context

import (
	"github.com/choerodon/c7nctl/pkg/config"
	"time"
)

type ChartValue struct {
	Name  string
	Value string
	Input Input
	Case  string
	Check string
}

type PreValue struct {
	Name  string
	Value string
	Check string
	Input Input
}

type PreValueList []*PreValue

type Input struct {
	Enabled  bool
	Regex    string
	Tip      string
	Password bool
	Include  []KV
	Exclude  []KV
	Twice    bool
}

type KV struct {
	Name  string
	Value string
}

// TaskInfo 包含了两类
type TaskInfo struct {
	Name      string
	Namespace string
	RefName   string
	Type      string
	// 资源对象的状态
	Status string
	// 错误原因
	Reason string
	Date   time.Time
	// 保存的配置项
	Values   []ChartValue
	Resource config.Resource

	TaskType string
	Version  string
	Prefix   string
}

type BackendTask struct {
	Name    string
	Success bool
}
