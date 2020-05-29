package action

import (
	"encoding/json"
	"github.com/choerodon/c7nctl/pkg/config"
	c7n_ctx "github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/resource"
	c7n_utils "github.com/choerodon/c7nctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	yaml_v2 "gopkg.in/yaml.v2"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"time"
)

// TODO is necessary move to consts pkg ?
const (
	DefaultRepoUrl = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	C7nLabelKey    = "c7n-usage"
	C7nLabelValue  = "c7n-installer"
)

type InstallC7n struct {
	cfg *Configuration
	// api versions
	Version string
	// Choerodon version
	PassVersion string
	// choerodon install configuration
	ConfigFile string
	// install
	ResourceFile       string
	Prefix             string
	NoTimeout          bool
	SkipInput          bool
	Namespace          string
	Timeout            int
	Mail               string
	CommonLabels       map[string]string
	DefaultAccessModes []v1.PersistentVolumeAccessMode `yaml:"accessModes"`
}

func NewInstall(cfg *Configuration) *InstallC7n {
	return &InstallC7n{cfg: cfg}
}

//
func (i *InstallC7n) Run() error {
	// 当 version 没有设置时，从 git repo 获取最新版本
	i.Version = c7n_utils.GetVersion(i.Version)

	userConfig := getUserConfig(i.ConfigFile)
	if userConfig == nil {
		log.Info("need user config file")
		os.Exit(127)
	}

	initContext(i, userConfig)
	id := i.getInstallDef(userConfig)

	// 检查硬件资源
	if !checkResource(&id.Spec.Resources) {
		os.Exit(126)
	}
	if !checkNamespace() {
		os.Exit(127)
	}

	stopCh := make(chan struct{})
	// TODO move method PrepareSlaver()
	s, err := id.PrepareSlaver(stopCh)
	if err != nil {
		return err
	}

	c7n_ctx.Ctx.Slaver = s
	defer func() {
		stopCh <- struct{}{}
	}()

	_ = c7n_ctx.Ctx.LoadJobInfoFromCM()
	// 渲染 Release
	for _, rls := range id.Spec.Release {
		// 传入参数的是 *Release
		renderRelease(rls)
	}

	releaseGraph := resource.NewReleaseGraph(id)
	installQueue := releaseGraph.TopoSortByKahn()

	for !installQueue.IsEmpty() {
		rls := installQueue.Dequeue()

		if err = rls.Install(); err != nil {
			log.Error(err)
		}
	}

	// 清理历史的job
	if err := cleanJobs(); err != nil {
		return err
	}

	c7n_ctx.Ctx.CheckExist(0)
	return nil
}

func (i *InstallC7n) InstallComponent(cname string) error {
	i.Version = c7n_utils.GetVersion(i.Version)

	id := i.getInstallDef(nil)
	c7n_ctx.Ctx.HelmClient = i.cfg.HelmClient
	c7n_ctx.Ctx.KubeClient = i.cfg.KubeClient
	c7n_ctx.Ctx.Namespace = i.Namespace
	c7n_ctx.Ctx.RepoUrl = DefaultRepoUrl

	/*	stopCh := make(chan struct{})
		// TODO move method PrepareSlaver()
		s, err := id.PrepareSlaver(stopCh)
		if err != nil {
			return err
		}

		c7n_ctx.Ctx.Slaver = s
		defer func() {
			stopCh <- struct{}{}
		}()*/

	for _, rls := range id.Spec.Component {
		if rls.Name == cname {
			renderComponent(rls)

			rls.Name = rls.Name + "-" + c7n_utils.RandomString(5)
			if err := rls.InstallComponent(); err != nil {
				return err
			} else {
				break
			}
		}
	}
	return nil
}

