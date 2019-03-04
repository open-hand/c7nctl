// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

var cfgFile string

var config c7nclient.C7NConfig

var clientConfig c7nclient.C7NPlatformContext

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "c7n",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		c7nclient.InitClient(&clientConfig)

		if _, err := os.Stat("test.yaml"); os.IsNotExist(err) {
			fmt.Println(err)
			return
		}
		b, err := ioutil.ReadFile("test.yaml")
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
		err, cluster := c7nclient.Client.GetCluster(cmd.OutOrStdout(), pro.OrganizationID, "testcobra4")
		if err != nil {
			return
		}
		err, app := c7nclient.Client.GetApp(release.Spec.ChartName, pro.ID)
		if err != nil {
			return
		}
		err, env := c7nclient.Client.GetEnv(cmd.OutOrStdout(), pro.ID, "zxzxzxzxz", cluster.ID)
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.c7n.yaml)")
	rootCmd.PersistentFlags().StringVarP(&orgCode, "orgCode", "o", "", "org code")
	rootCmd.PersistentFlags().StringVarP(&proCode, "proCode", "p", "", "pro code")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println(cfgFile)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".c7n" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".c7n")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	//序列化配置文件为CONTEXT结构
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(config.Contexts) == 0 {
		fmt.Println("No C7nConfig Context")
	}

	if config.CurrentContext == "" {
		clientConfig = config.Contexts[0].Context
		clientConfig.Name = config.Contexts[0].Name
	}
}
