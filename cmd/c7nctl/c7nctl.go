// Copyright © 2018 choerodon <EMAIL ADDRESS>
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
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/c7nclient"
	"github.com/choerodon/c7nctl/pkg/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	clientPlatformConfig c7nclient.C7NConfig
	clientConfig         c7nclient.C7NContext

	envSettings = cli.New()
	cmdLog      = &log.Entry{}
)

func init() {
	cmdLog = log.New().WithFields(log.Fields{
		"pkg": "github.com/choerodon/c7nctl/cmd",
	})
}

func main() {
	actionConfig := action.NewCfg()

	cmd := newRootCmd(actionConfig, os.Stdout, os.Args[1:])
	cobra.OnInitialize(initConfig)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
	defer viper.WriteConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if envSettings.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if envSettings.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(envSettings.CfgFile)
	} else {
		// set default configuration is $HOME/.c7n/config.yml
		viper.AddConfigPath("$HOME/.c7n")
		viper.SetConfigType("yml")
		viper.SetConfigName("config")
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// TODO new don's used config.yml. so only check it existing ?
		if err := viper.Unmarshal(&clientPlatformConfig); err != nil {
			cmdLog.Error(err)
			os.Exit(1)
		} /*else {
			if clientPlatformConfig.CurrentContext == "" {
				cmdLog.Error(" You have to define current context!")
				os.Exit(1)
			}
			for _, context := range clientPlatformConfig.Contexts {
				if context.Name == clientPlatformConfig.CurrentContext {
					if context.Server == "" {
						cmdLog.Error(" You should define a server under the current context!")
						os.Exit(1)
					}
					clientConfig = context
				}
			}
			if clientConfig.Name == "" {
				log.Info(" The current context is not exist!")
				os.Exit(1)
			}
		}*/
	}
}
