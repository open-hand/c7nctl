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
	"k8s.io/kubernetes/pkg/util/maps"
	"os"
	"text/template"
	syserr "errors"
	"github.com/choerodon/c7n/pkg/common"
)

type Install struct {
	Version      string
	Metadata     Metadata
	Spec         Spec
	Client       kubernetes.Interface
	UserConfig   *config.Config
	HelmClient   *helm.Client
	CommonLabels map[string]string
	Namespace    string
	Timeout      int
	Prefix       string
	SkipInput    bool
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
	Timeout      int
	Prefix       string
	SkipInput    bool
}

type Health struct {
	HttpGet   []HttpGetCheck `yaml:"httpGet"`
	Socket    []SocketCheck
	PodStatus []PodCheck     `yaml:"podStatus"`
}

type PodCheck struct {
	Name      string
	Status    v1.PodPhase
	Namespace string
	Client    kubernetes.Interface
}

func (p *PodCheck) MustRunning() error {
	po, err := p.Client.CoreV1().Pods(p.Namespace).Get(p.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	if status := po.Status.Phase; status != p.Status {
		return syserr.New(fmt.Sprintf("[ %s ] pod status is %s, need %s", p.Name, status, p.Status))
	}

	return nil
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
	Basic       Basic
	Resources   v1.ResourceRequirements
	Infra       []*InfraResource
	Framework   []*InfraResource
	DevOps      []*InfraResource `json:"devOps"`
	Agile       []*InfraResource `json:"agile"`
	TestManager []*InfraResource `json:"testManager"`
	Front       []*InfraResource `json:"front"`
	Wiki        []*InfraResource `json:"wiki"`
	Runner      *InfraResource   `json:"runner"`
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

	news := Ctx.GetSucceedTask(pi.Name, infra.Name, SqlTask)
	if news != nil {
		log.Successf("task %s had executed", pi.Name)
		return nil
	}

	news = &News{
		Name:     pi.Name,
		RefName:  infra.Name,
		Type:     TaskType,
		Status:   SucceedStatus,
		TaskType: SqlTask,
		Version:  infra.Version,
	}

	defer Ctx.SaveNews(news)

	for k, v := range pi.Commands {
		pi.Commands[k] = infra.renderValue(v)
	}
	r := infra.GetResource(pi.InfraRef)
	s := Ctx.Slaver
	if err := s.ExecuteRemoteSql(pi.Commands, r); err != nil {
		news.Status = FailedStatus
		news.Reason = err.Error()
		return err
	}
	return nil
}

func (pi *PreInstall) ExecuteRequests(infra *InfraResource) error {
	if pi.Request == nil {
		return nil
	}
	news := Ctx.GetSucceedTask(pi.Name, infra.Name, HttpGetTask)
	if news != nil {
		log.Successf("task %s had executed", pi.Name)
		return nil
	}

	news = &News{
		Name:     pi.Name,
		RefName:  infra.Name,
		Type:     TaskType,
		Status:   SucceedStatus,
		TaskType: HttpGetTask,
		Version:  infra.Version,
	}

	defer Ctx.SaveNews(news)

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
	_, err := s.ExecuteRemoteRequest(f)
	if err != nil {
		news.Status = FailedStatus
		news.Reason = err.Error()
	}
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
	Input common.Input
	Case  string
}

type PreValue struct {
	Name  string
	Value string
	Check string
}

func (p *PreValue) RandomToken(length int) string {
	return RandomToken(length)
}

func (p *PreValue) RandomLowCaseToken(length int) string {
	return GenerateRunnerToken(length)
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
			log.Errorf("请检查您的域名: %s 已正确解析到集群", val)
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
	Ctx.CheckExist(188)
	return nil
}

func (i *Install) CleanJobs() error {
	jobInterface := i.Client.BatchV1().Jobs(i.UserConfig.Metadata.Namespace)
	jobList, err := jobInterface.List(meta_v1.ListOptions{})
	if err != nil {
		return err
	}
	log.Info("clean history jobs...")
	delOpts := &meta_v1.DeleteOptions{}
	for _, job := range jobList.Items {
		if job.Status.Active > 0 {
			log.Infof("job %s still active ignored..", job.Name)
		} else {
			if err := jobInterface.Delete(job.Name, delOpts); err != nil {
				return err
			}
			log.Successf("deleted job %s", job.Name)
		}
		log.Info(job.Name)
	}
	return nil
}

func (i *Install) Install(apps []*InfraResource) error {
	// 安装基础组件
	for _, infra := range apps {
		log.Infof("start install %s",infra.Name)

		infra.SkipInput = i.SkipInput

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
		infra.Timeout = i.Timeout
		infra.Prefix = i.Prefix
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

	Ctx.Metrics.Memory = clusterMemory
	Ctx.Metrics.CPU = clusterCpu
	fmt.Println(i.Client.CoreV1().RESTClient().APIVersion().Version)

	serverVersion, err := i.Client.Discovery().ServerVersion()
	if err != nil {
		log.Error("can't get your cluster version")
	}
	Ctx.Metrics.Version = serverVersion.String()
	if clusterMemory < reqMemory {
		log.Errorf("cluster memory not enough, request %dGi", reqMemory/(1024*1024*1024))
		return false
	}
	if clusterCpu < reqCpu {
		log.Errorf("cluster cpu not enough, request %dc", reqCpu/1000)
		return false
	}
	if !i.SkipInput {
		Ctx.Metrics.Mail, err = common.AcceptUserInput(common.Input{
			Password: false,
			Tip:      "请输入您的邮箱：\nPlease enter your email address:\n",
			Regex:    "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
		})
		if err != nil {
			log.Error(err)
		}
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
	if i.UserConfig == nil {
		return "",nil
	}
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

func (i *Install) PrepareSlaver(stopCh <-chan struct{}) (*slaver.Slaver, error) {
	// prepare slaver to execute sql or make directory ..

	s := &i.Spec.Basic.Slaver
	s.Client = i.Client
	// be care of use point
	s.CommonLabels = maps.CopySS(i.CommonLabels)
	s.Namespace = i.Namespace

	if pvcName, err := i.PrepareSlaverPvc(); err != nil {
		return s,err
	} else {
		s.PvcName = pvcName
	}

	if _, err := s.CheckInstall(); err != nil {
		return s,err
	}
	port := s.ForwardPort("http", stopCh)
	grpcPort := s.ForwardPort("grpc", stopCh)
	s.Address = fmt.Sprintf("http://127.0.0.1:%d", port)
	s.GRpcAddress = fmt.Sprintf("127.0.0.1:%d", grpcPort)
	return s,nil
}

func (i *Install) Run(args ...string) error {

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
		Metrics:      Ctx.Metrics,
	}

	stopCh := make(chan struct{})

	s,err := i.PrepareSlaver(stopCh)
	if err != nil {
		return err
	}
	Ctx.Slaver = s
	defer func() {
		stopCh <- struct{}{}
	}()

	// 清理历史的job
	if err := i.CleanJobs(); err != nil {
		return err
	}
	// install 基础组件
	if err := i.Install(i.Spec.Infra); err != nil {
		return err
	}

	// install 框架微服务
	log.Info("start install choerodon:framework")
	if err := i.Install(i.Spec.Framework); err != nil {
		return err
	}

	// install 框架持续交付
	log.Info("start install choerodon:devops")
	if err := i.Install(i.Spec.DevOps); err != nil {
		return err
	}

	// install 前端服务
	log.Info("start install choerodon:front ")
	if err := i.Install(i.Spec.Front); err != nil {
		return err
	}

	// install 敏捷服务
	log.Info("start install choerodon:agile")
	if err := i.Install(i.Spec.Agile); err != nil {
		return err
	}

	// install 测试管理服务
	log.Info("start install choerodon:test manager")
	if err := i.Install(i.Spec.TestManager); err != nil {
		return err
	}

	// install 知识管理服务
	log.Info("start install choerodon:wiki manager")
	if err := i.Install(i.Spec.Wiki); err != nil {
		return err
	}

	Ctx.CheckExist(0)

	return nil
}
