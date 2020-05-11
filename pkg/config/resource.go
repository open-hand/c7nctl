package config

import (
	"github.com/vinkdong/gox/log"
	yaml_v2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/pkg/helm"
	"os"
)

const (
	versionPath       = "version.yml"
	installConfigPath = "install.yml"
	upgradeConfigPath = "upgrade.yml"
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
	Namespace    string
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
	versions := r.getVersions(version)
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

func (r *ResourceDefinition) getVersions(version string) Versions {
	data := r.requireRemoteResource(versionPath, version)
	versions := Versions{}
	yaml_v2.Unmarshal(data, &versions)
	return versions
}

func (r *ResourceDefinition) requireRemoteResource(resource string, version string) []byte {
	/*url := fmt.Sprintf(consts.RemoteInstallResourceRootUrl, version, resource)
	return utils.GetRemoteResource(url)*/
	return nil
}

func (r *ResourceDefinition) GetResourceData(resourceFile string, version string) (data []byte, err error) {
	if resourceFile == "" {
		// request network resource
		currentVersion := r.getVersion(version)
		data = r.requireRemoteResource(installConfigPath, currentVersion.Version)
	} else if resourceFile != "" {
		data, err = ioutil.ReadFile(resourceFile)
		if err != nil {
			log.Error("read resource file error")
			os.Exit(127)
		}
	}
	return data, err
}

func (r *ResourceDefinition) GetUpgradeResourceDate(version string) ([]byte, error) {
	// request network resource
	currentVersion := r.getVersion(version)
	var data []byte
	var err error
	if r.LocalFile == "" {
		data = r.requireRemoteResource(upgradeConfigPath, currentVersion.Version)
	}
	if r.LocalFile != "" {
		data, err = ioutil.ReadFile(r.LocalFile)
		if err != nil {
			log.Error("read resource file error")
			os.Exit(127)
		}
	}
	return data, err
}
