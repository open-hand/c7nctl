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

var instanceId int
var appCode string

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(envCmd)
	getCmd.AddCommand(authEnvCmd)
	getCmd.AddCommand(instanceCmd)
	getCmd.AddCommand(instanceConfig)
	getCmd.AddCommand(proCmd)
	getCmd.AddCommand(orgCmd)
	getCmd.AddCommand(serviceCmd)
	getCmd.AddCommand(instanceResources)
	getCmd.AddCommand(ingressCmd)
	getCmd.AddCommand(appVersionCmd)
	getCmd.AddCommand(appTemplateCmd)
	getCmd.AddCommand(appCmd)
	getCmd.AddCommand(clusterNodeCmd)
	getCmd.AddCommand(clusterCmd)

	appVersionCmd.Flags().StringVarP(&appCode, "appCode", "a", "", "app code")
	instanceCmd.Flags().StringVarP(&envCode, "env", "", "", "env code")
	serviceCmd.Flags().StringVarP(&envCode, "env", "", "", "env code")
	ingressCmd.Flags().StringVarP(&envCode, "env", "", "", "env code")
	instanceCmd.Flags().StringVarP(&clusterCode, "cluster", "", "", "cluster code")
	serviceCmd.Flags().StringVarP(&clusterCode, "cluster", "", "", "cluster code")
	ingressCmd.Flags().StringVarP(&clusterCode, "cluster", "", "", "cluster code")
	clusterNodeCmd.Flags().StringVarP(&clusterCode, "clusterCode", "c", "", "cluster id")
	instanceConfig.Flags().IntVar(&instanceId, "instance-id", 0, "instance id")
	instanceResources.Flags().IntVar(&instanceId, "instance-id", 0, "instance id")
	clusterNodeCmd.MarkFlagRequired("clusterCode")
	appVersionCmd.MarkFlagRequired("appCode")
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "The command to get choerodon resource",
	Long:  `The command to get choerodon resource.such as organization, project, app, instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		if len(args) > 0 {
			fmt.Printf("don't have the resource %s, you can user c7nctl get --help to see the resource you can use!", args[0])
		} else {
			cmd.Help()
		}
	},
}

// get env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	Cobra is a CLI library for Go that empowers applications.
	This application`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userinfo.ID, proCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListEnvs(cmd.OutOrStdout(), pro.ID)
	},
}

// get orginazation command
var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "get organization",
	Long:  `list the organizations `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		c7nclient.Client.ListOrganization(cmd.OutOrStdout(), userinfo.ID)
	},
}

// get project command
var proCmd = &cobra.Command{
	Use:   "pro",
	Short: "get project",
	Long:  `list the projects `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		c7nclient.Client.ListProject(cmd.OutOrStdout(), userinfo.ID)
	},
}

// get auth env  command
var authEnvCmd = &cobra.Command{
	Use:   "authEnv",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userinfo.ID, proCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListAuthEnvs(cmd.OutOrStdout(), pro.ID)
	},
}

// get instance command
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
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
		if err != nil {
			return
		}
		err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListEnvsInstance(cmd.OutOrStdout(), env.ID)
	},
}

// get instance config  command
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
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		c7nclient.Client.InstanceConfig(cmd.OutOrStdout(), instanceId)
	},
}

// get instance resource command
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
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		c7nclient.Client.InstanceResources(cmd.OutOrStdout(), instanceId) // get app templates command

	},
}

// get application template command
var appTemplateCmd = &cobra.Command{
	Use:   "appTemplate",
	Short: "Get AppTemplate",
	Long:  `Get Devops App Templates List`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, organizationId := c7nclient.Client.GetOrganization(cmd.OutOrStdout(), userinfo.ID, orgCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListAppTemplates(cmd.OutOrStdout(), organizationId)
	},
}

// get application command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Get Application",
	Long:  `Get Devops Application List`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userinfo.ID, proCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListApps(cmd.OutOrStdout(), pro.ID)
	},
}

// get application version command
var appVersionCmd = &cobra.Command{
	Use:   "appVersion",
	Short: "Get Application version",
	Long:  `Get Devops Application Version List`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userinfo.ID, proCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListAppVersions(cmd.OutOrStdout(), &appCode, pro.ID)
	},
}

// get cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Get Clusters",
	Long:  `Get Clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, organizationId := c7nclient.Client.GetOrganization(cmd.OutOrStdout(), userinfo.ID, orgCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListClusters(cmd.OutOrStdout(), organizationId)
	},
}

// get cluster node command
var clusterNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Get Cluster Nodes",
	Long:  `Get Cluster Nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userinfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userinfo.ID)
		if err != nil {
			return
		}
		err, organizationId := c7nclient.Client.GetOrganization(cmd.OutOrStdout(), userinfo.ID, orgCode)
		if err != nil {
			return
		}
		err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), organizationId, clusterCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListClusterNode(cmd.OutOrStdout(), organizationId, cluster.ID)
	},
}

// getCmd represents the get command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "get env pipeline",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		//c7nclient.Client.ListService(cmd.OutOrStdout(), envId)
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
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		//c7nclient.Client.ListIngress(cmd.OutOrStdout(), envId)
	},
}
