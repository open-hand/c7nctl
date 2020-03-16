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
	"github.com/choerodon/c7nctl/cmd/c7nctl/app"
	"github.com/spf13/cobra"
	"github.com/vinkdong/gox/log"
	"io"
)

// installCmd represents the install command
func newDeleteCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete Choerodon",
		Long:  `Delete Choerodon quickly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if debug, _ := cmd.Flags().GetBool("debug"); debug {
				log.EnableDebug()
			}
			err := app.Delete(cmd, args)
			if err != nil {
				log.Error(err)
				log.Error("delete failed")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&ResourceFile, "resource-file", "r", "", "resource file to read from, It provide which app should be installed")
	cmd.Flags().StringVarP(&ConfigFile, "config-file", "c", "", "user Config file to read from, User define config by this file")
	cmd.Flags().Bool("debug", false, "enable debug output")
	cmd.Flags().StringP("namespace", "n", "", "select namespace")

	return cmd
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
