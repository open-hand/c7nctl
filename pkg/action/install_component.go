package action

import (
	"github.com/choerodon/c7nctl/pkg/resource"
)

type InstallComponent struct {
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

func NewInstallComponent(cfg *C7nConfiguration) *InstallComponent {
	return &InstallComponent{
		cfg: cfg,
	}
}

func (ic *InstallComponent) InstallComponent(rls *resource.Release, namespace string) error {

	//vals, err := rls.InstallComponent(rls, ic.ResourcePath, c7nconsts.DefaultHelmValuesPath)
	//rls.Name = rls.Name + "-" + c7nutils.RandomString(5)
	return nil
}
