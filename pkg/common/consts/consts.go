package consts

import (
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	"path/filepath"
)

const (
	// 默认的获取远程资源文件的 url
	DefaultResource              = "https://cdn.jsdelivr.net/gh/yidaqiang/c7nctl@master/manifests/install.yaml"
	RemoteInstallResourceRootUrl = "https://cdn.jsdelivr.net/gh/yidaqiang/c7nctl@%s/manifests/values/%s.yaml"

	// 默认数据库连接信息
	DatabaseUrlTpl = "jdbc:mysql://%s:3306/%s?useUnicode=true&characterEncoding=utf-8&useSSL=false&useInformationSchema=true&remarks=true&allowMultiQueries=true&serverTimezone=Asia/Shanghai"

	// 默认的一些配置项
	DefaultImageRepository = "registry.cn-shanghai.aliyuncs.com/c7n/"
	DefaultRepoUrl         = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	DefaultHelmValuesPath  = "values"

	DefaultGitBranch = "master"
	C7nLabelKey      = "c7n-usage"
	C7nLabelValue    = "c7n-installer"

	Version = "0.22"

	ResourcePath    = "https://gitee.com/choerodon/c7nctl/raw/%s/manifests/"
	ImageRepository = "registry.cn-shanghai.aliyuncs.com/c7n"
	ChartRepository = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	DatasourceTpl   = "jdbc:mysql://%s:3306/%s?useUnicode=true&characterEncoding=utf-8&useSSL=false&useInformationSchema=true&remarks=true&allowMultiQueries=true&serverTimezone=Asia/Shanghai"
)

var (
	CommonLabels = map[string]string{
		C7nLabelKey: C7nLabelValue,
	}

	DefaultConfigPath     = filepath.Join(c7nutils.HomeDir(), ".c7n")
	DefaultConfigFileName = "config"
)

// 默认的资源文件名
const (
	VersionPath       = "version.yml"
	InstallConfigPath = "install.yml"
	UpgradeConfigPath = "upgrade.yml"
)

// TaskInfo 常量
const (
	StaticLogsCM        = "c7n-logs"
	StaticReleaseKey    = "release"
	StaticTaskKey       = "task"
	StaticPersistentKey = "persistent"
	PvType              = "pv"
	PvcType             = "pvc"
	CRDType             = "crd"
	ReleaseTYPE         = "helm"
	TaskType            = "task"
	UninitializedStatus = "uninitialized"
	SucceedStatus       = "succeed"
	FailedStatus        = "failed"
	InstalledStatus     = "installed"
	RenderedStatus      = "rendered"
	// if have after process while wait
	CreatedStatus      = "created"
	staticInstalledKey = "installed"
	staticExecutedKey  = "execute"
	SqlTask            = "sql"
	HttpGetTask        = "httpGet"
)

// 退出码
const (
	SuccessCode int = iota
	InitConfigErrorCode
)

// 服务列表
const (
	ChartMuseum          = "chartmuseum"
	Redis                = "c7n-redis"
	Mysql                = "c7n-mysql"
	Gitlab               = "gitlab"
	Harbor               = "harbor"
	Sonarqube            = "sonarqube"
	ChoerodonRegister    = "choerodon-register"
	ChoerodonPlatform    = "choerodon-platform"
	ChoerodonAdmin       = "choerodon-admin"
	ChoerodonIam         = "choerodon-iam"
	ChoerodonOauth       = "choerodon-Oauth"
	ChoerodonGateWay     = "choerodon-gateway"
	ChoerodonAsgard      = "choerodon-asgard"
	ChoerodonSwagger     = "choerodon-swagger"
	ChoerodonMessage     = "choerodon-message"
	ChoerodonMonitor     = "choerodon-monitor"
	ChoerodonFile        = "choerodon-file"
	DevopsService        = "devops-service"
	GitlabService        = "gitlab-service"
	WorkflowService      = "workflow-service"
	AgileService         = "agile-service"
	TestManagerService   = "test-manager-service"
	KnowledgebaseService = "knowledgebase-service"
	ChoerodonFrontHzero  = "choerodon-front-hzero"
	ChoerodonFront       = "choerodon-front"
)
