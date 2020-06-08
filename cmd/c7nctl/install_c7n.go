package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/context"
	c7n_utils "github.com/choerodon/c7nctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
)

const installC7nDesc = `One-click installation choerodon, When your prepared k8s, helm and NFS.
To install choerodon, you must set up the choerodon install configuration file
and specify the file with "--c7n-config <install-c7n-config.yaml>".

Ensure you run this within server can vista k8s.
`

func newInstallC7nCmd(cfg *action.Configuration, out io.Writer, args []string) *cobra.Command {
	install := action.NewInstall(cfg)

	cmd := &cobra.Command{
		Use:              "c7n",
		Short:            "One-click installation choerodon",
		Long:             installC7nDesc,
		PreRunE:          func(_ *cobra.Command, _ []string) error { return cfg.HelmClient.SetupConnection() },
		PersistentPreRun: func(*cobra.Command, []string) { cfg.HelmClient.InitSettings() },
		RunE: func(_ *cobra.Command, args []string) error {
			cfg.InitCfg()
			if err := installC7n(install); err != nil {
				return err
			}
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) { cfg.HelmClient.Teardown() },
	}

	settings := cfg.HelmClient.Settings()

	addInstallFlags(cmd.Flags(), install)

	flags := cmd.PersistentFlags()
	settings.AddFlags(flags)
	_ = flags.Parse(args)

	// set defaults from environment
	settings.Init(flags)

	return cmd
}

func installC7n(install *action.Choerodon) error {

	// set user config. default is $HOME/.c7n/config.yaml
	setUserConfig(install)

	if err := install.Run(); err != nil {
		log.Error("Choerodon failed")
		return err
	}
	log.Info("Choerodon succeed")

	return nil
}

func addInstallFlags(fs *pflag.FlagSet, client *action.Choerodon) {
	fs.StringVarP(&client.ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")
	fs.StringVarP(&client.ConfigFile, "c7n-config", "c", "", "User Config file to read from, User define config by this file")
	fs.StringVar(&client.Version, "version", "", "specify a version")
	fs.BoolVar(&client.NoTimeout, "no-timeout", false, "disable resource job timeout")
	fs.StringVar(&client.Prefix, "prefix", "", "add prefix to all helm release")
	fs.BoolVar(&client.SkipInput, "skip-input", false, "use default username and password to avoid user input")
}

func setUserConfig(client *action.Choerodon) {
	// 在 c7nctl.initConfig() 中 viper 获取了默认的配置文件
	c := config.Cfg
	if !c.Terms.Accepted && !client.SkipInput {
		c7n_utils.AskAgreeTerms()
		mail := inputUserMail()
		c.Terms.Accepted = true
		c.OpsMail = mail
		viper.Set("terms", c.Terms)
		viper.Set("opsMail", mail)
	} else {
		log.Info("your are execute job by skip input option, so we think you had allowed we collect your information")
	}
}

func inputUserMail() string {
	mail, err := c7n_utils.AcceptUserInput(context.Input{
		Password: false,
		Tip:      "请输入您的邮箱以便通知您重要的更新(Please enter your email address):  ",
		Regex:    "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	})
	c7n_utils.CheckErr(err)
	return mail
}
