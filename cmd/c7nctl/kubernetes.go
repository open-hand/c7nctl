package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

const kubernetesCmdDesc = ``

func newKubernetesCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "kubernetes Related operation.",
		Long:  kubernetesCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	addKubernetesFlag(cmd.Flags())
	cmd.AddCommand(
		newKubernetesJoinCmd(out, args),
		newKubernetesInitCmd(out, args),
	)

	return cmd
}

func addKubernetesFlag(flags *pflag.FlagSet) {

}
