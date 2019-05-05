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
	"github.com/choerodon/c7n/cmd/app"
	"github.com/choerodon/c7n/pkg/c7nclient"
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


var username string
var password string
var url string

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	loginCmd.Flags().StringVarP(&username, "username", "U", "", "username")
	loginCmd.Flags().StringVarP(&password, "password", "P", "", "password")
	loginCmd.Flags().StringVarP(&url, "url", "", "", "");
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")
	loginCmd.MarkFlagRequired("url")
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
