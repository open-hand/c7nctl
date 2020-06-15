package client

import "testing"

func TestSsh(t *testing.T) {
	ssh := NewSSHClient("192.168.56.201", "root", "yishuida", 22)
	if err := ssh.connect(); err != nil {
		t.Error(err)
	}
	if result, err := ssh.Run("ls -al"); err != nil {
		t.Log(err)
	} else {
		t.Log(result)
	}
}
