package helm

import (
	"github.com/choerodon/c7nctl/pkg/kube"
	"github.com/choerodon/c7nctl/pkg/utils"
	"testing"
)

var (
	testRepoURL = "https://openchart.choerodon.com.cn/choerodon/c7n/"
	testVersion = ""
	testChart   = "minio"
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
	_, err := client.locateChartPath(chartArgs)
	if err != nil {
		t.Error("locate chart failed")
	}
}

func TestInstallRelease(t *testing.T) {
	if !utils.ConditionSkip() {
		return
	}

	client := Client{
		Tunnel: kube.GetTunnel(),
	}
	client.InitClient()

	vals := []string{"pv.name=abc", "pbc[0].abc[1].als.s=rerws.sfds"}

	chartArgs := ChartArgs{
		ReleaseName: "",
		Namespace:   "",
		RepoUrl:     testRepoURL,
		Verify:      false,
		Version:     testVersion,
		ChartName:   testChart,
	}
	err := client.InstallRelease(vals, "", chartArgs)
	if err != nil {
		t.Fatal(err)
	}
}
