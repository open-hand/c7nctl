package app

import (
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/kube"
	c_kube "github.com/choerodon/c7n/cmd/kube"
	"fmt"
	"k8s.io/helm/pkg/helm"
	"io/ioutil"
	"github.com/vinkdong/gox/log"
	"os"
	"k8s.io/apimachinery/pkg/util/yaml"
	"encoding/json"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	tlsServerName string // overrides the server name used to verify the hostname on the returned certificates from the server.
	tlsCaCertFile string // path to TLS CA certificate file
	tlsCertFile   string // path to TLS certificate file
	tlsKeyFile    string // path to TLS key file
	tlsVerify     bool   // enable TLS and verify remote certificates
	tlsEnable     bool   // enable TLS

	tlsCaCertDefault = "$HELM_HOME/ca.pem"
	tlsCertDefault   = "$HELM_HOME/cert.pem"
	tlsKeyDefault    = "$HELM_HOME/key.pem"

	tillerTunnel *kube.Tunnel
	settings     helm_env.EnvSettings
	installFile  string
	installConfig = &InstallDefine{}
)

func PrepareEnv()  {
	
}

func setupConnection() {
	tunnel := c_kube.GetTunnel()
	settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", tunnel.Local)
	settings.TillerConnectionTimeout = 300
}

func newClient() *helm.Client {
	options := []helm.Option{helm.Host(settings.TillerHost), helm.ConnectTimeout(settings.TillerConnectionTimeout)}
	return helm.NewClient(options...)
}

func getInstallInfo() {
	var (
		data []byte
		err  error
	)
	if installFile != "" {
		data, err = ioutil.ReadFile(installFile)
		if err != nil {
			log.Error("read install file error")
			os.Exit(127)
		}
		data2, err := yaml.ToJSON(data)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(data2, installConfig)
	}
}

func getClusterMemoryAndCpu() (int64, int64)  {
	var sumMemory int64
	var sumCpu int64
	client := c_kube.GetClient()
	list, _ := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	for _,v := range list.Items{
		sumMemory += v.Status.Capacity.Memory().Value()
		sumCpu += v.Status.Capacity.Cpu().Value()
	}
	return sumMemory,sumCpu
}

func CheckResource(file string) bool {
	installFile = file
	getInstallInfo()
	request := installConfig.Spec.Resources.Requests
	reqMemory := request.Memory().Value()
	reqCpu := request.Cpu().Value()
	clusterMemory,clusterCpu := getClusterMemoryAndCpu()
	fmt.Print(reqMemory)
	if clusterMemory < reqMemory {
		log.Errorf("clusterMemory not Enough! require %dGi",reqMemory / (1024*1024*1024))
		return false
	}
	if clusterCpu < reqCpu {
		log.Errorf("clusterCpu not Enough! require %dc", reqCpu/1000)
		return false
	}
	return true
}