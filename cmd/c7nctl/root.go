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
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/spf13/cobra"
	"io"
)

//TODO
var c7nctlDesc = `c7nctl is a powerful command line tool contains Choerodon related operations.

Complete sources is available at https://github.com/choerodon/c7nctl/.
`

func newRootCmd(actionConfig *action.Configuration, out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "c7nctl",
		Short: "The Choerodon Command Tool Line",
		Long:  c7nctlDesc,
	}

	flags := cmd.PersistentFlags()
	envSettings.AddFlags(flags)

	// Add sub command
	cmd.AddCommand(
		newConfigCmd(out),
		newCreateCmd(out),
		newDeleteCmd(out),
		newGetCmd(out),
		newInstallCmd(actionConfig, out, args),
		newLoginCmd(out),
		newLogoutCmd(out),
		newContextCmd(out),
		newUpgradeCmd(out),
		newUseCmd(out),
	)

	return cmd
}
