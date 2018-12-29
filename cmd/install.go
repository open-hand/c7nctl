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
	"github.com/choerodon/c7n/cmd/app"
	"github.com/choerodon/c7n/pkg/common"
	"github.com/spf13/cobra"
	"github.com/vinkdong/gox/log"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Choerodon",
	Long:  `Install Choerodon quickly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			log.EnableDebug()
		}
		skip, _ := cmd.Flags().GetBool("skip-input")

		var (
			mail   string
			err    error
		)

		c, err := common.GetConfig()
		if err != nil {
			return err
		}
		if c.Terms.Accepted {
			goto start
		}

		if !skip {
			common.AskAgreeTerms()
			mail, err = common.AcceptUserInput(common.Input{
				Password: false,
				Tip:      "请输入您的邮箱以便通知您重要的更新(Please enter your email address):  ",
				Regex:    "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
			})
			if err != nil {
				return err
			}
			c.Terms.Accepted = true
			c.User.Mail = mail
			c.Write()
		} else {
			log.Info("your are execute job by skip input option, so we think you had allowed we collect your information")
		}

start:
		err = app.Install(cmd, args, mail)
		if err != nil {
			log.Error(err)
			log.Error("install failed")
		}
		log.Success("Install succeed")
		return nil
	},
}

var (
	ConfigFile   string
	ResourceFile string
)

func init() {
	installCmd.Flags().StringVarP(&ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")
	installCmd.Flags().StringVarP(&ConfigFile, "config-file", "c", "", "User Config file to read from, User define config by this file")
	installCmd.Flags().String("version", "", "specify a version")
	installCmd.Flags().Bool("debug", false, "enable debug output")
	installCmd.Flags().Bool("no-timeout", false, "disable install job timeout")
	installCmd.Flags().String("prefix", "", "add prefix to all helm release")
	installCmd.Flags().Bool("skip-input", false, "use default username and password to avoid user input")
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
