package consts

const (
	RemoteInstallResourceRootUrl = "https://cdn.jsdelivr.net/gh/yidaqiang/c7nctl@%s/manifests/values/%s.yaml"
	// RemoteInstallResourceRootUrl = "https://file.choerodon.com.cn/choerodon-install"
	//RemoteInstallResourceRootUrl = "http://localhost/choerodon-install"
)

const (
	DefaultConfigFileName = "c7nctl"
	DefaultConfigPath     = "$HOME/.c7n"

	DefaultGitBranch = "master"
	DefaultRepoUrl   = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	C7nLabelKey      = "c7n-usage"
	C7nLabelValue    = "c7n-installer"
)
