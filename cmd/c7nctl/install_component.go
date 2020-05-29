package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const installComponentDesc = `
`

func newInstallComponentCmd(cfg *action.Configuration, args []string) *cobra.Command {
	install := action.NewInstall(cfg)
	cmd := &cobra.Command{
		Use:              "component [ARG]",
		Short:            "Install common components to k8s",
		Long:             installComponentDesc,
		Args:             minimumNArgs(1),
		PreRunE:          func(_ *cobra.Command, _ []string) error { return cfg.HelmClient.SetupConnection() },
		PersistentPreRun: func(*cobra.Command, []string) { cfg.HelmClient.InitSettings() },
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg.InitCfg()
			cmdLog.Info("Starting install component ", args[0])
			if err := install.InstallComponent(args[0]); err != nil {
				cmdLog.Error("Install component failed")
				return err
			}
			cmdLog.Info("Install component succeed")
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) { cfg.HelmClient.Teardown() },
	}

	settings := cfg.HelmClient.Settings()

	addInstallComponentFlags(cmd.Flags(), install)

	flags := cmd.PersistentFlags()
	settings.AddFlags(flags)
	_ = flags.Parse(args)

	// set defaults from environment
	settings.Init(flags)
	return cmd
}

func addInstallComponentFlags(fs *pflag.FlagSet, i *action.InstallC7n) {
	fs.StringVarP(&i.Namespace, "namespace", "n", "default", "Namespace Which installed component")
	fs.StringVarP(&i.ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")

}

// minimumNArgs returns an error if there is not at least N args.
func minimumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return errors.Errorf(
				"%q requires at least %d %s\n\nUsage:  %s",
				cmd.CommandPath(),
				n,
				pluralize("argument", n),
				cmd.UseLine(),
			)
		}
		return nil
	}
}

func pluralize(word string, n int) string {
	if n == 1 {
		return word
	}
	return word + "s"
}
