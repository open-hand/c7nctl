package main

import (
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

func newInstallK8sCmd(out io.Writer, args []string) *cobra.Command {
	install := &action.InstallK8s{}
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "Simplest way to init your kubernets HA cluster",
		Long:  `c7nctl install k8s --master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 --node 192.168.0.5 --user root --passwd your-server-password`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := install.RunInstallK8s(); err != nil {
				log.Error(err)
			}
		},
	}
	addInstallK8sFlag(cmd.Flags(), install)
	return cmd
}

func addInstallK8sFlag(fs *pflag.FlagSet, install *action.InstallK8s) {
	fs.StringVar(&install.Ssh.AnsibleUser, "user", "root", "servers user name for ssh")
	fs.IntVar(&install.Ssh.AnsiblePort, "port", 22, "password for ssh")
	fs.StringVar(&install.Ssh.AnsiblePassword, "password", "", "password for ssh")

	fs.StringVar(&install.VIP, "vip", "", "virtual ip")
	fs.StringSliceVar(&install.MasterIPs, "master", []string{}, "kubernetes multi-masters")
	fs.StringSliceVar(&install.NodeIPs, "node", []string{}, "kubernetes multi-nodes")

	fs.StringVar(&install.Version, "version", "1.16.9", "version is kubernetes version")
	/*fs.StringVar(&install.Repo, "repo", "k8s.gcr.io", "choose a container registry to pull control plane images from")
	fs.StringVar(&install.PodCIDR, "podcidr", "100.64.0.0/10", "Specify range of IP addresses for the pod network")
	fs.StringVar(&install.SvcCIDR, "svccidr", "10.96.0.0/12", "Use alternative range of IP address for service VIPs")
	fs.StringVar(&install.Interface, "interface", "eth.*|en.*|em.*", "name of network interface")*/

	/*fs.BoolVar(&install.WithoutCNI, "without-cni", false, "If true we not install cni plugin")*/
	fs.StringVar(&install.Network, "network", "calico", "cni plugin, calico..")
	/*fs.BoolVar(&install.IPIP, "ipip", true, "ipip mode enable, calico..")
	fs.StringVar(&install.MTU, "mtu", "1440", "mtu of the ipip mode , calico..")

	fs.StringVar(&install.LvscareImage.Image, "lvscare-image", "fanux/lvscare", "lvscare image name")
	fs.StringVar(&install.LvscareImage.Tag, "lvscare-tag", "latest", "lvscare image tag name")

	fs.IntVar(&install.Vlog, "vlog", 0, "kubeadm log level")*/
}