func renderRelease(rls *resource.Release) {
	ji := c7n_ctx.Ctx.GetJobInfo(rls.Name)

	if ji.Name == "" {
		ji = c7n_ctx.JobInfo{
			Name:      rls.Name,
			Namespace: c7n_ctx.Ctx.Namespace,
			Type:      c7n_ctx.ReleaseTYPE,
			Status:    c7n_ctx.UninitializedStatus,
			Date:      time.Now(),
		}
		c7n_ctx.Ctx.UpdateJobInfo(ji)
	}
	if ji.Status == c7n_ctx.UninitializedStatus {
		// 传入的参数是指针
		mergerResource(rls)
		ji.Resource = *rls.Resource
		c7n_ctx.Ctx.UpdateJobInfo(ji)

		renderValues(rls)
		ji.Values = rls.Values
		ji.Status = c7n_ctx.InputtedStatus
		c7n_ctx.Ctx.UpdateJobInfo(ji)
	}
	if ji.Status == c7n_ctx.InputtedStatus {
		rlsByte, _ := yaml_v2.Marshal(rls)
		renderedRls := c7n_utils.RenderRelease(rls.Name, string(rlsByte))
		_ = yaml_v2.Unmarshal(renderedRls, rls)
		// 检查域名
		_ = checkReleaseDomain(rls)

		ji.Resource = *rls.Resource
		ji.Values = rls.Values
		ji.Status = c7n_ctx.RenderedStatus
		// 保存渲染完成的rls
		c7n_ctx.Ctx.UpdateJobInfo(ji)
	}
	// 当 rls 渲染完成但是没有完成安装——c7nctl install 会中断，二次执行
	if ji.Status == c7n_ctx.RenderedStatus {
		rls.Values = ji.Values
		rls.Resource = &ji.Resource
	}
}

func renderComponent(rls *resource.Release) {
	renderValues(rls)
	rlsByte, _ := yaml_v2.Marshal(rls)
	renderedRls := c7n_utils.RenderRelease(rls.Name, string(rlsByte))
	_ = yaml_v2.Unmarshal(renderedRls, rls)
}

func checkReleaseDomain(rls *resource.Release) error {
	for _, v := range rls.Values {
		if v.Check == "clusterdomain" {
			log.Debugf("Value %s: %s, checking: %s", v.Name, v.Value, v.Check)
			if err := c7n_ctx.Ctx.Slaver.CheckClusterDomain(v.Value); err != nil {
				log.Errorf("请检查您的域名: %s 已正确解析到集群", v.Value)
				return err
			}
		}
	}
	return nil
}

// 传指针的方式好呢，还是返回值的方式好？
//
// 在渲染 release 前将 values 渲染完成
func renderValues(rls *resource.Release) {
	if rls.Values == nil {
		log.Debugf("release %s values is empty", rls.Name)
		return
	}
	for idx, v := range rls.Values {
		// 输入 value
		if v.Input.Enabled && !c7n_ctx.Ctx.SkipInput {
			var err error
			var value string
			if v.Input.Password {
				v.Input.Twice = true
				value, err = c7n_utils.AcceptUserPassword(v.Input)
			} else {
				value, err = c7n_utils.AcceptUserInput(v.Input)
			}
			// v.Values 是复制
			rls.Values[idx].Value = value
			if err != nil {
				log.Error(err)
				os.Exit(128)
			}
		} else {
			rls.Values[idx].Value = c7n_utils.RenderReleaseValue(v.Name, v.Value)
		}
	}
}

// 将 config.yml 中的值合并到 Release.Resource
func mergerResource(r *resource.Release) {
	cnf := c7n_ctx.Ctx.GetConfig()
	if res := cnf.GetResource(r.Name); res == nil {
		log.Warnf("There is no resource in config.yaml of Release %s", r.Name)
	} else {
		// 直接使用外部配置
		if res.External {
			r.Resource = res
		} else {
			// TODO 有没有更加简便的方式
			if res.Domain != "" {
				r.Resource.Domain = res.Domain
			}
			if res.Schema != "" {
				r.Resource.Schema = res.Schema
			}
			if res.Url != "" {
				r.Resource.Url = res.Url
			}
			if res.Host != "" {
				r.Resource.Host = res.Host
			}
			if res.Port > 0 {
				r.Resource.Port = res.Port
			}
			if res.Username != "" {
				r.Resource.Username = res.Username
			}
			if res.Password != "" {
				r.Resource.Password = res.Password
			}
			if res.Persistence != nil {
				r.Resource.Persistence = res.Persistence
			}
		}
	}
}

