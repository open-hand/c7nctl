package gitlab

import (
	"github.com/choerodon/c7n/pkg/install"
	"github.com/choerodon/c7n/pkg/kube"
)

const RunnerTokenLength = 30

type Runner struct {
	Namespace  string
	InstallDef *install.Install
}


func (runner *Runner) InstallRunner() error {

	runner.InstallDef.Client = kube.GetClient()
	i := runner.InstallDef

	ctx := install.Context{
		Client:       i.Client,
		Namespace:    i.Namespace,
		CommonLabels: i.CommonLabels,
		UserConfig:   i.UserConfig,
	}
	install.Ctx = ctx

	stopCh := make(chan struct{})
	s, err := runner.InstallDef.PrepareSlaver(stopCh)
	if err != nil {
		return err
	}

	install.Ctx.Slaver = s

	installDef := runner.InstallDef
	runnerConfig := installDef.Spec.Runner
	err = installDef.Install([]*install.InfraResource{runnerConfig})

	defer func() {
		stopCh <- struct{}{}
	}()
	return err
}
