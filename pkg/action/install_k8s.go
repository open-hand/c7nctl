package action

const (
	kubeadmHaRepoUrl = "https://github.com/TimeBye/kubeadm-ha.git"
	repoPath         = "~/.c7n/cache/kubeadm-ha"
)

type InstallK8s struct {
	Cfg     Configuration
	Version string
}

func RunInstallK8s() error {

}

func cloneRepo() {

}
