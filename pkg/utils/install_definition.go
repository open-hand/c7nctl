package utils

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/context"
	v1 "k8s.io/api/core/v1"
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
	fileResource      = "https://file.choerodon.com.cn/choerodon-install/%s/%s"
)

type Versions struct {
	Versions []Version
}

type Version struct {
	Version string
	Status  string
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

func GetVersion(version string) string {
	if version == "" {
		version = "master"
	}
	url := fmt.Sprintf(githubResourceUrl, version, VersionPath)
	// TODO CheckNetWork error
	vd := GetRemoteResource(url)
	versions := Versions{}
	yaml_v2.Unmarshal(vd, &versions)

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

func CreateNamespace() bool {
	ns := &v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: context.Ctx.UserConfig.Metadata.Namespace,
		},
	}
	namespace, err := (*context.Ctx.KubeClient).CoreV1().Namespaces().Create(ns)
	log.Infof("creating namespace %s", namespace.Name)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
