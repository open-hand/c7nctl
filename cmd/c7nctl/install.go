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

package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/resource"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pflag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	yaml_v2 "gopkg.in/yaml.v2"
	"helm.sh/helm/v3/cmd/helm/require"
	"io"
	"io/ioutil"
)

const installDesc = `
This command install a set of instances.

The install argument must be a install reference.

To specify configuration files or resource file path, use the '--config/--resource' flag and pass in a file.

	$ c7nctl install c7n -c config.yaml -r ./

To check the generated manifests of a release without installing the chart,
the '--debug' and '--client-only' flags can be combined.
`

// installCmd represents the resource command
func newInstallCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	client := action.NewInstall(cfg)
	metrics := c7nclient.Metrics{}

	cmd := &cobra.Command{
		Use:   "install [NAME] [flags]",
		Short: "One-click installation choerodon or other component",
		Long:  installDesc,
		Args:  require.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			setUserConfig(settings.SkipInput)
			if err := runInstall(args, client, out); err != nil {
				log.Errorf("Install Choerodon failed: %s", err)
				metrics.ErrorMsg = []string{err.Error()}
			} else {
				log.Info("Install Choerodon succeed")
			}
			metrics.Send()
		},
	}

	flags := cmd.PersistentFlags()
	addInstallFlags(flags, client)

	return cmd
}

func runInstall(args []string, client *action.Install, out io.Writer) error {
	name, err := client.GetName(args)
	if err != nil {
		return err
	}
	client.Name = name

	userConfig, err := getUserConfig(settings.ConfigFile)
	if err != nil {
		return err
	}
	client.InitInstall(userConfig)
	log.Infof("The current installing choerodon version is %s", client.Version)

	instDef := &resource.InstallDefinition{}
	if err = instDef.GetInstallDefinition(client.ResourcePath); err != nil {
		return std_errors.WithMessage(err, "Failed to get install configuration file")
	}
	if !instDef.IsName(name) {
		return std_errors.New("Please input right release name!")
	}
	instDef.MergerConfig(userConfig)
	client.Namespace = settings.Namespace
	return client.Run(instDef)
}

func addInstallFlags(fs *pflag.FlagSet, client *action.Install) {
	fs.StringVarP(&client.ResourcePath, "resource-path", "r", "", "choerodon install definition file")
	fs.StringVarP(&client.Version, "version", "v", "0.23", "version of choerodon which will installation")
	fs.StringVar(&client.Prefix, "prefix", "", "add prefix to all helm release")
	fs.StringVar(&client.ImageRepository, "image-repo", "", "default image repository of all release")
	fs.StringVar(&client.ChartRepository, "chart-repo", "", "chart repository url")
	fs.StringVar(&client.DatasourceTpl, "datasource-url", "", "datasource url template")

	fs.BoolVar(&client.ThinMode, "thin-mode", false, "install choerodon using Low resource consumption")
	fs.BoolVar(&client.ClientOnly, "client-only", false, "simulate an install")
}

func setUserConfig(skipInput bool) {
	// 在 c7nctl.initConfig() 中 viper 获取了默认的配置文件
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Debug(err)
	}
	if !skipInput {
		// 当用户没有接受条款时，让其输入
		if !cfg.Terms.Accepted {
			c7nutils.AskAgreeTerms()
			mail := inputUserMail()
			cfg.Terms.Accepted = true
			cfg.OpsMail = mail
			viper.Set("terms", cfg.Terms)
			viper.Set("opsMail", cfg.OpsMail)
			viper.WriteConfig()
		}
	} else {
		log.Info("your are execute job by skip input option, so we think you had allowed we collect your information")
	}
}

func inputUserMail() string {
	mail, err := c7nutils.AcceptUserInput(c7nutils.Input{
		Password: false,
		Tip:      "请输入您的邮箱以便通知您重要的更新(Please enter your email address):  ",
		Regex:    "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	})
	c7nutils.CheckErr(err)
	return mail
}

func getUserConfig(filePath string) (*config.C7nConfig, error) {
	// TODO 如果不需要 config.yaml
	if filePath == "" {
		return nil, std_errors.New("No user config defined by `-c`")
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, std_errors.WithMessage(err, "Read config file failed: ")
	}

	userConfig := &config.C7nConfig{}
	if err = yaml_v2.Unmarshal(data, userConfig); err != nil {
		return nil, std_errors.WithMessage(err, "Unmarshal config failed")
	}
	log.Infof("The user profile %s was read successfully", filePath)

	return userConfig, nil
}
