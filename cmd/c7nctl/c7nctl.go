// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"github.com/choerodon/c7nctl/pkg/cli"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	settings = cli.New()
)

func init() {
	log.SetPrefix("[c7nctl] ")
	log.SetFlags(log.Lshortfile)
}

func main() {
	cmd := newRootCmd(os.Stdout, os.Args[1:])
	cobra.OnInitialize(initConfig)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if settings.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(settings.CfgFile)
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
