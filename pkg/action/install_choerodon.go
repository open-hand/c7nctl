package action

import (
	"fmt"
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	c7nconsts "github.com/choerodon/c7nctl/pkg/common/consts"
	"github.com/choerodon/c7nctl/pkg/common/graph"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/resource"
	c7nslaver "github.com/choerodon/c7nctl/pkg/slaver"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"os"
	"reflect"
)

type Install struct {
	cfg *C7nConfiguration

	// 安装资源的路径，helm 的 values.yaml 文件在其路径下的 values 文件夹中
	ResourcePath string
	Version      string
	HelmValues   string
	//
	C7nGatewayUrl string
	ClientOnly    bool

	// 以下都是初始化到 InstallDefinition 的配置项
	Prefix          string
	ImageRepository string
	ChartRepository string
	DatasourceTpl   string
	// TODO 是否支持

	ThinMode bool
}

func NewInstall(cfg *C7nConfiguration) *Install {
	return &Install{
		cfg: cfg,
	}
}

// 设置 install 的值，
func (i *Install) InitInstall(c *config.C7nConfig) {
	log.Debug("Initialize config to Install")
	// 配置优先级 flag > config.yaml > 默认值
	// 当 i.Version 不存在时，设置成 c.Version 或者默认值
	if c.Version == "" {
		// TODO 在打包时根据 TAG 设置 version
		c.Version = c7nconsts.Version
	}
	if i.Version == "" {
		i.Version = c.Version
	}
	log.Debugf("Choerodon version is %s", i.Version)
	if c.Spec.ResourcePath == "" {
		// 默认到 github 上获取资源文件
		c.Spec.ResourcePath = fmt.Sprintf(c7nconsts.ResourcePath, i.Version)
	}
	if i.ResourcePath == "" {
		i.ResourcePath = c.Spec.ResourcePath
	}
	log.Debugf("Install file path is %s", i.ResourcePath)
	if i.HelmValues == "" {
		i.HelmValues = c7nconsts.DefaultHelmValuesPath
	}
	log.Debugf("Helm values dir is %s", i.ResourcePath)

	// 配置优先级 flag > config.yaml > install.yml > 默认值
	log.Debug("Initialize flag to C7nConfig")
	if i.ImageRepository != "" {
		c.Spec.ImageRepository = i.ImageRepository
	}
	log.Debugf("Image repository is %s", c.GetImageRepository())
	if i.ChartRepository != "" {
		c.Spec.ChartRepository = i.ChartRepository
	}
	log.Debugf("Chart repository is %s", c.GetImageRepository())

	if i.DatasourceTpl != "" {
		c.Spec.DatasourceTpl = i.DatasourceTpl
	}
	log.Debugf("Datasource template is %s", c.GetImageRepository())

	if i.Prefix != "" {
		c.Spec.Prefix = i.Prefix
	}
	log.Debugf("Prefix is %s", c.GetImageRepository())

	if i.ThinMode {
		c.Spec.ThinMode = true
	}
}

func (i *Install) RenderChoerodon(inst *resource.InstallDefinition, namespace string) (*resource.InstallDefinition, error) {
	if inst == nil {
		return nil, std_errors.New("InstallChoerodon definition can't be empty.")
	}
	// 初始化安装记录
	c7nclient.InitC7nLogs(i.cfg.KubeClient.GetClientSet(), namespace)
	for _, rls := range inst.Spec.Release {
		if err := inst.RenderRelease(rls); err != nil {
			log.Errorf("Release %s render failed: %+v", rls.Name, err)
		}
		if i.ClientOnly {
			fmt.Printf("---------- Helm Release %s ----------\n", rls.Name)
			c7nutils.PrettyPrint(*rls)
		}
	}

	return inst, nil
}

func (i *Install) InstallChoerodon(inst *resource.InstallDefinition, namespace string) error {
	releaseGraph := graph.NewReleaseGraph(inst.Spec.Release)
	installQueue := releaseGraph.TopoSortByKahn()

	for !installQueue.IsEmpty() {
		rls := installQueue.Dequeue()
		log.Infof("start installing release %s", rls.Name)
		// 获取的 values.yaml 必须经过渲染，只能放在 id 中
		vals, err := inst.RenderHelmValues(rls, i.ResourcePath, i.HelmValues)
		if err != nil {
			return err
		}
		args := c7nclient.ChartArgs{
			RepoUrl:     i.ChartRepository,
			Namespace:   namespace,
			ReleaseName: inst.GetReleaseName(rls.Name),
			ChartName:   rls.Chart,
			Version:     rls.Version,
		}

		if i.ClientOnly {
			fmt.Printf("------------- Installingg helm release %s -------------", rls.Name)
			fmt.Println("\nargs:")
			c7nutils.PrettyPrint(args)
			fmt.Println("\nhelm values:")
			c7nutils.PrettyPrint(vals)

			continue
		}
		if err = i.InstallRelease(rls, vals, args, &inst.Spec.Basic.Slaver); err != nil {
			return std_errors.WithMessage(err, fmt.Sprintf("Release %s install failed", rls.Name))
		}
	}
	return nil
}

