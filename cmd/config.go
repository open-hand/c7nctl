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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vinkdong/gox/log"
	"github.com/choerodon/c7n/pkg/gitlab"
	"github.com/choerodon/c7n/cmd/app"
)

// installCmd represents the install command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config Choerodon",
	Long:  `Config Choerodon quickly.`,
}

var gitlabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "config gitlab",
	Long:  `Config gitlab quickly.`,
}

var gitlabRunnerCmd = &cobra.Command{
	Use:   "runner",
	Short: "config gitlab",
	Long:  `Config gitlab quickly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			log.EnableDebug()
		}
		installDef := app.GetInstall(cmd, args)

		ns ,err := cmd.Flags().GetString("namespace")
		if err!=nil {
			return err
		}
		runner := gitlab.Runner{
			InstallDef: installDef,
			Namespace: ns,
		}
		err = runner.InstallRunner()
		if err != nil {
			log.Error(err)
			log.Error("install failed")
		}
		log.Success("config gitlab succeed")
		return nil
	},
}

func init() {

	gitlabRunnerCmd.Flags().StringVarP(&ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")
	gitlabRunnerCmd.Flags().StringVarP(&ConfigFile, "config-file", "c", "", "User Config file to read from, User define config by this file")
	gitlabRunnerCmd.Flags().Bool("debug", false, "enable debug output")
	gitlabRunnerCmd.Flags().StringP("namespace","n","c7n-system","the namespace you install choerodon")
	gitlabRunnerCmd.Flags().String("prefix","","add prefix to all helm release")
	gitlabCmd.AddCommand(gitlabRunnerCmd)
	configCmd.AddCommand(gitlabCmd)
	rootCmd.AddCommand(configCmd)
}
