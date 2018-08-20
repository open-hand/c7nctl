package app

import (
	"testing"
	"fmt"
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
