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
	c7nconsts "github.com/choerodon/c7nctl/pkg/consts"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

func newConfigCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "config Choerodon",
		Long:  `Config Choerodon quickly.`,
	}
	cmd.AddCommand(newGitlabCmd(cfg, out))

	return cmd
}

func newGitlabCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitlab",
		Short: "config gitlab",
		Long:  `Config gitlab quickly.`,
	}
	cmd.AddCommand(newGitlabRunnerCmd(cfg, out))
	return cmd
}

func newGitlabRunnerCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	c := action.NewChoerodon(cfg, settings)

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
	grcAddFlags(flags)

	return cmd
}

func grcRun(c *action.Choerodon) error {
	id, err := c.GetInstallDef(settings.ConfigFile, settings.ResourceFile)
	if err != nil {
		return errors.WithMessage(err, "Failed to get install configration file")
	}
	if id.Spec.Basic.RepoURL != "" {
		c.RepoUrl = id.Spec.Basic.RepoURL
	} else {
		c.RepoUrl = c7nconsts.DefaultRepoUrl
	}
	c.DefaultAccessModes = id.DefaultAccessModes
	c.Slaver = &id.Spec.Basic.Slaver

	stopCh := make(chan struct{})
	_, err = c.PrepareSlaver(stopCh)
	if err != nil {
		return errors.WithMessage(err, "Create Slaver failed")
	}
	defer func() {
		stopCh <- struct{}{}
	}()

	// 渲染 Release 和 runner
	if err := c.RenderReleases(id); err != nil {
		return err
	}
	if err := c.RenderGitlabRunner(id); err != nil {
		return err
	}
	log.Infof("start install %s", id.Spec.Runner.Name)
	// 获取的 values.yaml 必须经过渲染，只能放在 id 中
	vals, err := id.RenderHelmValues(id.Spec.Runner, c.UserConfig)
	if err != nil {
		return err
	}
	if err = c.InstallRelease(id.Spec.Runner, vals); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Release %s install failed", id.Spec.Runner.Name))
	}
	return nil
}

func grcAddFlags(fs *pflag.FlagSet) {
	fs.Bool("debug", false, "enable debug output")
	fs.String("prefix", "", "add prefix to all helm release")
	fs.String("version", "", "specify a version")
	fs.Bool("skip-input", false, "use default username and password to avoid user input")
}
