package install

import (
	"bytes"
	"fmt"
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/helm"
	"github.com/choerodon/c7n/pkg/kube"
	"github.com/choerodon/c7n/pkg/slaver"
	"github.com/vinkdong/gox/log"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"text/template"
	"time"
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

type InfraResource struct {
	Name         string
	Chart        string
	Namespace    string
	RepoURL      string
	Version      string
	Values       []ChartValue
	Persistence  []*Persistence
	Client       *helm.Client
	Home         *Install
	Resource     *config.Resource
	PreInstall   []PreInstall
	AfterInstall []PreInstall
	PreValues    PreValueList
	Requirements []string
	Health       Health
}

type Health struct {
	HttpGet []HttpGetCheck `yaml:"httpGet"`
	Socket  []SocketCheck
}

type SocketCheck struct {
	Name string
	Host string
	Port int32
	Path string
}

type HttpGetCheck struct {
	Name string
	Host string
	Port int32
	Path string
}

type Spec struct {
	Basic        Basic
	Resources    v1.ResourceRequirements
	Infra        []InfraResource
	Framework []InfraResource
}

type Basic struct {
	RepoURL string
	Slaver  slaver.Slaver
}

type PreInstall struct {
	Name     string
	Commands []string
	Request  *Request
	InfraRef string `yaml:"infraRef"`
	Opens    []string
}

type Request struct {
	Header []ChartValue
	Url    string
	Body   string
	Method string
}

func (r *Request) Render(infra *InfraResource) error {
	r.Url = infra.renderValue(r.Url)
	r.Body = infra.renderValue(r.Body)
	for k, v := range r.Header {
		v.Value = infra.renderValue(v.Value)
		r.Header[k] = v
	}
	return nil
}

func (pi *PreInstall) ExecuteCommands(infra *InfraResource) error {
	for k, v := range pi.Commands {
		pi.Commands[k] = infra.renderValue(v)
	}
	r := infra.GetResource(pi.InfraRef)
	s := Ctx.Slaver
	s.ExecuteRemoteSql(pi.Commands, r)
	return nil
}

func (pi *PreInstall) ExecuteRequests(infra *InfraResource) error {
	if pi.Request == nil {
		return nil
	}
	pi.Request.Render(infra)
	req := pi.Request
	s := Ctx.Slaver
	header := make(map[string][]string)
	for _, h := range req.Header {
		header[h.Name] = []string{h.Value}
	}
	f := slaver.Forward{
		Url:    req.Url,
		Body:   req.Body,
		Header: header,
		Method: req.Method,
	}
	err := s.ExecuteRemoteRequest(f)
	return err
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

func (p *PreValue) RandomToken(length int) string {
	return RandomToken(length)
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
	val := data.String()
	log.Debugf("PreValue %s: %s, checking: %s", p.Name, val, p.Check)

	switch p.Check {
	case "clusterdomain":
		//todo: add check domain
		if err := Ctx.Slaver.CheckClusterDomain(val); err != nil {
			return err
		}
	}

	p.Value = val
	return nil
}

// 获取基础组件信息
func (p *PreValue) GetResource(key string) *config.Resource {
	news := Ctx.GetSucceed(key, ReleaseTYPE)
	// get info from succeed
	if news != nil {
		return &news.Resource
	} else {
		// 从用户配置文件中读取
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

func (i *Install) Install(apps []InfraResource) error {
	// 安装基础组件
	for _, infra := range apps {
		if r := i.UserConfig.GetResource(infra.Name); r != nil && r.External {
			log.Infof("using external %s", infra.Name)
			continue
		}
		// 准备pv和pvc
		if err := infra.preparePersistence(i.Client, i.UserConfig, i.CommonLabels); err != nil {
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

func (i *Install) PrepareSlaverPvc() (string, error) {
	pvs := i.UserConfig.Spec.Persistence.GetPersistentVolumeSource("")
	persistence := Persistence{
		Client:       i.Client,
		CommonLabels: i.CommonLabels,
		AccessModes:  []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
		Size:         "1Gi",
		Mode:         "755",
		PvcEnabled:   true,
		Name:         "slaver",
	}
	err := persistence.CheckOrCreatePv(pvs)
	if err != nil {
		return "", err
	}

	persistence.Namespace = i.UserConfig.Metadata.Namespace

	if err := persistence.CheckOrCreatePvc(); err != nil {
		return "", err
	}
	return persistence.RefPvcName, nil
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

	if pvcName, err := i.PrepareSlaverPvc(); err != nil {
		return err
	} else {
		s.PvcName = pvcName
	}

	if _, err := s.CheckInstall(); err != nil {
		return err
	}

	stopCh := make(chan struct{})

	port := s.ForwardPort("http", stopCh)
	grpcPort := s.ForwardPort("grpc", stopCh)
	s.Address = fmt.Sprintf("http://127.0.0.1:%d", port)
	s.GRpcAddress = fmt.Sprintf("127.0.0.1:%d", grpcPort)

	Ctx.SlaverAddress = fmt.Sprintf("http://127.0.0.1:%d", port)
	defer func() {
		stopCh <- struct{}{}
	}()

	// install 基础组件
	if err := i.Install(i.Spec.Infra); err != nil {
		return err
	}

	// install 框架微服务
	log.Info("start install choerodon-framework")
	if err := i.Install(i.Spec.Framework) ; err != nil{
		return err
	}

loop:
	for {
		select {
		case <-time.Tick(time.Second * 3):
			if !Ctx.HasBackendTask() {
				break loop
			}
		}
	}

	return nil
}
