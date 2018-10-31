package gitlab

import (
	"time"
	"github.com/vinkdong/gox/random"
	"math/rand"
	"github.com/choerodon/c7n/pkg/install"
	"github.com/choerodon/c7n/pkg/kube"
)

const RunnerTokenLength = 30

type Runner struct {
	Namespace  string
	InstallDef *install.Install
}

func (runner *Runner) GenerateRunnerToken() string {
	bytes := make([]byte, RunnerTokenLength)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < RunnerTokenLength; i++ {
		random.Seed(time.Now().UnixNano())
		op := random.RangeIntInclude(random.Slice{Start: 48, End: 57},
			random.Slice{Start: 97, End: 122})
		bytes[i] = byte(op) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func (runner *Runner) InstallRunner() error {

	runner.InstallDef.Client = kube.GetClient()
	i := runner.InstallDef



	ctx := install.Context{
		Client:       i.Client,
		Namespace:    runner.Namespace,
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
