package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	c7nconsts "github.com/choerodon/c7nctl/pkg/common/consts"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const installComponentDesc = `
`

func newInstallComponentCmd(cfg *action.C7nConfiguration) *cobra.Command {
	c7n := action.NewChoerodon(cfg)
	cmd := &cobra.Command{
		Use:   "component [ARG]",
		Short: "Install common components to k8s",
		Long:  installComponentDesc,
		Args:  minimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Starting install component ", args[0])
			if err := InstallComponent(c7n, args[0]); err != nil {
				log.Error("Install component failed")
				return err
			}
			log.Info("Install component succeed")
			return nil
		},
	}

	return cmd
}

func InstallComponent(c *action.Choerodon, cname string) error {
	c.Namespace = settings.Namespace
	c.Version = c7nutils.GetVersion(c.Version)

	id, _ := c.GetInstallDef("", c7nconsts.DefaultResource)

	for _, rls := range id.Spec.Component {
		if rls.Name == cname {
			err := id.RenderComponent(rls)
			if err != nil {
				return err
			}
			vals, err := id.RenderHelmValues(rls, nil)
			rls.Name = rls.Name + "-" + c7nutils.RandomString(5)
			if err := c.InstallRelease(rls, vals); err != nil {
				return err
			} else {
				break
			}
		}
	}
	return nil
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
