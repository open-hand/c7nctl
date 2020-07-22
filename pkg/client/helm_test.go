package client

import (
	"os"
	"testing"
)

func TestVals(t *testing.T) {
}

func TestRunHelmInstall(t *testing.T) {
	cfg := InitConfiguration("", "default")
	helmClient := NewHelm3Client(cfg)
	arg := ChartArgs{
		RepoUrl:     "http://chart.choerodon.com.cn/hzero/choerodon-ops",
		Version:     "5.0.4",
		Namespace:   "default",
		ReleaseName: "minio-test",
		Verify:      false,
		Keyring:     "",
		CertFile:    "",
		KeyFile:     "",
		CaFile:      "",
		ChartName:   "minio",
	}
	var vals map[string]interface{}
	if _, err := helmClient.Install(arg, vals, os.Stdout); err != nil {
		t.Error(err)
	}
}
