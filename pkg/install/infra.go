package install

import (
	"bytes"
	"fmt"
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/helm"
	pb "github.com/choerodon/c7n/pkg/protobuf"
	"github.com/choerodon/c7n/pkg/slaver"
	"github.com/pkg/errors"
	"github.com/vinkdong/gox/log"
	"k8s.io/client-go/kubernetes"
	"os"
	"text/template"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	core_v1 "k8s.io/api/core/v1"
	"time"
)

func (infra *InfraResource) executePreCommands() error {
	err := infra.executeExternalFunc(infra.PreInstall)
	return err
}

func (infra *InfraResource) executeExternalFunc(c []PreInstall) error {
	for _, pi := range c {
		if len(pi.Commands) > 0 {
			if err := pi.ExecuteCommands(infra); err != nil {
				return err
			}
		}
		if pi.Request != nil {
			if err := pi.ExecuteRequests(infra); err != nil {
				return err
			}
		}
	}
	return nil
}

func (infra *InfraResource) preparePersistence(client kubernetes.Interface, config *config.Config, commonLabel map[string]string) error {
	getPvs := config.Spec.Persistence.GetPersistentVolumeSource
	namespace := config.Metadata.Namespace
	commonLabel["app"] = infra.Name
	for _, persistence := range infra.Persistence {
		persistence.Client = client

		persistence.CommonLabels = commonLabel
		persistence.CommonLabels["pv"] = persistence.Name

		// check or create dir
		dir := slaver.Dir{
			Mode: persistence.Mode,
			Path: persistence.Path,
			Own:  persistence.Own,
		}
		if err := Ctx.Slaver.MakeDir(dir); err != nil {
			return err
		}

		if err := persistence.CheckOrCreatePv(getPvs(persistence.Path)); err != nil {
			return err
		}
		if persistence.PvcEnabled {
			persistence.Namespace = namespace
			err := persistence.CheckOrCreatePvc()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (infra *InfraResource) applyUserResource() error {
	r := Ctx.UserConfig.GetResource(infra.Name)
	if r == nil {
		log.Infof("no user config resource for %s", infra.Name)
		return nil
	}
	if r.External {
		infra.Resource = r
		return nil
	}
	// just override domain
	if r.Domain != "" {
		infra.Resource.Domain = r.Domain
	}

	return nil
}

func (infra *InfraResource) executePreValues() error {
	return infra.PreValues.prepareValues()
}

func (infra *InfraResource) renderValue(tplString string) string {
	tpl, err := template.New(infra.Name).Parse(tplString)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	var data bytes.Buffer
	err = tpl.Execute(&data, infra)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	return data.String()
}

func (infra *InfraResource) GetPods() (*core_v1.PodList, error)  {
	selectLabel := make(map[string]string)
	selectLabel["choerodon.io/release"] = infra.Name
	set := labels.Set(selectLabel)
	opts := v1.ListOptions{
		LabelSelector: set.AsSelector().String(),
	}
	return Ctx.Client.CoreV1().Pods(infra.Namespace).List(opts)
}

func (infra *InfraResource) GetPodIp() string  {
reget:
	poList, err := infra.GetPods()
	if err != nil || len(poList.Items) < 1 {
		log.Errorf("can't get a pod from %s, retry...",infra.Name)
		goto reget
	}
	for _, po := range poList.Items {
		if po.Status.Phase == core_v1.PodRunning {
			return po.Status.PodIP
		}
	}
	log.Errorf("can't get a running pod from %s, retry...",infra.Name)
	goto reget
}

func (infra *InfraResource) GetPreValue(key string) string {
	return infra.PreValues.getValues(key)
}

func (infra *InfraResource) GetRequire(app string) *InfraResource {
	news := Ctx.GetSucceed(app, ReleaseTYPE)
	i := &InfraResource{
		Name:      app,
		Namespace: infra.Namespace,
		Client:    infra.Client,
		Home:      infra.Home,
		PreValues: news.PreValue,
		Values:    news.Values,
	}
	return i
}

func (infra *InfraResource) GetRequireResource(app string) config.Resource {
	res := Ctx.UserConfig.Spec.Resources
	if r, ok := res[app]; ok {
		return *r
	}
	news := Ctx.GetSucceed(app, ReleaseTYPE)
	if news == nil {
		log.Errorf("require [%s] not right installed or defined", app)
		os.Exit(121)
	}
	return news.Resource
}

func (infra *InfraResource) GetRequirePreValue(app string) config.Resource {
	res := Ctx.UserConfig.Spec.Resources
	if r, ok := res[app]; ok {
		return *r
	}
	news := Ctx.GetSucceed(app, ReleaseTYPE)
	if news == nil {
		log.Errorf("require [%s] not right installed or defined", app)
		os.Exit(121)
	}
	return news.Resource
}

// convert yml values to values list as xxx=yyy
func (infra *InfraResource) HelmValues() ([]string, []ChartValue) {
	values := make([]string, len(infra.Values))
	// store values for feature use
	cvList := make([]ChartValue, len(infra.Values))
	for k, v := range infra.Values {
		value := ""
		if v.Input.Enabled {
			password, err := AcceptUserPassword(v.Input)
			if err != nil {
				log.Error(err)
				os.Exit(128)
			}
			value = password
		} else {
			value = infra.renderValue(v.Value)
		}
		values[k] = fmt.Sprintf("%s=%s", v.Name, value)
		v.Value = value
		cvList[k] = v
	}
	// todo: no return cvList ?
	infra.Values = cvList
	return values, cvList
}

func (infra *InfraResource) GetValue(key string) string {
	for _, v := range infra.Values {
		if v.Name == key {
			return v.Value
		}
	}
	log.Infof("can't get value '%s' of %s", key, infra.Name)
	return ""
}

// only used for save log
func (infra *InfraResource) renderResource() config.Resource {
	//todo: just render password now, add more
	r := infra.Resource
	tpl, err := template.New(fmt.Sprintf("r-%s-%s", infra.Name, "password")).Parse(r.Password)
	if err != nil {
		log.Info(err)
		os.Exit(125)
	}
	var data bytes.Buffer
	if err := tpl.Execute(&data, infra); err != nil {
		log.Error(err)
		os.Exit(125)
	}
	r.Password = data.String()
	if r.Password != "" {
		log.Debugf("%s: resource password is %s",infra.Name,r.Password)
	}
	r.Url = infra.renderValue(r.Url)
	if r.Url != "" {
		log.Debugf("%s: resource url is %s",infra.Name,r.Url)
	}
	return *r
}

// install infra
func (infra *InfraResource) Install(ch chan error) error {
	values, cvList := infra.HelmValues()
	chartArgs := helm.ChartArgs{
		ReleaseName: infra.Name,
		Namespace:   infra.Namespace,
		RepoUrl:     infra.RepoURL,
		Verify:      false,
		Version:     infra.Version,
		ChartName:   infra.Chart,
	}
	log.Infof("installing %s", infra.Name)
	err := infra.Client.InstallRelease(values, chartArgs)

	news := &News{
		Name:      infra.Name,
		Namespace: infra.Namespace,
		RefName:   infra.Name,
		Status:    FailedStatus,
		Type:      ReleaseTYPE,
		Resource:  infra.renderResource(),
		Values:    cvList,
		PreValue:  infra.PreValues,
	}
	defer Ctx.SaveNews(news)

	if err != nil {
		news.Reason = err.Error()
		ch <- err
		return err
	}

	if len(infra.AfterInstall) > 0 {
		news.Status = CreatedStatus
		infra.CheckExecuteAfterTasks()
	} else {
		news.Status = SucceedStatus
	}
	ch <- nil
	return nil
}

func (infra *InfraResource) CheckExecuteAfterTasks() error {
	task := &BackendTask{
		Success: false,
		Name:    infra.Name,
	}
	Ctx.AddBackendTask(task)
	go infra.executeAfterTasks(task)
	return nil
}

func (infra *InfraResource) executeAfterTasks(task *BackendTask) error {
	err := infra.CheckRunning()
	if err != nil {
		log.Error(err)
	}
	log.Successf("%s: started, will execute required commands and requests", infra.Name)
	err = infra.executeExternalFunc(infra.AfterInstall)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := Ctx.UpdateCreated(infra.Name, infra.Namespace); err == nil {
		task.Success = true
	}
	return nil
}

// get server definition
func (infra *InfraResource) GetInfra(key string) *InfraResource {
	spec := infra.Home.Spec
	if app := infra.getAppFromList(key, spec.Infra); app != nil {
		return app
	}
	if app := infra.getAppFromList(key, spec.Framework); app != nil {
		return app
	}
	if app := infra.getAppFromList(key, spec.DevOps); app != nil {
		return app
	}
	if app := infra.getAppFromList(key, spec.Agile); app != nil {
		return app
	}
	if app := infra.getAppFromList(key, spec.TestManager); app != nil {
		return app
	}
	return nil
}

func (infra *InfraResource) getAppFromList(appName string, resourceList []*InfraResource) *InfraResource {
	for _, v := range resourceList {
		if v.Name == appName {
			return v
		}
	}
	return nil
}

// just search the key
func (infra *InfraResource) CheckRunning() error {
	log.Infof("Checking %s is running", infra.Name)
	var err error

	// check http
	for _, h := range infra.Health.HttpGet {
		if !Ctx.Slaver.CheckHealth(
			infra.Name,
			&pb.Check{
				Type:   "httpGet",
				Host:   infra.renderValue(h.Host),
				Port:   h.Port,
				Schema: "http",
				Path:   h.Path,
			},
		) {
			err = errors.Errorf("Waiting %s running timeout", infra.Name)
		}
	}

	// check socket
	for _, s := range infra.Health.Socket {
		if !Ctx.Slaver.CheckHealth(
			infra.Name,
			&pb.Check{
				Type:   "socket",
				Host:   infra.renderValue(s.Host),
				Port:   s.Port,
				Schema: "",
			},
		) {
			err = errors.Errorf("Waiting %s running timeout", infra.Name)
		}
	}

	return err
}

// 获取基础组件信息
/**
读取安装成功或者用户配置的信息
*/
func (infra *InfraResource) GetResource(key string) *config.Resource {
	news := Ctx.GetSucceed(key, ReleaseTYPE)
	// get info from succeed
	if news != nil {
		return &news.Resource
	} else {
		if r, ok := Ctx.UserConfig.Spec.Resources[key]; ok {
			return r
		}
	}
	log.Errorf("can't get required resource [%s]", key)
	os.Exit(188)
	return nil
}

func (infra *InfraResource) CheckInstall() error {
	news := Ctx.GetSucceed(infra.Name, ReleaseTYPE)

	// check requirement started
	for _, r := range infra.Requirements {
		i := infra.GetInfra(r)
		if err := i.CheckRunning(); err != nil {
			return err
		}
	}
	// apply resource
	if err := infra.applyUserResource(); err != nil {
		return err
	}
	// 初始化value
	if err := infra.executePreValues(); err != nil {
		return err
	}
	if news != nil {
		log.Successf("using exist release %s", news.RefName)
		if news.Status == CreatedStatus {
			infra.CheckExecuteAfterTasks()
		}
		return nil
	}
	// 执行安装前命令
	if err := infra.executePreCommands(); err != nil {
		return err
	}

	statusCh := make(chan error)

	go infra.Install(statusCh)

	for {
		select {
		case err:=<-statusCh:
			return err
		case <-time.Tick(time.Second * 10):
			log.Infof("still install %s", infra.Name)
		}

	}
	return nil
}
