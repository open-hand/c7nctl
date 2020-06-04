package utils

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/config"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

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

func GetInstallDefinition(file string, version string) (rd []byte) {
	if file != "" {
		url := fmt.Sprintf(githubResourceUrl, version, InstallConfigPath)
		rd = GetRemoteResource(url)
	} else {
		var err error
		rd, err = ioutil.ReadFile(file)
		CheckErrAndExit(err, 127)
	}
	return rd
}

func GetResourceFile(isRemote bool, version, filepath string) (rd []byte) {
	if isRemote {
		if version == "" {
			version = GetVersion("")
		}
		url := fmt.Sprintf(githubResourceUrl, version, filepath)
		rd = GetRemoteResource(url)
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
	vd := GetRemoteResource(url)
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

func GetClusterResource(client kubernetes.Interface) (int64, int64) {
	var sumMemory int64
	var sumCpu int64
	list, _ := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	for _, v := range list.Items {
		sumMemory += v.Status.Capacity.Memory().Value()
		sumCpu += v.Status.Capacity.Cpu().Value()
	}
	return sumMemory, sumCpu
}
