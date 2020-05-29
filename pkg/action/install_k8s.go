package action

import (
	"bytes"
	"github.com/choerodon/c7nctl/pkg/consts"
	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/vinkdong/gox/log"
	"html/template"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)

const (
	kubeadmHaRepoUrl = "https://github.com/TimeBye/kubeadm-ha.git"
	repoPath         = ".c7n/cache/kubeadm-ha"
	installHelmCmd   = `kubectl create serviceaccount --namespace kube-system helm-tiller
kubectl create clusterrolebinding helm-tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:helm-tiller
curl -L -o /tmp/helm-v2.16.3-linux-amd64.tar.gz https://file.choerodon.com.cn/kubernetes-helm/v2.16.3/helm-v2.16.3-linux-amd64.tar.gz
tar -zxvf /tmp/helm-v2.16.3-linux-amd64.tar.gz
sudo mv /tmp/linux-amd64/helm /usr/bin/helm
helm init \
    --history-max=3 \
    --tiller-image=registry.aliyuncs.com/google_containers/tiller:v2.16.3 \
    --stable-repo-url=https://mirror.azure.cn/kubernetes/charts/ \
    --service-account=helm-tiller
`
)

type InstallK8s struct {
	Cfg       Configuration
	Ssh       SSHConfig
	Hosts     []string
	MasterIPs []string
	NodeIPs   []string
	Version   string
	Network   string
	VIP       string
}

type SSHConfig struct {
	SshPort  int
	Username string
	Password string
	PkFile   string
}

func (i InstallK8s) RunInstallK8s() error {
	// 检查输入参数是否有效
	i.checkValid()

	// clone repo kubeadm-ha
	cloneRepo()

	log.Info("starting write host.ini to file")
	home, _ := homedir.Dir()
	hostPath := home + string(os.PathSeparator) + ".c7n/host.ini"
	if checkFileIsExist(hostPath) {
		log.Info("host.ini already existing, skip up generate host.ini")
	} else {
		hostFile := i.renderHosts()

		var hostByte = []byte(hostFile)
		err := ioutil.WriteFile(hostPath, hostByte, 0644) //写入文件(字节数组)
		if err != nil {
			log.Error(err)
		}
	}

	log.Info("Starting install Necessary components")
	kubeadmHaPath := home + string(os.PathSeparator) + repoPath
	_, err1 := exec.LookPath("ansible")
	_, err2 := exec.LookPath("netaddr")
	if err1 == nil && err2 == nil {
		log.Error("command ansible and netaddr is Is already installed")
	} else {
		installAnsible := exec.Command(kubeadmHaPath + string(os.PathSeparator) + "install-ansible.sh")
		installAnsible.Stdout = os.Stdout
		if err := installAnsible.Run(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	log.Info("String install kubeadm-ha")
	intK8s := exec.Command("ansible-playbook", "-i"+hostPath, kubeadmHaPath+string(os.PathSeparator)+"90-init-cluster.yml")
	intK8s.Stdout = os.Stdout

	if err := intK8s.Run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info("sucessed install kubernetes")

	/*	log.Info("Starting install helm")
		ssh := client.NewSSHClient(i.MasterIPs[0], i.Ssh.Username, i.Ssh.Password, 22)
		result ,err := ssh.Run(installHelmCmd)
		if err != nil {
			log.Error(err)
		} else {
			log.Info(result)
		}*/
	return nil
}

func (i InstallK8s) checkValid() {
	// 检查输入的 IP 格式
	var hosts = append(i.MasterIPs, i.NodeIPs...)
	checkIP(hosts)
	if i.Ssh.Username == "" {
		log.Error("用户不能为空")
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
	if len(ips) == 0 {
		log.Error("节点数量不能为空")
		os.Exit(1)
	}
	for _, ip := range ips {
		if address := net.ParseIP(ip); address == nil {
			log.Errorf("IP %s 格式错误\n", ip)
			os.Exit(1)
		}
	}
	// TODO 检查重复 IP
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
