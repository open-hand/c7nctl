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
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/spf13/cobra"
)

var name string

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(contextCmd)

	contextCmd.Flags().StringVarP(&name, "name", "", "", "context name")
	contextCmd.MarkFlagRequired("name")
}

// installCmd represents the install command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to Choerodon",
	Long:  `Login to Choerodon.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.Login(cmd.OutOrStdout())

	},
}

// getCmd represents the get command

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "The command to logout choerodon",
	Long:  `you can use use command to logout choerodon , after you logout ,you can not use some c7n command,such as c7n create,c7n get.....`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		c7nclient.Client.Logout(cmd.OutOrStdout())
	},
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "The command to switch context",
	Long:  `you can use use command to switch context, after you swith ,the current context is changed!`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.SwitchContext(cmd.OutOrStdout(), name)
	},
}
