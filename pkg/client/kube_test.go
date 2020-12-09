package client

import (
	"path/filepath"
	"testing"
)

func TestK8sClient_ExecCommand(t *testing.T) {
	defaultKubeconfigPath := filepath.Join(homeDir(), ".kube", "config")

	cs, _ := GetKubeClient(defaultKubeconfigPath)
	kc := NewK8sClient(cs, "default")

	_ = kc.ExecCommand("busybox", "ping baidu.com")
}
