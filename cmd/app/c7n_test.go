package app

import (
	"testing"
	"fmt"
	"github.com/choerodon/c7n/cmd/kube"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewClient(t *testing.T) {

	setupConnection()
	client := newClient()

	res, err :=  client.ListReleases()
	fmt.Println(err)
	for _,v := range res.Releases{
		fmt.Println(v.Name,v.Chart.Metadata.Name,v.Chart.Metadata.Version)
	}
}

func TestGetInstallInfo(t *testing.T)  {
	installFile = "/Users/vink/temp/install.yaml"
	getInstallInfo()
}

func TestGetNoInfo(t *testing.T){

	var sumMemory int64
	var sumCpu int64
	client := kube.GetClient()
	list, _ := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	for _,v := range list.Items{
		fmt.Printf("node %s: %d \n",v.Name,v.Status.Capacity.Memory().Value())
		sumMemory += v.Status.Capacity.Memory().Value()
		sumCpu += v.Status.Capacity.Cpu().Value()
	}
	//fmt.Print(sumMemory)
	fmt.Print(sumCpu)
}

func TestCheckResource(t *testing.T) {
	//fmt.Print(CheckResource())
}
