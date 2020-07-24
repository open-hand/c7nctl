package action

import (
	c7ncfg "github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/slaver"
)

type GitlabRunner struct {
	cfg        *C7nConfiguration
	Slaver     *slaver.Slaver
	UserConfig *c7ncfg.C7nConfig

	Version string

	// choerodon install configuration
	ConfigFile string
	// install resource
	ResourceFile string
	Namespace    string
}
