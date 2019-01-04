package kube

import (
	"fmt"
	"github.com/vinkdong/gox/log"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetConfig(t *testing.T) {
	c, err := getConfig()
	if err != nil {
		log.Error(err)
	}
	log.Info(c.Host)
}

func TestGetClientset(t *testing.T) {
	//config,_ := getConfig()
	//_ ,client, _ := getClientset(config)
	client := GetClient()
	_, err := client.CoreV1().Nodes().List(meta_v1.ListOptions{})
	//for _, v := range list.Items {
	//	//fmt.Printf("node %s: %d \n",v.Name,v.Status.Capacity.Memory().Value())
	//	fmt.Printf("node %s: %d \n", v.Name, v.Status.Capacity.Cpu().Value())
	//}
	if err != nil {
		t.Error("get node failed")
	}
}
