package app

import (
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/kube"
	c_kube "github.com/choerodon/c7n/cmd/kube"
	"fmt"
	"k8s.io/helm/pkg/helm"
	"io/ioutil"
	"github.com/containous/traefik/log"
	"os"
	"k8s.io/apimachinery/pkg/util/yaml"
	"encoding/json"
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
	client
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

func checkResource() bool {

	request := installConfig.Spec.Resources.Requests
	reqMemory := request.Memory().Value()
	reqCpu := request.Cpu().Value()

}