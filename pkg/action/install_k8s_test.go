package action

import (
	"os/exec"
	"testing"
)

func TestEnv(t *testing.T) {
	env := exec.Command("export", "ANSIBLE_HOST_KEY_CHECKING=False")
	err := env.Run()
	t.Log(err)
}
