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
	"github.com/choerodon/c7n/pkg/c7nclient"
	"github.com/spf13/cobra"
)

var orgCode string
var proCode string

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.AddCommand(useOrgCmd)
	useCmd.AddCommand(useProCmd)

	useOrgCmd.Flags().StringVarP(&orgCode, "orgCode", "o", "", "org code")
	useProCmd.Flags().StringVarP(&proCode, "proCode", "p", "", "pro code")
}


// getCmd represents the get command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "The command to use organization or project",
	Long:  `you can use use command to define a default organization or a default project, then you can use other command with the default organization or the default project`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
	},
}

var useOrgCmd = &cobra.Command{
	Use:   "org",
	Short: "The command to use organization",
	Long:  `you can use use command to define a default organization ,then you can use other command with the default organization`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
	},
}


var useProCmd = &cobra.Command{
	Use:   "pro",
	Short: "The command to use project",
	Long:  `you can use use command to define a default project ,then you can use other command with the default project`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
	},
}

