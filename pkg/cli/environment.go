package cli

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

type EnvSettings struct {
	Namespace    string
	Debug        bool
	ConfigFile   string
	ResourceFile string
	KubeConfig   string
}

func New() *EnvSettings {
	return &EnvSettings{
		Namespace: os.Getenv("C7N_NAMESPACE"),
	}
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	var defaultKubeconfigPath string
	if home := homeDir(); home != "" {
		defaultKubeconfigPath = filepath.Join(home, ".kube", "config")
	}
	fs.StringVarP(&s.Namespace, "namespace", "n", "c7n-system", "namespace scope for this request")
	fs.BoolVar(&s.Debug, "debug", false, "enable verbose output")
	fs.StringVarP(&s.ConfigFile, "config", "c", "config.yaml", "choerodon configuration file")
	fs.StringVarP(&s.ResourceFile, "resource", "r", "", "choerodon install definition file")

	fs.StringVar(&s.KubeConfig, "kubeconfig", defaultKubeconfigPath, "(optional) absolute path to the kubeconfig file")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
