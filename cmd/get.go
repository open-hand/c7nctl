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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
)

var appCode string
var generic string

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(envCmd)
	getCmd.AddCommand(instanceCmd)
	getCmd.AddCommand(proCmd)
	getCmd.AddCommand(orgCmd)
	getCmd.AddCommand(serviceCmd)
	getCmd.AddCommand(ingressCmd)
	getCmd.AddCommand(appVersionCmd)
	getCmd.AddCommand(appCmd)
	getCmd.AddCommand(clusterNodeCmd)
	getCmd.AddCommand(clusterCmd)
	getCmd.AddCommand(valueCmd)
	getCmd.AddCommand(certCmd)
	getCmd.AddCommand(configMapCmd)
	getCmd.AddCommand(secretCmd)
	getCmd.AddCommand(customCmd)
	getCmd.AddCommand(pvcCmd)
	getCmd.AddCommand(pvCmd)

	appVersionCmd.Flags().StringVarP(&appCode, "appCode", "c", "", "app code")
	instanceCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	instanceCmd.Flags().StringVarP(&clusterCode, "cluster", "", "", "cluster code")
	valueCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	serviceCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	ingressCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	ingressCmd.Flags().StringVarP(&clusterCode, "cluster", "c", "", "cluster code")
	clusterNodeCmd.Flags().StringVarP(&clusterCode, "clusterCode", "c", "", "cluster code")
	certCmd.Flags().StringVarP(&generic, "generic", "g", "n", "whether to list generic certification,\"y\" means generic certification,\"n\" means not generic certification")
	certCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	configMapCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	secretCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	customCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")
	pvcCmd.Flags().StringVarP(&envCode, "env", "e", "", "env code")

	instanceCmd.MarkFlagRequired("env")
	clusterNodeCmd.MarkFlagRequired("clusterCode")
	appVersionCmd.MarkFlagRequired("appCode")
	configMapCmd.MarkFlagRequired("env")
	valueCmd.MarkFlagRequired("env")
	certCmd.MarkFlagRequired("env")
	ingressCmd.MarkFlagRequired("env")
	customCmd.MarkFlagRequired("env")
	secretCmd.MarkFlagRequired("env")
	serviceCmd.MarkFlagRequired("env")
	pvcCmd.MarkFlagRequired("env")

}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "The command to get choerodon resource",
	Long:  `The command to get choerodon resource.such as organization, project, app, instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
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
	Short: "Get env",
	Long:  "Get env",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
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
		c7nclient.Client.ListEnvs(cmd.OutOrStdout(), pro.ID)
	},
}

// get orginazation command
var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Get organization",
	Long:  `List the organizations `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		c7nclient.Client.ListOrganization(cmd.OutOrStdout(), userInfo.ID)
	},
}

// get project command
var proCmd = &cobra.Command{
	Use:   "proj",
	Short: "Get project",
	Long:  `List the projects `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		c7nclient.Client.ListProject(cmd.OutOrStdout(), userInfo.ID)
	},
}

// get instance command
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Get an instance",
	Long:  `Get an instance in a specific environment`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
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

// get deploy value
var valueCmd = &cobra.Command{
	Use:   "value",
	Short: "Get deploy value",
	Long:  `Get deploy value in a specific environment`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
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
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		valueDir := home + "/c7nctl/value/"
		DirectoryCheck(valueDir)
		c7nclient.Client.ListValue(cmd.OutOrStdout(), env.ID, valueDir)
	},
}

// get application command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Get Application",
	Long:  `Get Devops Application List`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
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
		c7nclient.Client.ListApps(cmd.OutOrStdout(), pro.ID)
	},
}

// get application version command
var appVersionCmd = &cobra.Command{
	Use:   "app-version",
	Short: "Get Application version",
	Long:  `Get Devops Application Version List`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
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
		c7nclient.Client.ListAppVersions(cmd.OutOrStdout(), &appCode, pro.ID)
	},
}

// get cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Get Clusters",
	Long:  `Get Clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}

		err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userInfo.ID)
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
		c7nclient.Client.ListClusters(cmd.OutOrStdout(), pro.ID)
	},
}

// get cluster node command
var clusterNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Get Cluster Nodes",
	Long:  `Get Cluster Nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		error := c7nclient.Client.CheckIsLogin()
		if error != nil {
			fmt.Println(error)
			return
		}
		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
		if err != nil {
			return
		}
		err = c7nclient.Client.SetOrganization(cmd.OutOrStdout(), userInfo.ID)
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
		err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.ID, clusterCode)
		if err != nil {
			return
		}
		c7nclient.Client.ListClusterNode(cmd.OutOrStdout(), pro.ID, cluster.ID)
	},
}

// get service
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "get service",
	Long:  "Get service value in a specific environment",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		c7nclient.Client.ListService(cmd.OutOrStdout(), env.ID)
	},
}

// get ingress
var ingressCmd = &cobra.Command{
	Use:   "ingress",
	Short: "Get ingress",
	Long:  "Get ingress in a specific environment",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		c7nclient.Client.ListIngress(cmd.OutOrStdout(), env.ID)
	},
}

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Get certification",
	Long:  "Get cert in a specific environment",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		if generic == "y" {
			c7nclient.Client.ListGenericCert(cmd.OutOrStdout(), pro.ID)
		} else {
			c7nclient.Client.ListCert(cmd.OutOrStdout(), pro.ID, env.ID)
		}
	},
}

var configMapCmd = &cobra.Command{
	Use:   "cm",
	Short: "Get configMap",
	Long:  "Get configMap in a specific environment",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		c7nclient.Client.ListConfigMap(cmd.OutOrStdout(), pro.ID, env.ID)
	},
}

var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Get secret",
	Long:  "Get secret in a specific environment",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		c7nclient.Client.ListSecret(cmd.OutOrStdout(), pro.ID, env.ID)
	},
}
var customCmd = &cobra.Command{
	Use:   "custom",
	Short: "Get custom resource",
	Long:  "Get custom resource in a specific environment",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		c7nclient.Client.ListCustom(cmd.OutOrStdout(), pro.ID, env.ID)
	},}

var pvcCmd = &cobra.Command{
	Use:   "pvc",
	Short: "Get pvc",
	Long:  "Get pvc in a specific environment",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		c7nclient.Client.ListPvc(cmd.OutOrStdout(), pro.ID, env.ID)
	},}

var pvCmd = &cobra.Command{
	Use:   "pv",
	Short: "Get pv",
	Long:  "Get pv",
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig, &clientPlatformConfig)
		err := c7nclient.Client.CheckIsLogin()
		if err != nil {
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
		c7nclient.Client.ListPv(cmd.OutOrStdout(), pro.ID)
	},}
