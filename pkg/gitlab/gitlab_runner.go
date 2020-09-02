package gitlab

import (
	"github.com/choerodon/c7nctl/pkg/resource"
)

const RunnerTokenLength = 30

type Runner struct {
	Namespace  string
	InstallDef *resource.InstallDefinition
}

func (runner *Runner) InstallRunner() error {

	/*	runner.InstallDef.Client = kube.GetClient()
		i := runner.InstallDef

		ctx := resource.Context{
			Client:       i.Client,
			Namespace:    i.Namespace,
			CommonLabels: i.CommonLabels,
			UserConfig:   i.UserConfig,
		}
		resource.Ctx = ctx

		stopCh := make(chan struct{})
		s, err := runner.InstallDef.PrepareSlaver(stopCh)
		if err != nil {
			return err
		}

		resource.Ctx.Slaver = s

		installDef := runner.InstallDef
		runnerConfig := installDef.Spec.Runner
		err = installDef.InstallChoerodon([]*resource.Release{runnerConfig})

		defer func() {
			stopCh <- struct{}{}
		}()
		return err*/
	return nil
}
