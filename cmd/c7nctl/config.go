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
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

func newConfigCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "config Choerodon",
		Long:  `Config Choerodon quickly.`,
	}
	cmd.AddCommand(newGitlabCmd(out))

	return cmd
}

func newGitlabCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitlab",
		Short: "config gitlab",
		Long:  `Config gitlab quickly.`,
	}
	cmd.AddCommand(newGitlabRunnerCmd(out))
	return cmd
}

func newGitlabRunnerCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runner",
		Short: "config gitlab",
		Long:  `Config gitlab quickly.`,
		RunE:  grcRun,
	}

	flags := cmd.Flags()
	grcAddFlags(flags)

	return cmd
}

func grcRun(cmd *cobra.Command, args []string) error {
	/*	if debug, _ := cmd.Flags().GetBool("debug"); debug {
			log.EnableDebug()
		}
		// TODO
		installDef := app.GetInstall(nil, args)

		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}
		runner := gitlab.Runner{
			InstallDef: installDef,
			Namespace:  ns,
		}
		err = runner.InstallRunner()
		if err != nil {
			log.Error("resource gitlab runner failed")
			return err
		}
		log.Success("config gitlab runner succeed")*/
	return nil
}

func grcAddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")
	fs.StringVarP(&ConfigFile, "config-file", "c", "", "User Config file to read from, User define config by this file")
	fs.Bool("debug", false, "enable debug output")
	fs.StringP("namespace", "n", "c7n-system", "the namespace you resource choerodon")
	fs.String("prefix", "", "add prefix to all helm release")
	fs.String("version", "", "specify a version")
	fs.Bool("skip-input", false, "use default username and password to avoid user input")
}
