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
	"github.com/choerodon/c7n/pkg/c7nclient"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var appTemplateName string
var appTemplateCode string
var appTemplateDescription string
var clusterName string
var clusterCode string
var clusterDescription string
var copyFrom string
var appName string
var appType string
var appTemplate string
var envCode string
var envName string
var envDescription string
var instanceContent string

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createAppTemplateCmd)
	createCmd.AddCommand(createClusterCmd)
	createCmd.AddCommand(createAppCmd)
	createCmd.AddCommand(createEnvCmd)
	createCmd.AddCommand(createInstanceCmd)

	createAppTemplateCmd.Flags().StringVar(&appTemplateName, "name", "", "appTemplate name")
	createAppTemplateCmd.Flags().StringVar(&appTemplateCode, "code", "", "appTemplate code")
	createAppTemplateCmd.Flags().StringVar(&appTemplateDescription, "description", "", "appTemplate description")
	createAppTemplateCmd.Flags().StringVar(&copyFrom, "copyFrom", "", "appTemplate copy from")
	createClusterCmd.Flags().StringVar(&clusterName, "name", "", "cluster name")
	createClusterCmd.Flags().StringVar(&clusterCode, "code", "", "cluster code")
	createClusterCmd.Flags().StringVar(&clusterDescription, "description", "", "cluster description")
	createAppCmd.Flags().StringVar(&appName, "name", "", "app name")
	createAppCmd.Flags().StringVar(&appCode, "code", "", "app code")
	createAppCmd.Flags().StringVar(&appType, "type", "", "the value can be normal or test")
	createAppCmd.Flags().StringVar(&appTemplate, "appTemplate", "", "the appTemplate code you want to use")
	createEnvCmd.Flags().StringVar(&envName, "name", "", "env name")
	createEnvCmd.Flags().StringVar(&envCode, "code", "", "env code")
	createEnvCmd.Flags().StringVar(&envDescription, "description", "", "env Description ")
	createEnvCmd.Flags().StringVar(&clusterCode, "cluster", "", "the cluster code you want to use")
	createInstanceCmd.Flags().StringVar(&envCode, "env", "", "the envCode you want to deploy")
	createInstanceCmd.Flags().StringVar(&clusterCode, "cluster", "", "the clusterCode you want to deploy")
	createInstanceCmd.Flags().StringVar(&instanceContent, "content", "", "the values you want to deploy")
	createAppTemplateCmd.MarkFlagRequired("name")
	createAppTemplateCmd.MarkFlagRequired("code")
	createAppTemplateCmd.MarkFlagRequired("description")
	createClusterCmd.MarkFlagRequired("name")
	createClusterCmd.MarkFlagRequired("code")
	createClusterCmd.MarkFlagRequired("description")
	createAppCmd.MarkFlagRequired("name")
	createAppCmd.MarkFlagRequired("code")
	createAppCmd.MarkFlagRequired("type")
	createAppCmd.MarkFlagRequired("appTemplate")
	createEnvCmd.MarkFlagRequired("cluster")
	createEnvCmd.MarkFlagRequired("description")
	createEnvCmd.MarkFlagRequired("code")
	createEnvCmd.MarkFlagRequired("name")
	createInstanceCmd.MarkFlagRequired("env")
	createInstanceCmd.MarkFlagRequired("content")
	createInstanceCmd.MarkFlagRequired("cluster")
}

