package config

import (
	"os"
	"io/ioutil"
	"fmt"
	"github.com/choerodon/c7n/pkg/helm"
	"github.com/spf13/pflag"
	"github.com/vinkdong/gox/log"
	yaml_v2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
)

const (
	remoteConfigUrlPrefix = "https://file.choerodon.com.cn/choerodon-install"
	versionPath           = "/version.yml"
	installConfigPath     = "/%s/install.yml"
	upgradeConfigPath     = "/%s/upgrade.yml"
)

type ResourceDefinition struct {
	LocalFile    string
	Version      string
	Metadata     Metadata
	Spec         Spec
	Client       kubernetes.Interface
	UserConfig   *Config
	HelmClient   *helm.Client
	CommonLabels map[string]string
	Timeout      int
}

type Versions struct {
	Versions []Version
}

type Version struct {
	Version string
	Status  string
}

func (v *Versions) GetLastStable() Version {
	for _, version := range v.Versions {
		if version.Status == "stable" {
			return version
		}
	}
	return Version{}
}

func (r *ResourceDefinition) getVersion(version string) Version {
	versions := r.getVersions()
	if version != "" {
		for _, v := range versions.Versions {
			if v.Version == version {
				return v
			}
		}
		log.Errorf("can't get version %s from remote server", version)
		os.Exit(1)
	}
	//todo: select version
	return versions.GetLastStable()
}

func (r *ResourceDefinition) getVersions() Versions {
	data := r.requireRemoteResource(versionPath)
	versions := Versions{}
	yaml_v2.Unmarshal(data, &versions)
	return versions
}

func (r *ResourceDefinition) requireRemoteResource(resourcePath string) []byte {
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

func (r *ResourceDefinition) GetResourceDate(version string) ([]byte, error) {
	// request network resource
	currentVersion := r.getVersion(version)
	var data []byte
	var err error
	if r.LocalFile == "" {
		data = r.requireRemoteResource(fmt.Sprintf(installConfigPath, currentVersion.Version))
	}
	if r.LocalFile != "" {
		data, err = ioutil.ReadFile(r.LocalFile)
		if err != nil {
			log.Error("read install file error")
			os.Exit(127)
		}
	}
	return data, err
}

func (r *ResourceDefinition) GetUpgradeResourceDate() ([]byte, error) {
	// request network resource
	var data []byte
	var err error
	if r.LocalFile == "" {
		currentVersion := r.getVersion(nil)
		data = r.requireRemoteResource(fmt.Sprintf(upgradeConfigPath, currentVersion.Version))
	}
	if r.LocalFile != "" {
		data, err = ioutil.ReadFile(r.LocalFile)
		if err != nil {
			log.Error("read install file error")
			os.Exit(127)
		}
	}
	return data, err
}
