package consts

import (
	"os"
	"path/filepath"
	"runtime"
)

const (

	// 默认数据库连接信息
	DatabaseUrlTpl = "jdbc:mysql://%s:3306/%s?useUnicode=true&characterEncoding=utf-8&useSSL=false&useInformationSchema=true&remarks=true&allowMultiQueries=true&serverTimezone=Asia/Shanghai"

	// 默认的一些配置项
	DefaultImageRepository = "registry.cn-shanghai.aliyuncs.com/c7n/"
	DefaultRepoUrl         = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	DefaultHelmValuesPath  = "values"

	DefaultGitBranch = "master"
	C7nLabelKey      = "c7n-usage"
	C7nLabelValue    = "c7n-installer"

	Version = "0.23"

	OpenSourceResourceURL      = "https://gitee.com/open-hand/"
	OpenSourceResourceBasePath = "c7nctl/raw/%s/manifests/"
	BusinessResourcePath       = "http://get.devops.hand-china.com/"

	ResourceInstallFile = "install.yml"

	BusinessResourceBasePath = "assets/biz/%s/%s?token=%v"
	ImageRepository          = "registry.cn-shanghai.aliyuncs.com/c7n"
	ChartRepository          = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	DatasourceTpl            = "jdbc:mysql://%s:3306/%s?useUnicode=true&characterEncoding=utf-8&useSSL=false&useInformationSchema=true&remarks=true&allowMultiQueries=true&serverTimezone=Asia/Shanghai"
)

var (
	CommonLabels = map[string]string{
		C7nLabelKey: C7nLabelValue,
	}

	DefaultConfigPath     = filepath.Join(HomeDir(), ".c7n")
	DefaultConfigFileName = "config"

	DefaultGiteeAccessToken = "14b8f261fabe031456cd48e7d76d407a"
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
