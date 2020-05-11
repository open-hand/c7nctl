package context

import "C"
import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/vinkdong/gox/log"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"sync"
	"time"
)

var Ctx Context

const (
	PvType      = "pv"
	PvcType     = "pvc"
	CRDType     = "crd"
	ReleaseTYPE = "helm"
	TaskType    = "task"

	// Release 安装成功
	SucceedStatus = "succeed"
	// Release 安装失败
	FailedStatus = "failed"
	// Release 输入配置完成
	InputtedStatus = "inputted"
	// Release 未初始化
	UninitializedStatus = "uninitialized"
	// Release 配置宣传成功
	RenderedStatus = "rendered"

	// if have after process while wait
	CreatedStatus      = "created"
	staticLogName      = "c7n-logs"
	staticLogKey       = "logs"
	staticTaskKey      = "tasks"
	staticInstalledKey = "installed"
	staticExecutedKey  = "execute"
	SqlTask            = "sql"
	HttpGetTask        = "httpGet"
)

type Context struct {
	// client for Release
	HelmClient *client.HelmClient
	KubeClient *kubernetes.Interface
	Slaver     *slaver.Slaver

	// configration
	UserConfig   *config.Config
	BackendTasks []*BackendTask
	JobInfo      map[string]JobInfo
	Mux          sync.Mutex
	Metrics      Metrics

	Namespace     string
	CommonLabels  map[string]string
	SlaverAddress string
	SkipInput     bool
	Prefix        string
	RepoUrl       string
}

func (ctx *Context) GetConfig() *config.Config {
	return ctx.UserConfig
}

// 当后台任务不存在时添加并返回 ture，否则之返回 false
func (ctx *Context) AddBackendTask(task *BackendTask) bool {
	for _, v := range ctx.BackendTasks {
		if v.Name == task.Name {
			return false
		}
	}
	ctx.BackendTasks = append(ctx.BackendTasks, task)
	return true
}

func (ctx *Context) AddJobInfo(ji JobInfo) {
	for _, jItem := range ctx.JobInfo {
		if ji.Name == jItem.Name {
			return
		}
	}

	if tmpJob := ctx.JobInfo[ji.Name]; tmpJob.Name != "" {
		log.Infof("JobInfo %s is already save to context\n", ji.Name)
		return
	}
	ctx.JobInfo[ji.Name] = ji
}

// 保存 jobInfo 修改到 cm
func (ctx *Context) UpdateJobInfo(ji JobInfo) {
	ctx.JobInfo[ji.Name] = ji
	// convert map to array
	var jiArray []JobInfo
	for _, ji := range ctx.JobInfo {
		jiArray = append(jiArray, ji)
	}
	// 直接保存所有的 JobInfo
	jiByte, err := yaml.Marshal(jiArray)
	if err != nil {
		log.Error(err)
		// return err
	}

	ctx.saveConfigMapData(string(jiByte), staticLogName, staticLogKey)
}

func (ctx *Context) GetJobInfo(jobName string) JobInfo {
	ctx.checkJobInfo()

	// 不存在返回 nil
	return ctx.JobInfo[jobName]
}

func (ctx *Context) checkJobInfo() {
	if ctx.JobInfo == nil {
		log.Error("Please load JobInfo from kubernetes config map")
	}
}

// TODO remove goto
func (ctx *Context) CheckExist(code int, errMsg ...string) {
	if code == 0 {
		ctx.Metrics.Status = "succeed"
	} else {
		ctx.Metrics.Status = "failed"
	}

	// 等待所有后台任务完成
	for {
		if !ctx.hasBackendTask() {
			log.Info("some backend task not finished yet wait it to be finished")
			break
		} else {
			time.Sleep(time.Second * 10)
		}
	}

	if len(errMsg) > 0 {
		for _, err := range errMsg {
			log.Error(err)
			ctx.Metrics.ErrorMsg = append(ctx.Metrics.ErrorMsg, err)
		}
	}
	ctx.Metrics.Send()
	if code == 0 {
		return
	}
	os.Exit(code)
}

func (ctx *Context) hasBackendTask() bool {
	for _, v := range ctx.BackendTasks {
		if v.Success == false {
			log.Infof("%s has task not finish", v.Name)
			return true
		}
	}
	return false
}

func (ctx *Context) LoadJobInfoFromCM() error {
	insLogs := ctx.GetOrCreateConfigMapData(staticLogName, staticLogKey)

	jil := []JobInfo{}
	if err := yaml.Unmarshal([]byte(insLogs), &jil); err != nil {
		log.Error(err)
	}
	for _, ji := range jil {
		ctx.JobInfo[ji.Name] = ji
	}
	return nil
}

// 将 news 保存到
func (ctx *Context) SaveNews(news *JobInfo) error {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	var key string

	// task 添加到 task key 下，release 添加到configMap的 logs 下
	if news.Type == TaskType {
		key = staticTaskKey
	} else {
		key = staticLogKey
	}
	data := ctx.GetOrCreateConfigMapData(staticLogName, key)

	var jil []JobInfo
	if err := yaml.Unmarshal([]byte(data), &jil); err != nil {
		return err
	}
	news.Date = time.Now()
	if news.RefName == "" {
		news.RefName = news.Name
	}
	jil = append(jil, *news)
	newData, err := yaml.Marshal(jil)
	if err != nil {
		log.Error(err)
		return err
	}

	ctx.saveConfigMapData(string(newData[:]), staticLogName, key)

	if news.Status == SucceedStatus || news.Status == CreatedStatus {
		ctx.SaveSucceed(news)
	}
	return nil
}

