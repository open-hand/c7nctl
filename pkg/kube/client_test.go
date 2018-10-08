package kube

import (
	"fmt"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetConfig(t *testing.T) {
	getConfig()
}

func TestGetClientset(t *testing.T) {
	//config,_ := getConfig()
	//_ ,client, _ := getClientset(config)
	client := GetClient()
	list, _ := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	for _, v := range list.Items {
		//fmt.Printf("node %s: %d \n",v.Name,v.Status.Capacity.Memory().Value())
		fmt.Printf("node %s: %d \n", v.Name, v.Status.Capacity.Cpu().Value())
	}
}
