package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/resource"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const installComponentDesc = `
`

func newInstallComponentCmd(cfg *action.C7nConfiguration) *cobra.Command {
	c7n := action.NewInstallComponent(cfg)
	cmd := &cobra.Command{
		Use:   "component [ARG]",
		Short: "InstallChoerodon common components to k8s",
		Long:  installComponentDesc,
		Args:  minimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Starting install component ", args[0])
			if err := runInstallComponent(c7n, args[0]); err != nil {
				log.Error("InstallChoerodon component failed")
				return err
			}
			log.Info("InstallChoerodon component succeed")
			return nil
		},
	}

	return cmd
}

func runInstallComponent(client *action.InstallComponent, cname string) error {

	instDef := &resource.InstallDefinition{
		Version:     client.Version,
		PaaSVersion: client.Version,
	}

	if err := instDef.GetInstallDefinition(client.ResourcePath); err != nil {
		return std_errors.WithMessage(err, "Failed to get install configuration file")
	}

	for _, rls := range instDef.Spec.Component {
		if rls.Name == cname {
			err := instDef.RenderComponent(rls)
			if err != nil {
				return err
			}
			if err := client.InstallComponent(rls, settings.Namespace); err != nil {
				return err
			}
			return nil
		}
	}
	return std_errors.New("Please make sure the release name is right.")
}

// minimumNArgs returns an error if there is not at least N args.
func minimumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return std_errors.Errorf(
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
