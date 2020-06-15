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

const deleteDesc = `Delete Choerodon quickly.`

// installCmd represents the resource command
func newDeleteCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete Choerodon",
		Long:  deleteDesc,
		RunE:  runDelete,
	}

	addDeleteFlags(cmd.Flags())

	return cmd
}

func runDelete(cmd *cobra.Command, args []string) error {
	/*	if debug, _ := cmd.Flags().GetBool("debug"); debug {
			log.EnableDebug()
		}
		err := app.Delete(cmd, args)
		if err != nil {
			log.Error(err)
			log.Error("delete failed")
		}*/
	return nil
}

func addDeleteFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&ResourceFile, "resource-file", "r", "", "resource file to read from, It provide which app should be installed")
	fs.StringVarP(&ConfigFile, "config-file", "c", "", "user Config file to read from, User define config by this file")
	fs.Bool("debug", false, "enable debug output")
	fs.StringP("namespace", "n", "", "select namespace")
}
