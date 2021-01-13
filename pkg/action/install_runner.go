package action

import (
	"fmt"
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	c7nconsts "github.com/choerodon/c7nctl/pkg/common/consts"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/resource"
	"k8s.io/client-go/kubernetes"
)

type InstallRunner struct {
	cfg          *C7nConfiguration
	ResourcePath string
	Version      string
}

func NewInstallRunner(cfg *C7nConfiguration) *InstallRunner {
	return &InstallRunner{
		cfg: cfg,
	}
}

func (ir *InstallRunner) InitInstallRunner(c *config.C7nConfig) {
	// config.yml 中的配置无效

	if c.Spec.ResourcePath == "" {
		// 默认到 github 上获取资源文件
		c.Spec.ResourcePath = fmt.Sprintf(c7nconsts.OpenSourceResourceURL+c7nconsts.OpenSourceResourceBasePath, ir.Version)
	}
	if ir.ResourcePath == "" {
		ir.ResourcePath = c.Spec.ImageRepository
	}
}

func (ir *InstallRunner) RenderGitlabRunner(id *resource.InstallDefinition, namespace string) {
	c7nclient.InitC7nLogs(ir.cfg.KubeClient.GetClientSet(), namespace)
	/*
		if err := id.renderRelease(id.Spec.Runner); err != nil {
			log.Errorf("Release gitlab runner render failed: %+v", err)
		}

	*/
}

func (ir *InstallRunner) InstallGitlabRunner(instDef *resource.InstallDefinition, namespace string) error {

	//log.Infof("start install %s", instDef.Spec.Runner.Name)
	// 获取的 values.yaml 必须经过渲染，只能放在 id 中
	/*
		vals, err := instDef.RenderHelmValues(instDef.Spec.Runner, ir.ResourcePath, c7nconsts.DefaultHelmValuesPath)
		if err != nil {
			return err
		}
	*/
	/*
		var vals map[string]interface{}
		runner := instDef.Spec.Runner
		slaver := instDef.Spec.Basic.Slaver
		args := c7nclient.ChartArgs{
			RepoUrl:     instDef.Spec.Basic.ChartRepository,
			Namespace:   namespace,
			ReleaseName: runner.Namespace,
			ChartName:   runner.Chart,
			Version:     runner.Version,
		}

		// 等待依赖项安装完成
		//for _, r := range runner.Requirements {
		//	runner.CheckReleasePodRunning(r)
		//}

		if err := runner.ExecutePreCommands(&slaver); err != nil {
			return std_errors.WithMessage(err, fmt.Sprintf("Release %s execute pre commands failed", runner.Name))
		}

		log.Infof("installing %s", runner.Name)
		// 使用 upgrade --install cmd
		_, err := ir.cfg.HelmClient.Install(args, vals, os.Stdout)
		if err != nil {
			return err
		}
		// 将异步的 afterInstall 改为同步，AfterInstall 其依赖检查依靠 release
		if len(runner.AfterInstall) > 0 {
			if err := runner.ExecuteAfterTasks(&slaver); err != nil {
				return std_errors.WithMessage(err, "Execute after task failed")
			}
		}
		log.Infof("Successfully installed %s", runner.Name)

		return err
	*/
	return nil
}

func (ir *InstallRunner) GetClientSet() *kubernetes.Clientset {
	return ir.cfg.KubeClient.GetClientSet()
}