// 为了避免循环依赖，从 install_definition.go 移到这里
func (i *InstallC7n) getInstallDef(uc *config.C7nConfig) *resource.InstallDefinition {
	var rd []byte
	if i.ResourceFile != "" {
		rd = c7n_utils.GetResourceFile(false, i.Version, i.ResourceFile)
	} else {
		rd = c7n_utils.GetResourceFile(true, i.Version, c7n_utils.InstallConfigPath)
	}
	// 只获取 installDef，不做全局渲染，gitlab Token 需要手动获取
	// rdRender := c7n_utils.RenderInstallDef(string(rd))

	installDef := &resource.InstallDefinition{}
	rdJson, err := yaml.ToJSON(rd)
	if err != nil {
		panic(err)
	}
	// slaver 使用了 core_v1.ContainerPort, 必须先转 JSON
	_ = json.Unmarshal(rdJson, installDef)

	// TODO PaasVersion and Timeout is necessary?
	installDef.PaaSVersion = i.Version
	if i.NoTimeout {
		installDef.Timeout = 60 * 60 * 24
	}

	if uc != nil {
		if accessModes := uc.Spec.Persistence.AccessModes; len(accessModes) > 0 {
			installDef.DefaultAccessModes = accessModes
		} else {
			installDef.DefaultAccessModes = []v1.PersistentVolumeAccessMode{"ReadWriteOnce"}
		}
	}

	return installDef
}

func getUserConfig(filePath string) *config.C7nConfig {
	if filePath == "" {
		log.Debugf("no user config defined by `-c`")
		return nil
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		os.Exit(124)
	}
	userConfig := &config.C7nConfig{}
	err = yaml_v2.Unmarshal(data, userConfig)
	if err != nil {
		log.Error(err)
		os.Exit(124)
	}
	return userConfig
}

// mv to client package
func cleanJobs() error {
	jobInterface := (*c7n_ctx.Ctx.KubeClient).BatchV1().Jobs(c7n_ctx.Ctx.UserConfig.Metadata.Namespace)
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
			log.Info("deleted job %s", job.Name)
		}
		log.Info(job.Name)
	}
	return nil
}

func initContext(i *InstallC7n, userConfig *config.C7nConfig) {
	commonLabels := map[string]string{
		C7nLabelKey: C7nLabelValue,
	}
	i.CommonLabels = commonLabels

	c7n_ctx.Ctx.Metrics.Mail = i.Mail
	c7n_ctx.Ctx = c7n_ctx.Context{
		// also init i.cfg
		HelmClient:   i.cfg.HelmClient,
		KubeClient:   i.cfg.KubeClient,
		Namespace:    userConfig.Metadata.Namespace,
		CommonLabels: i.CommonLabels,
		UserConfig:   userConfig,
		Metrics:      c7n_ctx.Metrics{},
		JobInfo:      map[string]c7n_ctx.JobInfo{},
		SkipInput:    i.SkipInput,
		Prefix:       i.Prefix,
		// TODO 根据 install.yaml 配置
		RepoUrl: DefaultRepoUrl,
	}
}

func checkResource(resources *v1.ResourceRequirements) bool {
	client := *c7n_ctx.Ctx.KubeClient
	request := resources.Requests

	reqMemory := request.Memory().Value()
	reqCpu := request.Cpu().Value()
	clusterMemory, clusterCpu := c7n_utils.GetClusterResource(client)

	c7n_ctx.Ctx.Metrics.Memory = clusterMemory
	c7n_ctx.Ctx.Metrics.CPU = clusterCpu

	serverVersion, err := client.Discovery().ServerVersion()
	if err != nil {
		log.Error("can't get your cluster version")
	}
	c7n_ctx.Ctx.Metrics.Version = serverVersion.String()
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

func checkNamespace() bool {
	_, err := (*c7n_ctx.Ctx.KubeClient).CoreV1().Namespaces().Get(c7n_ctx.Ctx.UserConfig.Metadata.Namespace, meta_v1.GetOptions{})
	if err != nil {
		if errorStatus, ok := err.(*errors.StatusError); ok {
			if errorStatus.ErrStatus.Code == 404 && c7n_utils.CreateNamespace() {
				return true
			}
		}
		log.Error(err)
		return false
	}
	log.Infof("namespace %s already exists", c7n_ctx.Ctx.UserConfig.Metadata.Namespace)
	return true
}
