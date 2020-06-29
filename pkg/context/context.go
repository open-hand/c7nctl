package context

import "C"
import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/slaver"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"sync"
)

var Ctx Context

const (
	PvType      = "pv"
	PvcType     = "pvc"
	CRDType     = "crd"
	ReleaseType = "helm"
	TaskType    = "task"

	// Release 未初始化
	UninitializedStatus = "uninitialized"
	// Release 输入配置完成
	InputtedStatus = "inputted"
	// Release 配置渲染成功
	RenderedStatus = "rendered"
	// Release 安装成功
	InstalledStatus = "installed"
	// Release 完成所有安装步骤，即 afterTask 完成
	SucceedStatus = "succeed"
	// Release 安装失败
	FailedStatus = "failed"

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

	// TODO 可以去掉，使用单例模式
	HelmClient *action.Install
	KubeClient *kubernetes.Interface
	Slaver     *slaver.Slaver

	// configration
	UserConfig *config.C7nConfig
	JobInfo    []*TaskInfo
	Metrics    Metrics

	Mux sync.Mutex

	Namespace     string
	CommonLabels  map[string]string
	SlaverAddress string
	SkipInput     bool
	Prefix        string
	RepoUrl       string
	Version       string
}

func (ctx *Context) GetConfig() *config.C7nConfig {
	return ctx.UserConfig
}

func (ctx *Context) AddJobInfo(ji *TaskInfo) {
	_, r := ctx.GetJobInfo(ji.Name)
	if r.Name != "" {
		log.WithField("release", ji.Name).Info("Release already existed")
		return
	}
	ctx.JobInfo = append(ctx.JobInfo, ji)
	ctx.saveJobInfo()
}

// 保存 jobInfo 修改到 cm
func (ctx *Context) UpdateJobInfo(ji *TaskInfo) {
	idx, _ := ctx.GetJobInfo(ji.Name)
	if idx > 0 {
		ctx.JobInfo[idx] = ji
	} else {
		ctx.JobInfo = append(ctx.JobInfo, ji)
	}

	ctx.saveJobInfo()
}

func (ctx *Context) GetJobInfo(jobName string) (int, *TaskInfo) {
	if ctx.JobInfo == nil {
		log.Error("Release job info can't be empty.")
		os.Exit(128)
	}
	for idx, r := range ctx.JobInfo {
		if r.Name == jobName {
			return idx, r
		}
	}
	// 不存在返回 nil
	return -1, &TaskInfo{}
}

func (ctx *Context) LoadJobInfo() error {
	insLogs := ctx.GetOrCreateConfigMapData(staticLogName, staticLogKey)

	jil := []*TaskInfo{}
	if err := yaml.Unmarshal([]byte(insLogs), &jil); err != nil {
		log.Error(err)
	}
	ctx.JobInfo = jil
	return nil
}

func (ctx *Context) CheckExist(code int, errMsg ...string) {
	if code == 0 {
		ctx.Metrics.Status = "succeed"
	} else {
		ctx.Metrics.Status = "failed"
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

func (ctx *Context) SendMetrics(err error) {
	if err == nil {
		ctx.Metrics.Status = "succeed"
	} else {
		ctx.Metrics.Status = "failed"
		ctx.Metrics.ErrorMsg = append(ctx.Metrics.ErrorMsg, err.Error())
	}

	ctx.Metrics.Send()
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

func (ctx *Context) saveJobInfo() {
	// 直接保存所有的 TaskInfo
	jiByte, err := yaml.Marshal(ctx.JobInfo)
	if err != nil {
		log.Error(err)
	}
	ctx.saveConfigMapData(string(jiByte), staticLogName, staticLogKey)
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
		log.Error(err)
		os.Exit(127)
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
