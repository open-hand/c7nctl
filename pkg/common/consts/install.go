package consts

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	// DefaultDatabaseUrlTpl 默认数据库连接信息模版
	DefaultDatabaseUrlTpl = "jdbc:mysql://%s:3306/%s?useUnicode=true&characterEncoding=utf-8&useSSL=false&useInformationSchema=true&remarks=true&allowMultiQueries=true&serverTimezone=Asia/Shanghai"
	// DefaultImageRepository 默认镜像仓库地址，因为同步 chart 包的时候会替换镜像仓库地址，所以不使用这个镜像地址
	DefaultImageRepository = "registry.cn-shanghai.aliyuncs.com/c7n/"
	// DefaultRepoUrl 默认 chart 仓库地址，当 install.yml 中没有定义时使用
	DefaultRepoUrl = "https://chart.choerodon.com.cn/choerodon/c7n-all"
	// Version DefaultRepoUrl = "http://chart.choerodon.com.cn/hand-yanqianyun/yq-devops"
	// Version c7nctl 版本号
	Version = "0.25"

	// OpenSourceResourceURL 默认的开源版和商业版资源获取路径
	OpenSourceResourceURL      = "https://gitee.com/open-hand/"
	OpenSourceResourceBasePath = "c7nctl/raw/%s/manifests"
	BusinessResourcePath       = "http://get.devops.hand-china.com/"
	BusinessResourceBasePath   = "assets/biz/%s/%s?token=%v"
	ResourceInstallFile        = "install.yml"
	// DefaultHelmValuesPath 默认 value.yaml 模版文件路径
	DefaultHelmValuesPath = "values"

	// C7nLabelKey 默认 label
	C7nLabelKey   = "c7n-usage"
	C7nLabelValue = "c7n-installer"

	MetricsUrl = "http://get.devops.hand-china.com/api/v1/metrics"
	IpAddr     = "ns1.dnspod.net:6666"
)

var (
	CommonLabels = map[string]string{
		C7nLabelKey: C7nLabelValue,
	}

	DefaultConfigPath     = filepath.Join(HomeDir(), ".c7n")
	DefaultConfigFileName = "config"
)

// 退出码
const (
	SuccessCode int = iota
	InitConfigErrorCode
)

// HomeDir returns the home directory for the current user
func HomeDir() string {
	if runtime.GOOS == "windows" {

		// First prefer the HOME environmental variable
		if home := os.Getenv("HOME"); len(home) > 0 {
			if _, err := os.Stat(home); err == nil {
				return home
			}
		}
		if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
			homeDir := homeDrive + homePath
			if _, err := os.Stat(homeDir); err == nil {
				return homeDir
			}
		}
		if userProfile := os.Getenv("USERPROFILE"); len(userProfile) > 0 {
			if _, err := os.Stat(userProfile); err == nil {
				return userProfile
			}
		}
	}
	return os.Getenv("HOME")
}
