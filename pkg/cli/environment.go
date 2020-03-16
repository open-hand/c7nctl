package cli

import "github.com/spf13/pflag"

type EnvSettings struct {
	CfgFile string
	OrgCode string
	ProCode string
}

func New() *EnvSettings {
	env := EnvSettings{}

	return &env
}

func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.CfgFile, "config", "", "config file (default is $HOME/.c7n.yaml)")
	fs.StringVarP(&s.OrgCode, "orgCode", "o", "", "org code")
	fs.StringVarP(&s.ProCode, "proCode", "p", "", "pro code")
	// 去掉了 toggle
}
