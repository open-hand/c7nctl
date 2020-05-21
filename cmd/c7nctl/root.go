// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/spf13/cobra"
	"io"
	"os"
)

//TODO
var globalUsage = `The Choerodon Command Tool Line

Usage:
  c7nctl [command] [flags]

Available Commands:
  install	InstallC7n Choerodon in kubernetes with the given config
  config	Configuration choerodon Component. eg: gitlab runner
  upgrade	Upgrade Choerodon to newer version
  backup	Backup data Of Choerodon
  delete    Delete Choerodon or component.

Flags:
  -c, --config
  -h, --help		help fro c7nctl

Use "c7nctl [command] --help" for more information about a command.
`

var cfgFile string

// TODO move to Configuration

var (
	clientPlatformConfig c7nclient.C7NConfig
	clientConfig         c7nclient.C7NContext
)

func newRootCmd(actionConfig *action.Configuration, out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "c7nctl",
		Short: "The Choerodon Command Tool Line",
		Long:  globalUsage,
		// Uncomment the following line if your bare application
		// has an action associated with it:

		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	flags := cmd.PersistentFlags()
	envSettings.AddFlags(flags)

	// Add subcommand
	cmd.AddCommand(newConfigCmd(out))
	cmd.AddCommand(newCreateCmd(out))
	cmd.AddCommand(newDeleteCmd(out))
	cmd.AddCommand(newGetCmd(out))
	cmd.AddCommand(newInstallCmd(actionConfig, out, args))
	cmd.AddCommand(newLoginCmd(out))
	cmd.AddCommand(newLogoutCmd(out))
	cmd.AddCommand(newContextCmd(out))
	cmd.AddCommand(newUpgradeCmd(out))
	cmd.AddCommand(newUseCmd(out))

	return cmd
}

func DirectoryCheck(dirName string) {
	_, err := os.Stat(dirName)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
