// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

var envId int
var instanceId int

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(envCmd)
	getCmd.AddCommand(authEnvCmd)
	getCmd.AddCommand(instanceCmd)
	getCmd.AddCommand(instanceConfig)
	getCmd.AddCommand(ServiceCmd)
	getCmd.AddCommand(instanceResources)
	getCmd.AddCommand(ingressCmd)

	instanceCmd.Flags().IntVar(&envId, "env-id", 0, "env id")
	ServiceCmd.Flags().IntVar(&envId, "env-id", 0, "env id")
	ingressCmd.Flags().IntVar(&envId, "env-id", 0, "env id")
	instanceConfig.Flags().IntVar(&instanceId, "instance-id", 0, "instance id")
	instanceResources.Flags().IntVar(&instanceId, "instance-id", 0, "instance id")
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
	},
}



// getCmd represents the get command
var envCmd = &cobra.Command{
	Use:   "all-env",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.ListEnvs(cmd.OutOrStdout())
	},
}
// getCmd represents the get command
var authEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.ListAuthEnvs(cmd.OutOrStdout())
	},
}

// getCmd represents the get command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.ListEnvsInstance(cmd.OutOrStdout(), envId)
	},
}

// getCmd represents the get command
var instanceConfig = &cobra.Command{
	Use:   "instance-config",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.InstanceConfig(cmd.OutOrStdout(), instanceId)
	},
}


// getCmd represents the get command
var instanceResources = &cobra.Command{
	Use:   "instance-resources",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.InstanceResources(cmd.OutOrStdout(), instanceId)
	},
}


// getCmd represents the get command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.ListService(cmd.OutOrStdout(), envId)
	},
}

// getCmd represents the get command
var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		c7nclient.Client.ListIngress(cmd.OutOrStdout(), envId)
	},
}
