package resource

import (
	"fmt"
	c7n_config "github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/vinkdong/gox/log"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/util/maps"
	"time"
)

const (
	DEFAULT_REPO_URL = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	C7nLabelKey      = "c7n-usage"
	C7nLabelValue    = "c7n-installer"
)

type InstallDefinition struct {
	// api 版本
	Version string
	// Choerodon 平台版本
	PaaSVersion string
	Metadata    Metadata
	Spec        Spec
	// TODO REMOVE
	Resource           *c7n_config.Resource
	CommonLabels       map[string]string
	DefaultAccessModes []v1.PersistentVolumeAccessMode `yaml:"accessModes"`
	SkipInput          bool
	Timeout            int
	Prefix             string
	Namespace          string
	Mail               string
}

type Metadata struct {
	Name      string
	Namespace string
}

type Spec struct {
	Basic Basic
	// Funcs       []utils.Func
	Resources v1.ResourceRequirements
	Release   []*Release
	Runner    *Release `json:"runner"`
}

type Basic struct {
	RepoURL string
	Slaver  slaver.Slaver
}

func (i *InstallDefinition) CleanJobs() error {
	jobInterface := (*context.Ctx.KubeClient).BatchV1().Jobs(context.Ctx.UserConfig.Metadata.Namespace)
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

func (i *InstallDefinition) Install(apps []*Release) error {
	// 安装基础组件
	for _, infra := range apps {
		log.Infof("start resource %s", infra.Name)

		infra.SkipInput = i.SkipInput

		if r := context.Ctx.UserConfig.GetResource(infra.Name); r != nil && r.External {
			log.Infof("using external %s", infra.Name)
			continue
		}
		// apply configs

		infra.Timeout = i.Timeout
		infra.Namespace = context.Ctx.UserConfig.Metadata.Namespace
		infra.Prefix = i.Prefix
		infra.PaaSVersion = i.Version

		// 准备pv和pvc
		if err := i.PreparePersistence(infra); err != nil {
			return err
		}

		if infra.RepoURL == "" {
			infra.RepoURL = i.Spec.Basic.RepoURL
		}
		if err := i.CheckInstall(infra); err != nil {
			return err
		}
	}
	return nil
}

func (i *InstallDefinition) CheckResource() bool {
	client := *context.Ctx.KubeClient
	request := i.Spec.Resources.Requests
	reqMemory := request.Memory().Value()
	reqCpu := request.Cpu().Value()
	clusterMemory, clusterCpu := getClusterResource(client)

	context.Ctx.Metrics.Memory = clusterMemory
	context.Ctx.Metrics.CPU = clusterCpu

	serverVersion, err := client.Discovery().ServerVersion()
	if err != nil {
		log.Error("can't get your cluster version")
	}
	context.Ctx.Metrics.Version = serverVersion.String()
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

func (i *InstallDefinition) CheckNamespace() bool {
	_, err := (*context.Ctx.KubeClient).CoreV1().Namespaces().Get(context.Ctx.UserConfig.Metadata.Namespace, meta_v1.GetOptions{})
	if err != nil {
		if errorStatus, ok := err.(*errors.StatusError); ok {
			if errorStatus.ErrStatus.Code == 404 && i.createNamespace() {
				return true
			}
		}
		log.Error(err)
		return false
	}
	log.Infof("namespace %s already exists", context.Ctx.UserConfig.Metadata.Namespace)
	return true
}

func (i *InstallDefinition) createNamespace() bool {
	ns := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: context.Ctx.UserConfig.Metadata.Namespace,
		},
	}
	namespace, err := (*context.Ctx.KubeClient).CoreV1().Namespaces().Create(ns)
	log.Infof("creating namespace %s", namespace.Name)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (i *InstallDefinition) PrepareSlaverPvc() (string, error) {
	if context.Ctx.UserConfig == nil {
		return "", nil
	}
	pvs := context.Ctx.UserConfig.Spec.Persistence.GetPersistentVolumeSource("")

	persistence := Persistence{
		Client:       *context.Ctx.KubeClient,
		CommonLabels: i.CommonLabels,
		AccessModes:  i.DefaultAccessModes,
		Size:         "1Gi",
		Mode:         "755",
		PvcEnabled:   true,
		Name:         "slaver",
	}
	err := persistence.CheckOrCreatePv(pvs)
	if err != nil {
		return "", err
	}

	persistence.Namespace = context.Ctx.UserConfig.Metadata.Namespace

	if err := persistence.CheckOrCreatePvc(); err != nil {
		return "", err
	}
	return persistence.RefPvcName, nil
}

func (i *InstallDefinition) PrepareSlaver(stopCh <-chan struct{}) (*slaver.Slaver, error) {
	// prepare slaver to execute sql or make directory ..

	s := &i.Spec.Basic.Slaver
	s.Client = *context.Ctx.KubeClient
	// be care of use point
	s.CommonLabels = maps.CopySS(context.Ctx.CommonLabels)
	s.Namespace = context.Ctx.Namespace

	if pvcName, err := i.PrepareSlaverPvc(); err != nil {
		return s, err
	} else {
		s.PvcName = pvcName
	}

	if _, err := s.CheckInstall(); err != nil {
		return s, err
	}
	port := s.ForwardPort("http", stopCh)
	grpcPort := s.ForwardPort("grpc", stopCh)
	s.Address = fmt.Sprintf("http://127.0.0.1:%d", port)
	s.GRpcAddress = fmt.Sprintf("127.0.0.1:%d", grpcPort)
	return s, nil
}

func (i *InstallDefinition) PreparePersistence(infra *Release) error {
	cfg := context.Ctx.UserConfig
	getPvs := cfg.Spec.Persistence.GetPersistentVolumeSource

	// use app defined persistence
	if res, ok := cfg.Spec.Resources[infra.Name]; ok && res.Persistence != nil {
		getPvs = res.Persistence.GetPersistentVolumeSource
	}

	namespace := cfg.Metadata.Namespace
	context.Ctx.CommonLabels["app"] = infra.Name
	for _, persistence := range infra.Persistence {
		persistence.Client = *context.Ctx.KubeClient

		persistence.CommonLabels = context.Ctx.CommonLabels
		persistence.CommonLabels["pv"] = persistence.Name
		if len(persistence.AccessModes) == 0 {
			persistence.AccessModes = i.DefaultAccessModes
		}

		if cfg.Spec.Persistence.GetStorageType() == c7n_config.PersistenceHostPathType {
			persistence.MountOptions = []string{}
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

// get server definition
func (id *InstallDefinition) GetRelease(key string) *Release {
	for _, rls := range id.Spec.Release {
		if rls.Name == key {
			return rls
		}
	}
	return nil
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

func (i *InstallDefinition) CheckInstall(infra *Release) error {
	news := context.Ctx.GetSucceed(infra.Name, context.ReleaseTYPE)

	// check requirement started
	for _, r := range infra.Requirements {
		i := i.GetRelease(r)
		if i.Prefix == "" {
			i.Prefix = infra.Prefix
		}
		/*if err := i.CheckRunning(); err != nil {
			return err
		}*/
	}
	// apply resource
	if err := infra.ApplyUserResource(); err != nil {
		return err
	}

	if news != nil {
		log.Successf("using exist release %s", news.RefName)
		if news.Status != context.SucceedStatus {
			//infra.PreValues = news.PreValue
			infra.ExecuteAfterTasks()
		}
		return nil
	}
	// 初始化value
	if err := infra.ExecutePreValues(); err != nil {
		return err
	}

	// 执行安装前命令
	if err := infra.ExecutePreCommands(); err != nil {
		return err
	}

	statusCh := make(chan error)

	go func() {
		err := infra.Install()
		statusCh <- err
	}()

	for {
		select {
		case <-time.Tick(time.Second * 10):
			_ = infra.CatchInitJobs()
		case err := <-statusCh:
			return err
		}
	}
}

func getAppFromList(appName string, resourceList []*Release) *Release {
	for _, v := range resourceList {
		if v.Name == appName {
			//v.convertInstalledValue()
			return v
		}
	}
	return nil
}

func (i *InstallDefinition) getValue(release, key string) string {
	rel := i.GetRelease(release)
	for _, v := range rel.Values {
		if v.Name == key {
			return v.Value
		}
	}
	log.Infof("can't get value '%s' of %s", key, rel.Name)
	return ""
}

func (i *InstallDefinition) getPreValue(release, key string) string {
	_ = i.GetRelease(release)
	// return rel.PreValues.getValues(key)
	return ""
}

func getResource(key string) *c7n_config.Resource {
	news := context.Ctx.GetSucceed(key, context.ReleaseTYPE)
	// get info from succeed
	if news != nil {
		return &news.Resource
	} else {
		// 从用户配置文件中读取
		if r, ok := context.Ctx.UserConfig.Spec.Resources[key]; ok {
			return r
		}
	}
	log.Errorf("can't get required resource [%s]", key)
	context.Ctx.CheckExist(188)
	return nil
}