func (ctx *Context) UpdateCreated(name, namespace string) error {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	jil := ctx.getSucceedData()
	isUpdate := false
	for k, v := range jil {
		if v.Name == name && v.Namespace == namespace && v.Status == CreatedStatus {
			v.Status = SucceedStatus
			jil[k] = v
			isUpdate = true
		}
	}
	if !isUpdate {
		log.Infof("nothing update with app %s in ns: %s", name, namespace)
	}
	newData, err := yaml.Marshal(jil)
	if err != nil {
		log.Error(err)
		return err
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, staticInstalledKey)
	return nil
}

func (ctx *Context) SaveSucceed(news *JobInfo) error {
	news.Date = time.Now()
	jil := ctx.getSucceedData()
	jil = append(jil, *news)
	newData, err := yaml.Marshal(jil)
	if err != nil {
		log.Error(err)
		return err
	}

	key := staticInstalledKey
	if news.Type == TaskType {
		key = staticExecutedKey
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, key)
	return nil
}

func (ctx *Context) GetSucceed(name string, resourceType string) *JobInfo {
	jil := ctx.getSucceedData()
	for _, v := range jil {
		if v.Name == name && v.Type == resourceType {
			// todo: make sure gc effort
			p := v
			return &p
		}
	}
	return nil
}

func (ctx *Context) GetSucceedTask(taskName, appName string, taskType string) *JobInfo {
	jil := ctx.getSucceedData(staticTaskKey)
	for _, v := range jil {
		if v.Name == taskName && v.RefName == appName && v.TaskType == taskType && v.Status == SucceedStatus {
			// todo: make sure gc effort
			p := v
			return &p
		}
	}
	return nil
}

func (ctx *Context) DeleteSucceedTask(appName string) error {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	jil := ctx.getSucceedData(staticExecutedKey)
	leftNews := make([]JobInfo, 0)
	for _, v := range jil {
		if v.RefName == appName {
			// todo: make sure gc effort
		} else {
			leftNews = append(leftNews, v)
		}
	}

	newData, err := yaml.Marshal(leftNews)
	if err != nil {
		log.Error(err)
		return err
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, staticExecutedKey)
	return nil
}

func (ctx *Context) DeleteSucceed(name, namespace, resourceType string) error {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	jil := ctx.getSucceedData()
	index := -1
	for k, v := range jil {
		if v.Name == name && v.Namespace == namespace && v.Type == resourceType {
			index = k
		}
	}

	if index == -1 {
		log.Infof("nothing delete with app %s in ns: %s", name, namespace)
		return nil
	}
	jil = append(jil[:index], jil[index+1:]...)
	newData, err := yaml.Marshal(jil)
	if err != nil {
		log.Error(err)
		return err
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, staticInstalledKey)
	// todo save delete to log
	return nil
}

func (ctx *Context) getSucceedData(key ...string) []JobInfo {
	cmKey := staticInstalledKey
	if len(key) > 0 {
		cmKey = key[0]
	}
	data := ctx.GetOrCreateConfigMapData(staticLogName, cmKey)
	jil := []JobInfo{}
	_ = yaml.Unmarshal([]byte(data), &jil)
	return jil
}

func (ctx *Context) GetOrCreateConfigMapData(cmName, cmKey string) string {
	cm, err := (*ctx.KubeClient).CoreV1().ConfigMaps(ctx.Namespace).Get(cmName, meta_v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("creating logs to cluster")
			cm = ctx.createInstallLogsCM()
		}
	}
	return cm.Data[cmKey]
}

func (ctx *Context) createInstallLogsCM() *v1.ConfigMap {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	data := make(map[string]string)
	data[staticLogKey] = ""
	data["user_info"] = fmt.Sprintf("email: %s", ctx.Metrics.Mail)
	cm := &v1.ConfigMap{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   staticLogName,
			Labels: ctx.CommonLabels,
		},
		Data: data,
	}
	configMap, err := (*ctx.KubeClient).CoreV1().ConfigMaps(ctx.Namespace).Create(cm)
	if err != nil {
		ctx.CheckExist(128, err.Error())
	}
	return configMap
}

func (ctx *Context) saveConfigMapData(data, cmName, cmKey string) *v1.ConfigMap {
	kubeClient := *ctx.KubeClient
	cm, err := kubeClient.CoreV1().ConfigMaps(ctx.Namespace).Get(cmName, meta_v1.GetOptions{})
	cm.Data[cmKey] = data
	configMap, err := kubeClient.CoreV1().ConfigMaps(ctx.Namespace).Update(cm)
	if err != nil {
		log.Error(err)
		os.Exit(122)
	}
	return configMap
}

func IsNotFound(err error) bool {
	errorStatus, ok := err.(*errors.StatusError)
	if ok && errorStatus.Status().Code == 404 {
		return true
	}
	return false
}

type Exclude struct {
	Start int
	End   int
}
