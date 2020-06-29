package cli

import (
	"github.com/spf13/pflag"
	"os"
)

type EnvSettings struct {
	Namespace    string
	Debug        bool
	ConfigFile   string
	ResourceFile string
}

func New() *EnvSettings {
	return &EnvSettings{
		Namespace: os.Getenv("C7N_NAMESPACE"),
	}
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.Namespace, "namespace", "n", "c7n-system", "namespace scope for this request")
	fs.BoolVar(&s.Debug, "debug", false, "enable verbose output")
	fs.StringVar(&s.ConfigFile, "config", "config.yaml", "choerodon configuration file")
	fs.StringVar(&s.ResourceFile, "resource", "install.yaml", "choerodon install definition file")
}
