package client

import (
	"context"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	c7nerrors "github.com/choerodon/c7nctl/pkg/common/errors"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/ghodss/yaml"
	stderrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
	"time"
)

type C7nLogs struct {
	client    *kubernetes.Clientset
	Name      string
	namespace string
	// 避免 yaml.Unmarshal 无法取地址
	Tasks map[string]*[]TaskInfo
}

/*
TaskInfo 用于保存安装过程中的信息，包括 Release 的配置项，以及其他的状态等

	保存到 k8s cm 中，Type 将 task 分为3类：release，task，persistent
*/
type TaskInfo struct {
	// 唯一值
	Name      string
	Namespace string
	RefName   string
	// 任务类型： Release, task, persistence
	Type string
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

var c7nLogs C7nLogs

func InitC7nLogs(client *kubernetes.Clientset, namespace string) {
	var once sync.Once
	once.Do(func() {
		if c7nLogs.client == nil {
			c7nLogs.client = client
			c7nLogs.namespace = namespace
			c7nLogs.Name = consts.StaticLogsCM
			c7nLogs.Tasks = map[string]*[]TaskInfo{}
		}
	})
}

func NewReleaseTask(release, namespace, version, prefix string) *TaskInfo {
	return &TaskInfo{
		Name:      release,
		Namespace: namespace,
		Type:      consts.StaticReleaseKey,
		Status:    consts.UninitializedStatus,
		Date:      time.Now(),
		Version:   version,
		Prefix:    prefix,
	}
}

func NewReleaseJobTask(rlsJob, taskType, version string) *TaskInfo {
	return &TaskInfo{
		Name:      rlsJob,
		Namespace: c7nLogs.namespace,
		Type:      consts.StaticTaskKey,
		Status:    consts.UninitializedStatus,
		Date:      time.Now(),
		TaskType:  taskType,
		Version:   version,
	}
}

func GetTask(task string) (*TaskInfo, error) {
	if err := getC7nLogs(c7nLogs.namespace, c7nLogs.Name); err != nil {
		panic(err)
	}
	for group := range c7nLogs.Tasks {
		tasks := *c7nLogs.Tasks[group]
		for idx, t := range tasks {
			if t.Name == task {
				return &tasks[idx], nil
			}
		}
		log.Debugf("Task %s isn't in group %s", task, group)
	}
	return nil, stderrors.WithMessage(c7nerrors.TaskInfoIsNotFoundError, fmt.Sprintf("Task %s isn't in configMaps c7n-logs", task))
}

func SaveTask(t TaskInfo) (*TaskInfo, error) {
	if t.Name != "" {
		task, err := GetTask(t.Name)
		// 错误为不存在时，任务追加到末尾。
		if err != nil {
			if stderrors.Is(err, c7nerrors.TaskInfoIsNotFoundError) {
				log.Debugf("Task %s isn't in c7n-logs,new add it", t.Name)
				*c7nLogs.Tasks[t.Type] = append(*c7nLogs.Tasks[t.Type], t)
			} else {
				return nil, stderrors.WithMessage(err, "Getting task failed when save Task: ")
			}
		} else {
			*task = t
			log.Debugf("Update task is %+v", *task)
		}
	} else {
		log.Debug("Task is empty，Please confirm that the task exists")
	}

	if err := saveC7nLogs(); err != nil {
		return nil, err
	}
	return &t, nil
}

func saveC7nLogs() error {
	cm, err := getConfigMaps(c7nLogs.namespace, c7nLogs.Name)
	if err != nil {
		return stderrors.WithMessage(err, "Save configMaps c7n-logs failed: ")
	}

	for key := range c7nLogs.Tasks {
		tbyte, err := yaml.Marshal(*c7nLogs.Tasks[key])
		if err != nil {
			log.Error(err)
		}
		if cm.Data == nil {
			cm.Data = map[string]string{}
		}
		cm.Data[key] = string(tbyte)
		log.Debugf("Saving ConfigMaps key %s", key)
	}
	if _, err = c7nLogs.client.CoreV1().ConfigMaps(c7nLogs.namespace).Update(context.Background(), cm,
		metav1.UpdateOptions{}); err != nil {
		return stderrors.WithMessage(err, "Save configMaps c7n-logs failed: ")
	}
	return nil
}

// 如果不存在 c7n-logs 就创建，之后将其数据初始化到 c7nLogs
func getC7nLogs(namespace, cmName string) error {
	cm, err := getConfigMaps(namespace, cmName)
	if err != nil {
		return stderrors.WithMessage(err, "Get configMaps c7n-logs failed: ")
	}

	for key := range cm.Data {
		if c7nLogs.Tasks[key] == nil {
			log.Debugf("Task type %s is empty", key)
			c7nLogs.Tasks[key] = new([]TaskInfo)
		}

		if err = yaml.Unmarshal([]byte(cm.Data[key]), c7nLogs.Tasks[key]); err != nil {
			panic(err)
		}
		log.Debugf("Getting ConfigMaps key %s", key)
	}
	return nil
}

func getConfigMaps(namespace, cmName string) (cm *v1.ConfigMap, err error) {
	cm, err = c7nLogs.client.CoreV1().ConfigMaps(c7nLogs.namespace).Get(context.Background(), c7nLogs.Name, metav1.GetOptions{})
	// 如果不存在则创建 cm
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Infof("Config map %s isn't existing, now Create it", cmName)

			cm = &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cmName,
					Namespace: namespace,
					Labels:    consts.CommonLabels,
				},
				Data: map[string]string{
					consts.StaticReleaseKey:    "",
					consts.StaticPersistentKey: "",
					consts.StaticTaskKey:       "",
				},
			}
			cm, err = c7nLogs.client.CoreV1().ConfigMaps(c7nLogs.namespace).Create(context.Background(), cm,
				metav1.CreateOptions{})
			if err != nil {
				return nil, stderrors.WithMessage(err, fmt.Sprintf("Failed to create configMaps %s in namespace %s",
					c7nLogs.Name, c7nLogs.namespace))
			}
			log.Infof("Successfully created ConfigMaps %s in namespace %s", c7nLogs.Name, c7nLogs.namespace)
		} else {
			return nil, stderrors.WithMessage(err, fmt.Sprintf("Failed to get configMaps %s in namespace %s",
				c7nLogs.Name, c7nLogs.namespace))
		}
	}
	log.Debugf("Using existing configMaps %s in namespace %s", c7nLogs.Name, c7nLogs.namespace)
	return cm, nil
}
