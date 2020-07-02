package client

import (
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/utils"
	"time"
)

/* TaskInfo 用于保存安装过程中的信息，包括 Release 的配置项，以及其他的状态等
   TaskInfo 保存到 k8s cm 中，Type 将 task 分为3类：Release，task，persistent
   TaskInfo 的操作方法都在 K8sClient 中，原因是它依赖于 kubernetes.ClientSet，并且在
   Release，Persistent，InstallDefinition 中都有使用。
   TaskInfo 包含了两类
*/
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

type ChartValue struct {
	Name  string
	Value string
	Input utils.Input
	Case  string
	Check string
}

type PreValue struct {
	Name  string
	Value string
	Check string
	Input utils.Input
}

type PreValueList []*PreValue

type BackendTask struct {
	Name    string
	Success bool
}
