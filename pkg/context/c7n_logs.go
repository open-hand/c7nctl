package context

import (
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// 所有的 taskName 都唯一
type c7nLogs struct {
	Release []TaskInfo
	Task    []TaskInfo
	Pvc     []TaskInfo
}

func GetReleaseTaskInfo(namespace, taskName string) (*TaskInfo, error) {
	task, err := GetTaskFromCM(namespace, taskName)
	if err != nil {
		log.Warn(err)
		task = &TaskInfo{
			Name:      taskName,
			Namespace: namespace,
			Type:      ReleaseType,
			Status:    UninitializedStatus,
			Date:      time.Now(),
		}
		err = AddTaskToCM(namespace, *task)
	}
	return task, err
}

func GetTaskFromCM(namespace, taskName string) (*TaskInfo, error) {
	cl, err := getC7nLogsCM(namespace)
	if err != nil {
		return nil, err
	}
	task := getTaskFromC7nLos(taskName, cl)
	if task == nil {
		return nil, errors.New("Task isn't existing")
	}
	return task, nil
}

func AddTaskToCM(namespace string, task TaskInfo) error {
	if IsTaskExisting(namespace, task.Name, task.Type) {
		return errors.New("Task already exists，skip the add.")
	}

	cl, err := getC7nLogsCM(namespace)
	if err != nil {
		return err
	}

	switch task.Type {
	case PvType, PvcType:
		{
			cl.Pvc = append(cl.Pvc, task)
		}
	case TaskType:
		{
			cl.Task = append(cl.Task, task)
		}
	case ReleaseType:
		{
			cl.Release = append(cl.Release, task)
		}
	}
	return saveC7nLogsCm(namespace, cl)
}

func UpdateTaskToCM(namespace string, task TaskInfo) error {
	if !IsTaskExisting(namespace, task.Name, task.Type) {
		return errors.New("Task isn't existing")
	}
	cl, err := getC7nLogsCM(namespace)
	if err != nil {
		return err
	}
	//TODO
	t := getTaskFromC7nLos(task.Name, cl)
	t = &task
	return saveC7nLogsCm(namespace, cl)
}

func UpdateTaskStatusToCM(namespace, taskName, status string) error {
	cl, err := getC7nLogsCM(namespace)
	if err != nil {
		return err
	}
	task := getTaskFromC7nLos(taskName, cl)
	if task == nil {
		return errors.New("Task isn't existing")
	}
	task.Status = status
	return saveC7nLogsCm(namespace, cl)
}

func IsTaskExisting(namespace, taskName, t string) bool {
	cl, err := getC7nLogsCM(namespace)
	if err != nil {
		log.Error(err)
		return false
	}
	switch t {
	case PvcType, PvType:
		{
			if task := getTaskFromArray(taskName, cl.Pvc); task != nil {
				return true
			}
		}
	case TaskType:
		{
			if task := getTaskFromArray(taskName, cl.Task); task != nil {
				return true
			}
		}
	case ReleaseType:
		{
			if task := getTaskFromArray(taskName, cl.Release); task != nil {
				return true
			}
		}
	}
	return false
}

func C7nlogsFunc(task TaskInfo, f func(info TaskInfo)) {
	/*
		switch task.Type {
		case PvcType, PvType: {
			f(task)
		}
		case TaskType: {
			if task := getTaskFromArray(taskName, cl.Task); task != nil {
				return true
			}
		}
		case ReleaseType: {
			if task := getTaskFromArray(taskName, cl.Release); task != nil {
				return true
			}
		}
		}*/
}

func getTaskFromArray(task string, tasks []TaskInfo) *TaskInfo {
	for _, t := range tasks {
		if t.Name == task {
			return &t
		}
	}
	return nil
}

func getTaskFromC7nLos(task string, logs *c7nLogs) *TaskInfo {
	if task := getTaskFromArray(task, logs.Release); task != nil {
		return task
	}
	if task := getTaskFromArray(task, logs.Pvc); task != nil {
		return task
	}
	if task := getTaskFromArray(task, logs.Task); task != nil {
		return task
	}
	return nil
}

func getC7nLogsCM(namespace string) (*c7nLogs, error) {
	cm, err := c7nclient.GetOrCreateCM(namespace, staticLogName)
	if err != nil {
		return nil, err
	}
	data := cm.Data[staticLogKey]
	var cl c7nLogs

	if err = yaml.Unmarshal([]byte(data), &cl); err != nil {
		return nil, err
	}
	return &cl, nil
}

func saveC7nLogsCm(namespace string, logs *c7nLogs) error {
	taskData, err := yaml.Marshal(logs)
	if err != nil {
		return err
	}
	cmData := map[string]string{
		staticLogName: string(taskData),
	}
	_, err = c7nclient.SaveToCM(namespace, staticLogName, cmData)
	if err != nil {
		err = errors.WithMessage(err, "Failed to save C7nLogs to configMaps")
	}
	return err
}