func (i *Install) InstallRelease(rls *resource.Release, vals map[string]interface{}, args c7nclient.ChartArgs, slaver *c7nslaver.Slaver) error {
	task, err := c7nclient.GetTask(rls.Name)
	if err != nil {
		return err
	}
	if task.Status == c7nconsts.SucceedStatus {
		log.Infof("Release %s is already installed", rls.Name)
		return nil
	}

	if task.Status == c7nconsts.RenderedStatus || task.Status == c7nconsts.FailedStatus {
		// 等待依赖项安装完成
		for _, r := range rls.Requirements {
			rls.CheckReleasePodRunning(r)
		}
		if err := rls.ExecutePreCommands(slaver); err != nil {
			task.Status = c7nconsts.FailedStatus
			return std_errors.WithMessage(err, fmt.Sprintf("Release %s execute pre commands failed", rls.Name))
		}

		log.Infof("installing %s", rls.Name)
		// TODO 使用统一的 io.writer
		// 使用 upgrade --install cmd
		_, err := i.cfg.HelmClient.Upgrade(args, vals, os.Stdout)
		if err != nil {
			task.Status = c7nconsts.FailedStatus
			return err
		}
		task.Status = c7nconsts.InstalledStatus
		// 将异步的 afterInstall 改为同步，AfterInstall 其依赖检查依靠 release
		if len(rls.AfterInstall) > 0 {
			if err := rls.ExecuteAfterTasks(slaver); err != nil {
				task.Status = c7nconsts.FailedStatus
				return std_errors.WithMessage(err, "Execute after task failed")
			}
		}
		task.Status = c7nconsts.SucceedStatus
		log.Infof("Successfully installed %s", rls.Name)
	}
	// 完成后更新 task 状态
	_, err = c7nclient.SaveTask(*task)
	return err
}

func (i *Install) CheckResource(resources *v1.ResourceRequirements, metrics *c7nclient.Metrics) error {
	request := resources.Requests

	reqMemory := request.Memory().Value()
	reqCpu := request.Cpu().Value()
	clusterMemory, clusterCpu := i.cfg.KubeClient.GetClusterResource()

	metrics.Memory = clusterMemory
	metrics.CPU = clusterCpu

	serverVersion, err := i.cfg.KubeClient.GetServerVersion()
	if err != nil {
		return std_errors.Wrap(err, "can't get your cluster version")
	}
	metrics.Version = serverVersion.String()
	if i.ClientOnly {
		log.Info("Running Client only, So skip up check cluster resource")
		return nil
	}
	if clusterMemory < reqMemory {
		return std_errors.New(fmt.Sprintf("cluster memory not enough, request %dGi", reqMemory/(1024*1024*1024)))
	}
	if clusterCpu < reqCpu {
		return std_errors.New(fmt.Sprintf("cluster cpu not enough, request %dc", reqCpu/1000))
	}
	return nil
}

func (i *Install) CheckNamespace(namespace string) error {
	_, err := i.cfg.KubeClient.GetNamespace(namespace)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return i.cfg.KubeClient.CreateNamespace(namespace)
		}
		return err
	}
	log.Infof("namespace %s already exists", namespace)
	return nil
}

func (i *Install) GetClientSet() *kubernetes.Clientset {
	return i.cfg.KubeClient.GetClientSet()
}

func (i *Install) setFiledValue(filed, value interface{}) {
	getType := reflect.TypeOf(i)
	getValue := reflect.ValueOf(&i).Elem()

	for i := 0; i < getType.NumField(); i++ {
		t := getType.Field(i)
		if t.Name == filed {
			switch getValue.FieldByName(t.Name).Kind() {
			case reflect.String:
				{
					getValue.Field(i).SetString(value.(string))
				}
			case reflect.Int:
				{
					getValue.Field(i).SetInt(value.(int64))

				}
			case reflect.Bool:
				{
					getValue.Field(i).SetBool(value.(bool))
				}
			default:
				{
					log.Debugf("Type %s is not support with InstallChoerodon")
				}
			}
		}
	}
}
