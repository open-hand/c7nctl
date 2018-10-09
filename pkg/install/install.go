package install

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"github.com/choerodon/c7n/pkg/kube"
	"github.com/choerodon/c7n/pkg/config"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"github.com/vinkdong/gox/log"
	"github.com/choerodon/c7n/pkg/helm"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"text/template"
	"bytes"
	"github.com/choerodon/c7n/pkg/slaver"
)

type Install struct {
	Version      string
	Metadata     Metadata
	Spec         Spec
	Client       kubernetes.Interface
	UserConfig   *config.Config
	HelmClient   *helm.Client
	CommonLabels map[string]string
}

type Metadata struct {
	Name string
}

type Spec struct {
	Basic     Basic
	Resources v1.ResourceRequirements
	Infra     []InfraResource
}

type Basic struct {
	RepoURL string
	Slaver  slaver.Slaver
}

type PreInstall struct {
	Name     string
	Commands []string
	InfraRef string `yaml:"infraRef"`
}

type InfraResource struct {
	Name        string
	Chart       string
	Namespace   string
	RepoURL     string
	Version     string
	Values      []ChartValue
	Persistence []*Persistence
	Client      *helm.Client
	Home        *Install
	Resource    config.Resource
	PreInstall  []PreInstall
	PreValues   PreValueList
}

type PreValueList []*PreValue

func (pl *PreValueList) prepareValues() error {
	for _, v := range *pl {
		if err := v.renderValue(); err != nil {
			return err
		}
	}
	return nil
}

