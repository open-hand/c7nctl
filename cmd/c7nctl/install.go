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
	"github.com/choerodon/c7nctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/vinkdong/gox/log"
	"io"
)

var (
	ConfigFile   string
	ResourceFile string
)

// installCmd represents the install command

func newInstallCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Choerodon",
		Long:  `Install Choerodon quickly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if debug, _ := cmd.Flags().GetBool("debug"); debug {
				log.EnableDebug()
			}
			skip, _ := cmd.Flags().GetBool("skip-input")

			var (
				mail string
				err  error
			)

			c, err := utils.GetConfig()
			if err != nil {
				return err
			}
			if c.Terms.Accepted {
				mail = c.OpsMail
				goto start
			}

			if !skip {
				utils.AskAgreeTerms()
				mail, err = utils.AcceptUserInput(utils.Input{
					Password: false,
					Tip:      "请输入您的邮箱以便通知您重要的更新(Please enter your email address):  ",
					Regex:    "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
				})
				if err != nil {
					return err
				}
				c.Terms.Accepted = true
				c.OpsMail = mail
				c.Write()
			} else {
				log.Info("your are execute job by skip input option, so we think you had allowed we collect your information")
			}

		start:
			err = app.Install(cmd, args, mail)
			if err != nil {
				log.Error("install failed")
				return err
			}
			log.Success("Install succeed")
			return nil
		},
	}

	cmd.Flags().StringVarP(&ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")
	cmd.Flags().StringVarP(&ConfigFile, "config-file", "c", "", "User Config file to read from, User define config by this file")
	cmd.Flags().String("version", "", "specify a version")
	cmd.Flags().Bool("debug", false, "enable debug output")
	cmd.Flags().Bool("no-timeout", false, "disable install job timeout")
	cmd.Flags().String("prefix", "", "add prefix to all helm release")
	cmd.Flags().Bool("skip-input", false, "use default username and password to avoid user input")

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
