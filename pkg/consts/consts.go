package consts

import (
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	"path/filepath"
)

const (
	RemoteInstallResourceRootUrl = "https://cdn.jsdelivr.net/gh/yidaqiang/c7nctl@%s/manifests/values/%s.yaml"
	// RemoteInstallResourceRootUrl = "https://file.choerodon.com.cn/choerodon-install"
	//RemoteInstallResourceRootUrl = "http://localhost/choerodon-install"
)

const (
	DefaultGitBranch = "master"
	DefaultRepoUrl   = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	C7nLabelKey      = "c7n-usage"
	C7nLabelValue    = "c7n-installer"
)

var (
	DefaultConfigPath     = filepath.Join(c7nutils.HomeDir(), ".c7n")
	DefaultConfigFileName = "config"
)

// c7nctl 默认信息
const (
	Version = 0.21
)

// 退出码
const (
	SuccessCode int = iota
	InitConfigErrorCode
)

// 服务列表
const (
	ChartMuseum   = "chartmuseum"
	Redis         = "redis"
	mysql         = "mysql"
	Gitlab        = "gitlab"
	Harbor        = "harbor"
	HzeroRegister = "hzero-register"
)
