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
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var globalUsage = `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

var cfgFile string

var clientPlatformConfig c7nclient.C7NConfig
var clientConfig c7nclient.C7NContext

func newRootCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "c7nctl",
		Short: "The Choerodon manager command tools",
		Long:  globalUsage,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flags := cmd.PersistentFlags()
	settings.AddFlags(flags)

	// Add subcommand
	cmd.AddCommand(newConfigCmd(out))
	cmd.AddCommand(newCreateCmd(out))
	cmd.AddCommand(newDeleteCmd(out))
	cmd.AddCommand(newGetCmd(out))
	cmd.AddCommand(newInstallCmd(out))
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
