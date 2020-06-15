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
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/spf13/cobra"
	"io"
)

var orgCode string
var proCode string
var defaultValue string

// getCmd represents the get command
func newUseCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "The command to use organization or project",
		Long:  `you can use use command to define a default organization or a default project, then you can use other command with the default organization or the default project`,
		Run: func(cmd *cobra.Command, args []string) {
			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			error := c7nclient.Client.CheckIsLogin()
			if error != nil {
				fmt.Println(error)
				return
			}
			if len(args) > 0 {
				fmt.Printf("don't have the resource %s, you can user c7nctl use --help to see the resource you can use!", args[0])
			} else {
				cmd.Help()
			}
		},
	}

	cmd.AddCommand(newUseOrgCmd(out))
	cmd.AddCommand(newUseProCmd(out))

	return cmd
}

func newUseOrgCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "The command to use organization",
		Long:  `you can use use command to define a default organization ,then you can use other command with the default organization`,
		Run: func(cmd *cobra.Command, args []string) {
			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			error := c7nclient.Client.CheckIsLogin()
			if error != nil {
				fmt.Println(error)
				return
			}
			error, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if error != nil {
				return
			}
			error = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userInfo.ID)
			if error != nil {
				return
			}
			c7nclient.Client.UseOrganization(cmd.OutOrStdout(), orgCode)
		},
	}
	cmd.Flags().StringVarP(&orgCode, "orgCode", "o", "", "org code")
	cmd.Flags().StringVarP(&defaultValue, "default", "", "", "default")
	cmd.MarkFlagRequired("orgCode")

	return cmd
}

func newUseProCmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proj",
		Short: "The command to use project",
		Long:  `you can use use command to define a default project ,then you can use other command with the default project`,
		Run: func(cmd *cobra.Command, args []string) {
			c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
			error := c7nclient.Client.CheckIsLogin()
			if error != nil {
				fmt.Println(error)
				return
			}
			error, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
			if error != nil {
				return
			}
			error = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
			if error != nil {
				return
			}
			c7nclient.Client.UseProject(cmd.OutOrStdout(), proCode)
		},
	}

	cmd.Flags().StringVarP(&proCode, "proCode", "p", "", "pro code")
	cmd.Flags().StringVarP(&defaultValue, "default", "", "", "default")

	cmd.MarkFlagRequired("proCode")

	return cmd
}
