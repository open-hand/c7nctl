package helm

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/kube"
	"testing"
)

var (
	testRepoURL = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	testVersion = ""
)

func TestLocateChartPath(t *testing.T) {
	client := Client{}
	client.InitClient()
	chartArgs := ChartArgs{
		ReleaseName: "",
		Namespace:   "",
		RepoUrl:     testRepoURL,
		Verify:      false,
		Version:     testVersion,
	}
	cp, err := client.locateChartPath(chartArgs)
	fmt.Println(cp, err)
}

func TestInstallRelease(t *testing.T) {
	client := Client{
		Tunnel: kube.GetTunnel(),
	}
	client.InitClient()

	vals := []string{"pv.name=abc"}

	chartArgs := ChartArgs{
		ReleaseName: "",
		Namespace:   "",
		RepoUrl:     testRepoURL,
		Verify:      false,
		Version:     testVersion,
	}
	err := client.InstallRelease(vals, chartArgs)
	if err != nil {
		t.Fatal(err)
	}
}
