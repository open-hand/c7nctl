package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
)

const JoinCmdDesc = `
`

func newKubernetesJoinCmd(out io.Writer, args []string) *cobra.Command {
	install := &action.InstallK8s{}

	cmd := &cobra.Command{
		Use:   "join",
		Short: "Join node into kubernetes",
		Long:  JoinCmdDesc,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(install.MasterIPs) == 0 && len(install.NodeIPs) == 0 {
				log.Error("this command is join feature,master and node is empty at the same time.please check your args in command.")
				_ = cmd.Help()
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			install.RunJoinNode()
		},
	}
	addJoinFlag(cmd.Flags(), install)
	return cmd
}

func addJoinFlag(fs *pflag.FlagSet, install *action.InstallK8s) {
	fs.StringSliceVar(&install.MasterIPs, "master", []string{}, "kubernetes multi-master ex. 192.168.0.5-192.168.0.5")
	fs.StringSliceVar(&install.NodeIPs, "node", []string{}, "kubernetes multi-nodes ex. 192.168.0.5-192.168.0.5")

}
