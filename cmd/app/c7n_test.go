package app

import (
	"github.com/choerodon/c7nctl/pkg/helm"
	"github.com/choerodon/c7nctl/pkg/kube"
	"testing"
)

func TestNewClient(t *testing.T) {

	tillerTunnel = kube.GetTunnel()
	helmClient := &helm.Client{
		Tunnel: tillerTunnel,
	}
	if tillerTunnel == nil {
		t.Log("skip...")
		return
	}
	helmClient.InitClient()
	_, err := helmClient.Client.ListReleases()
	if err != nil {
		t.Error("Test get client failed")
	}
}

func TestCheckResource(t *testing.T) {
	//fmt.Print(CheckResource())
}

func TestGetUserConfig(t *testing.T) {
	getUserConfig("~/go/src/github.com/choerodon/c7nctl/install.yaml")
	if UserConfig.Spec.Resources["mysql"].Host != "192.168.12.88" {
		t.Error("read user config error")
	}
}