func (pl *PreValueList) getValues(key string) string {
	for _, v := range *pl {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}

func (infra *InfraResource) executePreCommands() error {
	s := Ctx.Slaver
	for _, pi := range infra.PreInstall {
		for _, c := range pi.Commands {
			if err := s.ExecuteSql(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (infra *InfraResource) preparePersistence(client kubernetes.Interface, config *config.Config) error {
	getPvs := config.Spec.Persistence.GetPersistentVolumeSource
	namespace := config.Metadata.Namespace
	for _, persistence := range infra.Persistence {
		persistence.Client = client

		// check or create dir
		dir := slaver.Dir{
			Mode: persistence.Mode,
			Path: persistence.Path,
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

func (infra *InfraResource) CheckInstall() error {
	news := Ctx.GetSucceed(infra.Name, ReleaseTYPE)
	// 初始化value
	if err := infra.executePreValues(); err != nil {
		return err
	}
	if news != nil {
		log.Infof("using exist release %s", news.RefName)
		return nil
	}
	// 执行安装前命令
	if err := infra.executePreCommands(); err != nil {
		return err
	}
	return infra.Install()
}

func (infra *InfraResource) executePreValues() error {
	return infra.PreValues.prepareValues()
}

// install infra
func (infra *InfraResource) Install() error {
	values := infra.HelmValues()
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
		Status:    FailedStatues,
		Type:      ReleaseTYPE,
		Resource:  infra.Resource,
	}
	defer Ctx.SaveNews(news)

	if err != nil {
		news.Reason = err.Error()
		return err
	}
	news.Status = SucceedStatus
	return nil
}

// convert yml values to values list as xxx=yyy
func (infra *InfraResource) HelmValues() []string {
	values := make([]string, len(infra.Values))
	for k, v := range infra.Values {
		if v.Input.Enabled {
			password, err := AcceptUserPassword(v.Input)
			if err != nil {
				log.Error(err)
				os.Exit(128)
			}
			values[k] = fmt.Sprintf("%s=%s", v.Name, password)
		} else {
			values[k] = fmt.Sprintf("%s=%s", v.Name, infra.renderValue(v.Value))
		}
	}
	return values
}

func (infra *InfraResource) GetPreValue(key string) string {
	return infra.PreValues.getValues(key)
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

type ChartValue struct {
	Name  string
	Value string
	Input Input
}

type PreValue struct {
	Name  string
	Value string
	Check string
}

func (p *PreValue) renderValue() error {
	tpl, err := template.New(p.Name).Parse(p.Value)
	if err != nil {
		return err
	}
	var data bytes.Buffer
	err = tpl.Execute(&data, p)
	if err != nil {
		return err
	}
	switch p.Check {
	case "domain":
		//todo: add check domain
		log.Infof("check domain of %s", data.String())
	}

	p.Value = data.String()
	return nil
}

// 获取基础组件信息
func (p *PreValue) GetResource(key string) *config.Resource {
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

type Input struct {
	Enabled bool
	Regex   string
	Tip     string
}

func (i *Install) InstallInfra() error {
	// 安装基础组件
	for _, infra := range i.Spec.Infra {
		// 准备pv和pvc
		if err := infra.preparePersistence(i.Client, i.UserConfig); err != nil {
			return err
		}
		infra.Client = i.HelmClient
		infra.Namespace = i.UserConfig.Metadata.Namespace
		infra.Home = i
		if infra.RepoURL == "" {
			infra.RepoURL = i.Spec.Basic.RepoURL
		}
		if err := infra.CheckInstall(); err != nil {
			return err
		}
	}
	return nil
}

func (i *Install) CheckResource() bool {
	request := i.Spec.Resources.Requests
	reqMemory := request.Memory().Value()
	reqCpu := request.Cpu().Value()
	clusterMemory, clusterCpu := getClusterResource(i.Client)
	if clusterMemory < reqMemory {
		log.Errorf("cluster memory not enough, request %dGi", reqMemory/(1024*1024*1024))
		return false
	}
	if clusterCpu < reqCpu {
		log.Errorf("cluster cpu not enough, request %dc", reqCpu/1000)
		return false
	}
	return true
}

func (i *Install) CheckNamespace() bool {
	_, err := i.Client.CoreV1().Namespaces().Get(i.UserConfig.Metadata.Namespace, meta_v1.GetOptions{})
	if err != nil {
		if errorStatus, ok := err.(*errors.StatusError); ok {
			if errorStatus.ErrStatus.Code == 404 && i.createNamespace() {
				return true
			}
		}
		log.Error(err)
		return false
	}
	log.Infof("namespace %s already exists", i.UserConfig.Metadata.Namespace)
	return true
}

func (i *Install) createNamespace() bool {
	ns := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: i.UserConfig.Metadata.Namespace,
		},
	}
	namespace, err := i.Client.CoreV1().Namespaces().Create(ns)
	log.Infof("creating namespace %s", namespace.Name)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func getClusterResource(client kubernetes.Interface) (int64, int64) {
	var sumMemory int64
	var sumCpu int64
	list, _ := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	for _, v := range list.Items {
		sumMemory += v.Status.Capacity.Memory().Value()
		sumCpu += v.Status.Capacity.Cpu().Value()
	}
	return sumMemory, sumCpu
}

func (i *Install) Run() error {

	if i.Client == nil {
		i.Client = kube.GetClient()
	}
	if !i.CheckResource() {
		os.Exit(126)
	}

	if !i.CheckNamespace() {
		os.Exit(127)
	}

	if i.HelmClient == nil {
		log.Info("reinit helm client")
		tunnel := kube.GetTunnel()
		i.HelmClient = &helm.Client{
			Tunnel: tunnel,
		}
	}

	Ctx = Context{
		Client:       i.Client,
		Namespace:    i.UserConfig.Metadata.Namespace,
		CommonLabels: i.CommonLabels,
		UserConfig:   i.UserConfig,
	}

	// prepare slaver to execute sql or make directory ..

	s := &i.Spec.Basic.Slaver
	s.Client = i.Client
	s.CommonLabels = i.CommonLabels
	s.Namespace = i.UserConfig.Metadata.Namespace

	Ctx.Slaver = s

	if _, err := s.CheckInstall(); err != nil {
		return err
	}

	// install 基础组件
	if err := i.InstallInfra(); err != nil {
		return err
	}

	return nil
}
