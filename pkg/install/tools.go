package install

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/choerodon/c7nctl/pkg/utils"
	"github.com/vinkdong/gox/log"
	"github.com/vinkdong/gox/random"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"os"
	"sync"
	"time"
)

var Ctx Context

const (
	PvType        = "pv"
	PvcType       = "pvc"
	CRDType       = "crd"
	ReleaseTYPE   = "helm"
	TaskType      = "task"
	SucceedStatus = "succeed"
	FailedStatus  = "failed"
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
	Client        kubernetes.Interface
	Namespace     string
	CommonLabels  map[string]string
	SlaverAddress string
	Slaver        *slaver.Slaver
	UserConfig    *config.Config
	BackendTasks  []*BackendTask
	Mux           sync.Mutex
	Metrics       utils.Metrics
}

type BackendTask struct {
	Name    string
	Success bool
}

// i want use log but it make ...
type News struct {
	Name      string
	Namespace string
	RefName   string
	Type      string
	Status    string
	Reason    string
	Date      time.Time
	Resource  config.Resource
	Values    []ChartValue
	PreValue  PreValueList
	TaskType  string
	Version   string
	Prefix    string
}

type NewsResourceList struct {
	News []News `yaml:"logs"`
}

func (ctx *Context) AddBackendTask(task *BackendTask) bool {
	for _, v := range ctx.BackendTasks {
		if v.Name == task.Name {
			return false
		}
	}
	ctx.BackendTasks = append(ctx.BackendTasks, task)
	return true
}

func (ctx *Context) CheckExist(code int, errMsg ...string) {
	if code == 0 {
		ctx.Metrics.Status = "succeed"
	} else {
		ctx.Metrics.Status = "failed"
	}
	if !ctx.HasBackendTask() {
		goto exit
	}
	log.Info("some backend task not finished yet wait it to be finished")
	for {
		select {
		case <-time.Tick(time.Second * 1):
			if !ctx.HasBackendTask() {
				goto exit
			}
		}
	}
exit:
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

func (ctx *Context) HasBackendTask() bool {
	for _, v := range ctx.BackendTasks {
		if v.Success == false {
			return true
		}
	}
	return false
}

func (ctx *Context) SaveNews(news *News) error {
	ctx.Mux.Lock()
	defer ctx.Mux.Unlock()
	var key string
	if news.Type == TaskType {
		key = staticTaskKey
	} else {
		key = staticLogKey
	}
	data := ctx.GetOrCreateConfigMapData(staticLogName, key)
	nr := &NewsResourceList{}
	yaml.Unmarshal([]byte(data), nr)
	news.Date = time.Now()
	if news.RefName == "" {
		news.RefName = news.Name
	}
	nr.News = append(nr.News, *news)
	newData, err := yaml.Marshal(nr)
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
	nr := ctx.getSucceedData()
	isUpdate := false
	for k, v := range nr.News {
		if v.Name == name && v.Namespace == namespace && v.Status == CreatedStatus {
			v.Status = SucceedStatus
			nr.News[k] = v
			isUpdate = true
		}
	}
	if !isUpdate {
		log.Infof("nothing update with app %s in ns: %s", name, namespace)
	}
	newData, err := yaml.Marshal(nr)
	if err != nil {
		log.Error(err)
		return err
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, staticInstalledKey)
	return nil
}

func (ctx *Context) SaveSucceed(news *News) error {
	news.Date = time.Now()
	nr := ctx.getSucceedData()
	nr.News = append(nr.News, *news)
	newData, err := yaml.Marshal(nr)
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

func (ctx *Context) GetSucceed(name string, resourceType string) *News {
	nr := ctx.getSucceedData()
	for _, v := range nr.News {
		if v.Name == name && v.Type == resourceType {
			// todo: make sure gc effort
			p := v
			return &p
		}
	}
	return nil
}

func (ctx *Context) GetSucceedTask(taskName, appName string, taskType string) *News {
	nr := ctx.getSucceedData(staticExecutedKey)
	for _, v := range nr.News {
		if v.Name == taskName && v.RefName == appName && v.TaskType == taskType {
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
	nr := ctx.getSucceedData(staticExecutedKey)
	leftNews := make([]News, 0)
	for _, v := range nr.News {
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
	nr := ctx.getSucceedData()
	index := -1
	for k, v := range nr.News {
		if v.Name == name && v.Namespace == namespace && v.Type == resourceType {
			index = k
		}
	}

	if index == -1 {
		log.Infof("nothing delete with app %s in ns: %s", name, namespace)
		return nil
	}
	nr.News = append(nr.News[:index], nr.News[index+1:]...)
	newData, err := yaml.Marshal(nr)
	if err != nil {
		log.Error(err)
		return err
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, staticInstalledKey)
	// todo save delete to log
	return nil
}

func (ctx *Context) getSucceedData(key ...string) *NewsResourceList {
	cmKey := staticInstalledKey
	if len(key) > 0 {
		cmKey = key[0]
	}
	data := ctx.GetOrCreateConfigMapData(staticLogName, cmKey)
	nr := &NewsResourceList{}
	yaml.Unmarshal([]byte(data), nr)
	return nr
}

func (ctx *Context) GetOrCreateConfigMapData(cmName, cmKey string) string {
	if ctx.Client == nil {
		log.Error("Get k8s client failed")
		os.Exit(127)
	}
	cm, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Get(cmName, meta_v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("creating logs to cluster")
			cm = ctx.createNewsData()
		}
	}
	return cm.Data[cmKey]
}

func (ctx *Context) createNewsData() *v1.ConfigMap {
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
	configMap, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Create(cm)
	if err != nil {
		ctx.CheckExist(128, err.Error())
	}
	return configMap
}

func (ctx *Context) saveConfigMapData(data, cmName, cmKey string) *v1.ConfigMap {

	cm, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Get(cmName, meta_v1.GetOptions{})
	cm.Data[cmKey] = data
	configMap, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Update(cm)
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

func RandomInt(min, max int, exclude ...Exclude) {
	randInt := min + rand.Intn(max)
	for _, e := range exclude {
		if randInt >= e.Start && randInt <= e.End {
			randInt += 1
		}
	}
}

func RandomToken(length int) string {
	bytes := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		random.Seed(time.Now().UnixNano())
		op := random.RangeIntInclude(random.Slice{Start: 48, End: 57},
			random.Slice{Start: 65, End: 90}, random.Slice{Start: 97, End: 122})
		bytes[i] = byte(op) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func GenerateRunnerToken(length int) string {
	bytes := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		random.Seed(time.Now().UnixNano())
		op := random.RangeIntInclude(random.Slice{Start: 48, End: 57},
			random.Slice{Start: 97, End: 122})
		bytes[i] = byte(op) //A=65 and Z = 65+25
	}
	return string(bytes)
}
