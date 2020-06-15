package action

import (
	"github.com/choerodon/c7nctl/pkg/client"
	"os/exec"
	"testing"
)

func TestEnv(t *testing.T) {
	env := exec.Command("export", "ANSIBLE_HOST_KEY_CHECKING=False")
	err := env.Run()
	t.Log(err)
}

func TestInstallHelm(t *testing.T) {
	ssh := client.NewSSHClient("192.168.56.201", "root", "yishuida", 22)
	for _, cmd := range installHelmCmd {
		result, err := ssh.Run(cmd)
		if err != nil {
			t.Error(err)
		} else {
			t.Log(result)
		}
	}
}
