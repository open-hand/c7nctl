package cli

import "github.com/spf13/pflag"

type EnvSettings struct {
	// choerodon configuration file
	CfgFile string
	// log level, default is info
	Debug bool
}

func New() *EnvSettings {
	return &EnvSettings{}
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.CfgFile, "config", "", "Specify the choerodon configuration file (default is $HOME/.c7n/config.yml)")
	fs.BoolVarP(&s.Debug, "debug", "d", false, "Set up the log level Whether it is debug")
}
