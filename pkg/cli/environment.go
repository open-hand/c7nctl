package cli

import (
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
)

type EnvSettings struct {
	Namespace  string
	ConfigFile string
	KubeConfig string
	Debug      bool
	SkipInput  bool
	Timeout    int
}

func New() *EnvSettings {
	return &EnvSettings{
		// TODO complete env default setting
		Namespace: os.Getenv("C7N_NAMESPACE"),
	}
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	var defaultKubeconfigPath string
	if home := homeDir(); home != "" {
		defaultKubeconfigPath = filepath.Join(home, ".kube", "config")
	}
	fs.StringVarP(&s.Namespace, "namespace", "n", "c7n-system", "namespace scope for this request")
	fs.StringVarP(&s.ConfigFile, "config", "c", "config.yaml", "choerodon install configuration file")
	fs.StringVar(&s.KubeConfig, "kubeconfig", defaultKubeconfigPath, "(optional) absolute path to the kubeconfig file")
	fs.BoolVar(&s.Debug, "debug", false, "enable verbose output")
	fs.BoolVar(&s.SkipInput, "skip-input", false, "skip up unnecessary input")
	fs.IntVar(&s.Timeout, "timeout", 0, "the number of seconds the Operation has time out")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
