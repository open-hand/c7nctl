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
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/docker/docker/pkg/fileutils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	clientPlatformConfig c7nclient.C7NConfig
	clientConfig         c7nclient.C7NContext

	settings = cli.New()
)

func main() {
	c7nCfg := new(action.C7nConfiguration)

	cmd := newRootCmd(c7nCfg, os.Stdout, os.Args[1:])
	cobra.OnInitialize(func() {
		initConfig()
		if settings.Debug {
			log.SetLevel(log.DebugLevel)
		}
		// 初始化 helm3Client
		cfg := client.InitConfiguration(settings.KubeConfig, settings.Namespace)
		c7nCfg.HelmClient = client.NewHelm3Client(cfg)
		// 初始化 kubeClient
		kubeclient, _ := client.GetKubeClient(settings.KubeConfig)
		c7nCfg.KubeClient = client.NewK8sClient(kubeclient)
	})
	if err := cmd.Execute(); err != nil {
		log.Debug(err)
	}
}

// 初始化 config 与 c7n api 操作有关
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// set default configuration is $HOME/.c7n/config.yml
	viper.AddConfigPath(consts.DefaultConfigPath)
	viper.SetConfigName(consts.DefaultConfigFileName)
	viper.SetConfigType("yaml")

	// read in environment variables that match
	viper.AutomaticEnv()

	viper.SetDefault("version", consts.Version)
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; Set default config to predefined path
			configPath := filepath.Join(consts.DefaultConfigPath, consts.DefaultConfigFileName+".yaml")
			if err = fileutils.CreateIfNotExists(configPath, false); err != nil {
				log.Debug(err)
			}
			log.Infof("Created default config file %s", file)
		} else {
			// Config file was found but another error was produced
			log.Error(err)
			os.Exit(consts.InitConfigErrorCode)
		}
	} else {
		var cfg config.Config
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Error(err)
			os.Exit(consts.InitConfigErrorCode)
		}
		// TODO 校验 c7n context 和 clientConfig.GetName 等是否存在
	}
}
