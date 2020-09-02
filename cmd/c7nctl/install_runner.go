// Copyright © 2018 VinkDong <dong@wenqi.us>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/resource"
	"github.com/pkg/errors"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
	"io"
)

func NewGitlabRunnerCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	c := action.NewInstallRunner(cfg)

	cmd := &cobra.Command{
		Use:   "runner",
		Short: "config gitlab",
		Long:  `Config gitlab quickly.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := grcRun(c); err != nil {
				log.Error(err)
				log.Error("resource gitlab runner failed")
			}
			log.Info("config gitlab runner succeed")

		},
	}

	flags := cmd.Flags()
	addInstallRunnerFlags(flags, c)
	//grcAddFlags(flags)

	return cmd
}

func grcRun(client *action.InstallRunner) error {
	userConfig, err := getUserConfig(settings.ConfigFile)
	client.InitInstallRunner(userConfig)

	instDef := &resource.InstallDefinition{
		Version:     client.Version,
		PaaSVersion: client.Version,
	}

	if err = instDef.GetInstallDefinition(client.ResourcePath); err != nil {
		return std_errors.WithMessage(err, "Failed to get install configuration file")
	}
	instDef.MergerConfig(userConfig)

	stopCh := make(chan struct{})
	_, err = instDef.Spec.Basic.Slaver.InitSalver(client.GetClientSet(), settings.Namespace, stopCh)
	if err != nil {
		return errors.WithMessage(err, "Create Slaver failed")
	}
	defer func() {
		stopCh <- struct{}{}
	}()

	// 渲染 Release 和 runner
	client.RenderGitlabRunner(instDef, settings.Namespace)

	if err = client.InstallGitlabRunner(instDef, settings.Namespace); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Release %s install failed", instDef.Spec.Runner.Name))
	}
	return nil
}

func addInstallRunnerFlags(fs *pflag.FlagSet, client *action.InstallRunner) {
	fs.StringVarP(&client.ResourcePath, "resource-path", "r", "", "choerodon install definition file")
	fs.StringVarP(&client.Version, "version", "v", "0.22", "version of choerodon which will installation")

}
