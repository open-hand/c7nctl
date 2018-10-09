package app

import (
	"fmt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"github.com/choerodon/c7n/pkg/kube"
	"github.com/choerodon/c7n/pkg/helm"
)

func TestNewClient(t *testing.T) {

	tillerTunnel = kube.GetTunnel()
	helmClient := &helm.Client{
		Tunnel: tillerTunnel,
	}

	res, err := helmClient.Client.ListReleases()
	fmt.Println(err)
	for _, v := range res.Releases {
		fmt.Println(v.Name, v.Chart.Metadata.Name, v.Chart.Metadata.Version)
	}
}

func TestGetNoInfo(t *testing.T) {

	var sumMemory int64
	var sumCpu int64
	client := kube.GetClient()
	list, _ := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	for _, v := range list.Items {
		fmt.Printf("node %s: %d \n", v.Name, v.Status.Capacity.Memory().Value())
		sumMemory += v.Status.Capacity.Memory().Value()
		sumCpu += v.Status.Capacity.Cpu().Value()
	}
	//fmt.Print(sumMemory)
	fmt.Print(sumCpu)
}

func TestCheckResource(t *testing.T) {
	//fmt.Print(CheckResource())
}

func TestGetUserConfig(t *testing.T)  {
	getUserConfig("/Users/vink/go/src/github.com/choerodon/c7n/install.yaml")
	if UserConfig.Spec.Resources["mysql"].Host!= "192.168.12.88" {
		t.Error("read user config error")
	}
}