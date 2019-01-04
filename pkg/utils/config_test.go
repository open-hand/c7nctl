package utils

import (
	"fmt"
	"testing"
)

func TestGetConfigMail(t *testing.T) {
	c, err := GetConfig()

	fmt.Println(c.Clusters)
	if err != nil {
		t.Fatal(err)
	}
	c.OpsMail = "test@test.com"
	c.Write()
	d, _ := GetConfig()
	mail := d.OpsMail
	if mail != "test@test.com" {
		t.Error("get config mail error")
	}
}
