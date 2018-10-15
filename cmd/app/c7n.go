package app

import (
	"encoding/json"
	"fmt"
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/helm"
	"github.com/choerodon/c7n/pkg/install"
	kube2 "github.com/choerodon/c7n/pkg/kube"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vinkdong/gox/log"
	yaml_v2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/kube"
	"net/http"
	"os"
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

	tillerTunnel     *kube.Tunnel
	settings         helm_env.EnvSettings
	ResourceFile     string
	currentVersion   Version
	client           *helm.Client
	defaultNamespace = "choerodon"
	UserConfig       *config.Config
)

const (
	remoteConfigUrlPrefix = "http://share.hd.wenqi.us/install"
	versionPath           = "/version.yml"
	installConfigPath     = "/%s/install.yml"
	repoUrl               = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	C7nLabelKey           = "c7n-usage"
	C7nLabelValue         = "c7n-installer"
)

func getVersions() Versions {
	data := requireRemoteResource(versionPath)
	versions := Versions{}
	yaml_v2.Unmarshal(data, &versions)
	return versions
}

func getVersion(set *pflag.FlagSet) Version {
	versions := getVersions()
	//todo: select version
	return versions.GetLastStable()
}

func requireRemoteResource(resourcePath string) []byte {
	log.Infof("getting resource %s", resourcePath)
	var (
		data []byte
		err  error
	)
	resp, err := http.Get(fmt.Sprintf("%s%s", remoteConfigUrlPrefix, resourcePath))
	if err != nil {
		log.Error(err)
		os.Exit(127)
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Get resource %s failed", resourcePath)
		log.Error(err)
		os.Exit(127)
	}
	return data
}

func getInstallConfig() *install.Install {
	var (
		data []byte
		err  error
	)
	var install = &install.Install{}

	// request network resource
	if ResourceFile == "" {
		data = requireRemoteResource(fmt.Sprintf(installConfigPath, currentVersion.Version))

	}
	if ResourceFile != "" {
		data, err = ioutil.ReadFile(ResourceFile)
		if err != nil {
			log.Error("read install file error")
			os.Exit(127)
		}
	}
	data2, err := yaml.ToJSON(data)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data2, install)
	return install
}

func getUserConfig(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		os.Exit(124)
	}
	userConfig := &config.Config{}
	err = yaml_v2.Unmarshal(data, userConfig)
	if err != nil {
		log.Error(err)
		os.Exit(124)
	}
	UserConfig = userConfig

}
func TearDown() {
	tillerTunnel.Close()
}

func Install(cmd *cobra.Command) error {
	var err error
	// get current version to
	currentVersion = getVersion(cmd.Flags())

	// get install configFile
	ResourceFile, err = cmd.Flags().GetString("resource-file")
	if err != nil {
		log.Error(err)
		os.Exit(123)
	}

	configFile, err := cmd.Flags().GetString("config-file")
	getUserConfig(configFile)

	// get install config
	installConfig := getInstallConfig()

	if installConfig.Version == "" {
		log.Error("get install config error")
		os.Exit(127)
	}

	defer TearDown()
	//tunnel.Close()

	installConfig.UserConfig = UserConfig

	commonLabels := make(map[string]string)
	commonLabels[C7nLabelKey] = C7nLabelValue
	installConfig.CommonLabels = commonLabels

	// prepare environment
	tillerTunnel = kube2.GetTunnel()
	helmClient := &helm.Client{
		Tunnel: tillerTunnel,
	}
	helmClient.InitClient()
	installConfig.HelmClient = helmClient

	// do install
	return installConfig.Run()
}
