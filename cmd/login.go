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
	"github.com/choerodon/c7nctl/cmd/app"
	"github.com/spf13/cobra"
	"github.com/vinkdong/gox/log"
)

// installCmd represents the install command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to Choerodon",
	Long:  `Login to Choerodon.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			log.EnableDebug()
		}
		err := app.Login(cmd, args)
		if err != nil {
			log.Error("login failed")
		}
		log.Success("login succeed")
		return err
	},
}

func init() {
	loginCmd.Flags().Bool("debug", false, "enable debug output")
	loginCmd.Flags().String("name", "", "define the cluster name")
	rootCmd.AddCommand(loginCmd)
}
