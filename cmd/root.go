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
	"bufio"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vinkdong/gox/log"
	"os"
	"regexp"
)

var cfgFile string

var clientPlatformConfig c7nclient.C7NConfig
var clientConfig c7nclient.C7NContext

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "c7nctl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
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
	initHosts()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.c7n.yaml)")
	rootCmd.PersistentFlags().StringVarP(&orgCode, "orgCode", "o", "", "org code")
	rootCmd.PersistentFlags().StringVarP(&proCode, "proCode", "p", "", "pro code")
	rootCmd.HasPersistentFlags()
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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
		//序列化配置文件为CONTEXT结构
		if err := viper.Unmarshal(&clientPlatformConfig); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			if clientPlatformConfig.CurrentContext == "" {
				fmt.Println(" You have to define current context!")
				os.Exit(1)
			}
			for _, context := range clientPlatformConfig.Contexts {
				if context.Name == clientPlatformConfig.CurrentContext {
					if context.Server == "" {
						fmt.Println(" You should define a server under the current context!")
						os.Exit(1)
					}
					clientConfig = context
				}
			}
			if clientConfig.Name == "" {
				fmt.Println(" The current context is not exist!")
				os.Exit(1)
			}
		}
	}
}

func DirectoryCheck(dirName string) {
	_, err := os.Stat(dirName)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func initHosts() {
	hosts := "/etc/hosts"
	file, err := os.OpenFile(hosts, os.O_RDWR, 0644)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		lineText := scan.Text()
		if isMatch, _ := regexp.Match("^199.232.28.133\\sraw.githubusercontent.com", []byte(lineText)); isMatch {
			log.Info("domain raw.githubusercontent.com existing in /etc/hosts")
			return
		}
	}
	// when raw.githubusercontent.com isn't in /etc/hosts, add it
	writer := bufio.NewWriter(file)
	if _, err = writer.WriteString("\n199.232.28.133\traw.githubusercontent.com\n"); err != nil {
		log.Error(err)
		os.Exit(0)
	}
	writer.Flush()
}