// getCmd represents the get command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "The command to create choerodon resource",
	Long:  `The command to create choerodon resource.such as organization, project, app, instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
	},
}

// create apptemplate command
var createAppTemplateCmd = &cobra.Command{
	Use:   "appTemplate",
	Short: "create application template",
	Long:  `you can use this command to create application template `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
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
		var appTemplateId int
		if copyFrom != "" {
			err, appTemplateInfo := c7nclient.Client.GetAppTemplate(cmd.OutOrStdout(), organizationId, copyFrom)
			if err != nil {
				return
			}
			appTemplateId = appTemplateInfo.ID
		}
		appTemplatePostInfo := model.AppTemplatePostInfo{appTemplateName, appTemplateCode, appTemplateDescription, appTemplateId}

		c7nclient.Client.CreateAppTemplate(cmd.OutOrStdout(), organizationId, &appTemplatePostInfo)
	},
}

// create cluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "create cluster",
	Long:  `you can use this command to create cluster `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
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
		clusterPostInfo := model.ClusterPostInfo{clusterName, clusterCode, clusterDescription, true}
		c7nclient.Client.CreateCluster(cmd.OutOrStdout(), organizationId, &clusterPostInfo)
	},
}

// create app command
var createAppCmd = &cobra.Command{
	Use:   "app",
	Short: "create app",
	Long:  `you can use this command to create app `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
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
		err, apptemplate := c7nclient.Client.GetAppTemplate(cmd.OutOrStdout(), pro.OrganizationID, appTemplate)
		if err != nil {
			return
		}
		appPostInfo := model.AppPostInfo{appName, appCode, appType, apptemplate.ID, true,}
		c7nclient.Client.CreateApp(cmd.OutOrStdout(), pro.ID, &appPostInfo)
	},
}

// create Env command
var createEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "create env",
	Long:  `you can use this command to create env `,
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)
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
		err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.OrganizationID, clusterCode)
		if !cluster.Connect {
			fmt.Println("the cluster you choose is not connected!")
			return
		}
		if err != nil {
			return
		}
		envPostInfo := model.EnvPostInfo{envName, envCode, envDescription, cluster.ID}
		c7nclient.Client.CreateEnv(cmd.OutOrStdout(), pro.ID, &envPostInfo)
	},
}

var createInstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "create instance",
	Long:  `you can use this command to create instance `,
	Run: func(cmd *cobra.Command, args []string) {

		c7nclient.InitClient(&clientConfig)

		if _, err := os.Stat(instanceContent); os.IsNotExist(err) {
			fmt.Println(err)
			return
		}
		b, err := ioutil.ReadFile(instanceContent)
		release := model.Release{}
		yaml.Unmarshal(b, &release)
		if err != nil {
			fmt.Print(err)
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
		err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.OrganizationID, clusterCode)
		if err != nil {
			return
		}
		err, app := c7nclient.Client.GetApp(release.Spec.ChartName, pro.ID)
		if err != nil {
			return
		}
		err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode, cluster.ID)
		if err != nil {
			return
		}
		err, version := c7nclient.Client.GetAppVersion(cmd.OutOrStdout(), pro.ID, release.Spec.ChartVersion, app.ID)
		if err != nil {
			return
		}
		instancePostInfo := model.InstancePostInfo{version.ID, env.ID, app.ID, release.Metadata.Name, release.Spec.Values, "create", false}
		c7nclient.Client.CreateInstance(cmd.OutOrStdout(), pro.ID, &instancePostInfo)
	},
}


//var createServiceCmd = &cobra.Command{
//	Use:   "service",
//	Short: "create servicese",
//	Long:  `you can use this command to create service `,
//	Run: func(cmd *cobra.Command, args []string) {
//
//		c7nclient.InitClient(&clientConfig)
//
//		if _, err := os.Stat(instanceContent); os.IsNotExist(err) {
//			fmt.Println(err)
//			return
//		}
//		b, err := ioutil.ReadFile(instanceContent)
//		release := kubeRelease{}
//		yaml.Unmarshal(b, &release)
//		if err != nil {
//			fmt.Print(err)
//			return
//		}
//		err, userInfo := c7nclient.Client.QuerySelf(cmd.OutOrStdout())
//		if err != nil {
//			return
//		}
//		err = c7nclient.Client.SetProject(cmd.OutOrStdout(), userInfo.ID)
//		if err != nil {
//			return
//		}
//		err, pro := c7nclient.Client.GetProject(cmd.OutOrStdout(), userInfo.ID, proCode)
//		if err != nil {
//			return
//		}
//		err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.OrganizationID, clusterCode)
//		if err != nil {
//			return
//		}
//		err, app := c7nclient.Client.GetApp(release.Spec.ChartName, pro.ID)
//		if err != nil {
//			return
//		}
//		err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, envCode, cluster.ID)
//		if err != nil {
//			return
//		}
//		err, version := c7nclient.Client.GetAppVersion(cmd.OutOrStdout(), pro.ID, release.Spec.ChartVersion, app.ID)
//		if err != nil {
//			return
//		}
//		instancePostInfo := model.InstancePostInfo{version.ID, env.ID, app.ID, release.Metadata.Name, release.Spec.Values, "create", false}
//		c7nclient.Client.CreateInstance(cmd.OutOrStdout(), pro.ID, &instancePostInfo)
//	},
//}
