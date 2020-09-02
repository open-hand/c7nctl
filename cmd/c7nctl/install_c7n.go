package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/resource"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	yaml_v2 "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

const installC7nDesc = `One-click installation choerodon, When your prepared k8s, helm and NFS.
To install choerodon, you must set up the choerodon install configuration file
and specify the file with "--c7n-config <install-c7n-config.yaml>".

Ensure you run this within server can vista k8s.
`

func newInstallC7nCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	i := action.NewInstall(cfg)

	cmd := &cobra.Command{
		Use:   "c7n",
		Short: "One-click installation choerodon",
		Long:  installC7nDesc,
		Run: func(_ *cobra.Command, args []string) {
			setUserConfig(settings.SkipInput)
			metrics := client.Metrics{}
			if err := runInstallC7n(i, &metrics); err != nil {
				log.Errorf("InstallChoerodon Choerodon failed: %s", err)
				metrics.ErrorMsg = []string{err.Error()}
			} else {
				log.Info("InstallChoerodon Choerodon succeed")
			}
			metrics.Send()
		},
	}

	flags := cmd.PersistentFlags()
	addInstallFlags(flags, i)

	// set defaults from environment
	return cmd
}

func runInstallC7n(client *action.Install, metrics *client.Metrics) error {
	userConfig, err := getUserConfig(settings.ConfigFile)
	if err != nil {
		return err
	}
	client.InitInstall(userConfig)
	log.Infof("The current installing choerodon version is %s", client.Version)

	instDef := &resource.InstallDefinition{
		Version:     client.Version,
		PaaSVersion: client.Version,
	}

	if err = instDef.GetInstallDefinition(client.ResourcePath); err != nil {
		return std_errors.WithMessage(err, "Failed to get install configuration file")
	}
	instDef.MergerConfig(userConfig)

	// 检查资源，并将现有集群的硬件信息保存到 metrics
	if err := client.CheckResource(&instDef.Spec.Resources, metrics); err != nil {
		return err
	}
	if err := client.CheckNamespace(settings.Namespace); err != nil {
		return err
	}

	stopCh := make(chan struct{})
	if _, err = instDef.Spec.Basic.Slaver.InitSalver(client.GetClientSet(), settings.Namespace, stopCh); err != nil {
		return std_errors.WithMessage(err, "Create Slaver failed")
	}
	defer func() {
		stopCh <- struct{}{}
	}()

	// 渲染 Release
	if _, err := client.RenderChoerodon(instDef, settings.Namespace); err != nil {
		return err
	}

	// 安装 release
	if err := client.InstallChoerodon(instDef, settings.Namespace); err != nil {
		return err
	}

	// 清理历史的job
	// c.Clean()
	return nil
}

func addInstallFlags(fs *pflag.FlagSet, client *action.Install) {
	fs.StringVarP(&client.ResourcePath, "resource-path", "r", "", "choerodon install definition file")
	fs.StringVarP(&client.Version, "version", "v", "0.22", "version of choerodon which will installation")
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
	log.WithField("profile", filePath).Info("The user profile was read successfully")

	return userConfig, nil
}
