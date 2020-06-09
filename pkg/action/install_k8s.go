package action

import (
	"bytes"
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/consts"
	c7n_utils "github.com/choerodon/c7nctl/pkg/utils"
	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)

const (
	kubeadmHaRepoUrl = "https://github.com/TimeBye/kubeadm-ha.git"
	repoPath         = ".c7n/cache/kubeadm-ha"
)

var (
	home, _           = homedir.Dir()
	hostPath          = home + string(os.PathSeparator) + ".c7n/host.yaml"
	kubeadmHaPath     = home + string(os.PathSeparator) + repoPath
	installScriptPath = kubeadmHaPath + string(os.PathSeparator) + "install-ansible.sh"
	installPlaybook   = "90-init-cluster.yml"
	installHelmCmd    = [6]string{
		"kubectl create serviceaccount --namespace kube-system helm-tiller",
		"kubectl create clusterrolebinding helm-tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:helm-tiller",
		"curl -L -o /tmp/helm-v2.16.3-linux-amd64.tar.gz https://file.choerodon.com.cn/kubernetes-helm/v2.16.3/helm-v2.16.3-linux-amd64.tar.gz",
		"tar -zxvf /tmp/helm-v2.16.3-linux-amd64.tar.gz -c /tmp",
		"mv /tmp/linux-amd64/helm /usr/bin/helm",
		"helm init --history-max=3 --tiller-image=registry.aliyuncs.com/google_containers/tiller:v2.16.3 --stable-repo-url=https://mirror.azure.cn/kubernetes/charts/ --service-account=helm-tiller",
	}
)

type InstallK8s struct {
	Cfg       Configuration
	Ssh       client.Host
	Hosts     []string
	MasterIPs []string
	NodeIPs   []string
	Version   string
	Network   string
	VIP       string
}

func (i InstallK8s) RunInstallK8s() error {
	// 检查输入参数是否有效
	i.checkValid()

	// clone repo kubeadm-ha
	cloneRepo()
	inventory := i.newInventory()

	log.Info("starting write host.ini to file")

	if checkFileIsExist(hostPath) {
		log.Info("host.ini already existing, skip up generate host.ini")
	} else {
		hostByte, err := yaml.Marshal(inventory)
		c7n_utils.CheckErrAndExit(err, 1)
		err = ioutil.WriteFile(hostPath, hostByte, 0644) //写入文件(字节数组)
		c7n_utils.CheckErrAndExit(err, 1)
	}

	log.Info("Starting install Necessary components")
	_, err1 := exec.LookPath("ansible")
	_, err2 := exec.LookPath("netaddr")
	if err1 == nil && err2 == nil {
		log.Info("command ansible and netaddr is Is already installed")
	} else {
		installAnsible := exec.Command(installScriptPath)
		installAnsible.Stdout = os.Stdout
		if err := installAnsible.Run(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	log.Info("String install kubeadm-ha")
	execAnsiblePlaybook(installPlaybook)
	log.Info("sucessed install kubernetes")

	log.Info("Starting install helm")
	ssh := client.NewSSHClient(i.MasterIPs[0], i.Ssh.AnsibleUser, i.Ssh.AnsiblePassword, i.Ssh.AnsiblePort)
	for _, cmd := range installHelmCmd {
		result, err := ssh.Run(cmd)
		if err != nil {
			log.Error(result)
		} else {
			log.Info(result)
		}
	}

	return nil
}

func execAnsiblePlaybook(playbook string) {
	intK8s := exec.Command("ansible-playbook", "-i", hostPath, kubeadmHaPath+string(os.PathSeparator)+playbook)
	intK8s.Stdout = os.Stdout

	if err := intK8s.Run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func (i InstallK8s) newInventory() client.Inventory {

	if i.Version != "" {
		client.DefaultAnsibleVar.KubeVersion = i.Version
	}
	if i.Network != "" {
		client.DefaultAnsibleVar.NetworkPlugin = i.Network
	}
	inventory := client.Inventory{All: client.Group{
		Hosts: map[string]client.Host{},
		Vars:  client.DefaultAnsibleVar,
		Children: client.Children{
			Etcd:       client.Hostname{Hosts: map[string]interface{}{}},
			KubeMaster: client.Hostname{Hosts: map[string]interface{}{}},
			KubeWorker: client.Hostname{Hosts: map[string]interface{}{}},
			NewEtcd:    client.Hostname{Hosts: map[string]interface{}{}},
			NewMaster:  client.Hostname{Hosts: map[string]interface{}{}},
			NewWorker:  client.Hostname{Hosts: map[string]interface{}{}},
			Lb:         client.Hostname{Hosts: map[string]interface{}{}},
		},
	}}
	for _, ip := range i.MasterIPs {
		host := i.Ssh
		inventory.All.Hosts[ip] = host
		inventory.All.Children.KubeMaster.Hosts[ip] = struct{}{}
		inventory.All.Children.KubeWorker.Hosts[ip] = struct{}{}
		inventory.All.Children.Etcd.Hosts[ip] = struct{}{}
	}
	for _, ip := range i.NodeIPs {
		host := i.Ssh
		inventory.All.Hosts[ip] = host
		inventory.All.Children.KubeWorker.Hosts[ip] = struct{}{}
	}
	return inventory
}

// 检查输入项是否有效
func (i InstallK8s) checkValid() {
	// 检查输入的 IP 格式
	checkMasterIP(i.MasterIPs)
	var hosts = append(i.MasterIPs, i.NodeIPs...)
	checkIP(hosts)
}

func checkMasterIP(m []string) {
	if len(m) == 0 {
		log.Error("The number cant's be zero")
		os.Exit(1)
	}
	// 节点数必须为奇数
	if len(m)%2 == 0 {
		log.Error("The number of master nodes must be odd")
		os.Exit(1)
	}
}

func (i InstallK8s) renderHosts() string {
	tpl, err := template.New("hosts").Parse(consts.HostFile)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	var result bytes.Buffer
	err = tpl.Execute(&result, i)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	return result.String()
}

func checkIP(ips []string) {
	dict := make(map[string]bool)
	for _, ip := range ips {
		if address := net.ParseIP(ip); address == nil {
			log.WithField("IP", ip).Errorf("IP %s format error\n", ip)
			os.Exit(1)
		}
		if _, ok := dict[ip]; !ok {
			dict[ip] = true //不冲突, 主机名加入字典
		} else {
			log.WithField("IP", ip).Error("duplicate IP is not allowed")
			os.Exit(1)
		}
		log.WithField("IP", ip).Info("Node check ok")
	}
}

func cloneRepo() {
	home, _ := homedir.Dir()

	path := home + string(os.PathSeparator) + repoPath
	if checkFileIsExist(path) {
		log.Info("kubeadm-ha is existing!")
	} else {
		log.Info("git clone kubeadm-ha to .c7n/cache")
		_ = os.MkdirAll(path, os.ModePerm)
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      kubeadmHaRepoUrl,
			Progress: os.Stdout,
		})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
