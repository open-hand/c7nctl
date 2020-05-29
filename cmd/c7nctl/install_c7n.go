package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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
			err := installC7n(install)
			if err != nil {
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

func installC7n(install *action.InstallC7n) error {
	c, err := utils.GetConfig()
	if err != nil {
		cmdLog.Error(err)
		return err
	}
	// set user config. default is $HOME/.c7n/config.yaml
	setUserConfig(c, install)

	if err := install.Run(); err != nil {
		cmdLog.Error("InstallC7n failed")
		return err
	}
	cmdLog.Info("InstallC7n succeed")
	return nil
}

func addInstallFlags(fs *pflag.FlagSet, client *action.InstallC7n) {
	fs.StringVarP(&client.ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")
	fs.StringVarP(&client.ConfigFile, "--c7n-config", "c", "", "User Config file to read from, User define config by this file")
	fs.StringVar(&client.Version, "version", "", "specify a version")
	fs.BoolVar(&client.NoTimeout, "no-timeout", false, "disable resource job timeout")
	fs.StringVar(&client.Prefix, "prefix", "", "add prefix to all helm release")
	fs.BoolVar(&client.SkipInput, "skip-input", false, "use default username and password to avoid user input")
}

func setUserConfig(c *utils.Config, client *action.InstallC7n) {
	if !c.Terms.Accepted && !client.SkipInput {
		utils.AskAgreeTerms()
		mail, err := utils.AcceptUserInput(context.Input{
			Password: false,
			Tip:      "请输入您的邮箱以便通知您重要的更新(Please enter your email address):  ",
			Regex:    "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
		})
		if err != nil {
			log.Error(err)
		}
		c.Terms.Accepted = true
		c.OpsMail = mail

		_ = c.Write()
	} else {
		log.Info("your are execute job by skip input option, so we think you had allowed we collect your information")
	}
}
