package utils

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/pkg/errors"
	"github.com/vinkdong/gox/log"
	yaml_v2 "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	VersionPath       = "version.yml"
	InstallConfigPath = "install.yml"
	UpgradeConfigPath = "upgrade.yml"

	githubResourceUrl = "https://cdn.jsdelivr.net/gh/yidaqiang/c7nctl@%s/manifests/%s"
	fileResourceUrl   = "https://file.choerodon.com.cn/choerodon-install/%s/%s"
)

func GetInstallDefinition(file string, version string) (rd []byte, err error) {
	if file == "" {
		url := fmt.Sprintf(githubResourceUrl, version, InstallConfigPath)
		rd, err = GetRemoteResource(url)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to get install.yaml")
		}
	} else {
		rd, err = ioutil.ReadFile(file)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to Read install.yaml")
		}
	}
	log.Info("")
	return rd, nil
}

func GetResourceFile(isRemote bool, version, filepath string) (rd []byte) {
	if isRemote {
		if version == "" {
			version = GetVersion("")
		}
		url := fmt.Sprintf(githubResourceUrl, version, filepath)
		rd, _ = GetRemoteResource(url)
	} else {
		var err error
		// TODO resolve filepath separator error
		rd, err = ioutil.ReadFile(filepath)
		if err != nil {
			log.Error("read resource file error")
			os.Exit(127)
		}
	}

	return rd
}

// 获取最新的
func GetVersion(branch string) string {
	if branch == "" {
		branch = "master"
	}
	url := fmt.Sprintf(githubResourceUrl, branch, VersionPath)
	vd, _ := GetRemoteResource(url)
	versions := config.Versions{}
	err := yaml_v2.Unmarshal(vd, &versions)
	CheckErr(err)

	for _, v := range versions.Versions {
		if v.Status == "stable" {
			return v.Version
		}
	}
	return ""
}
