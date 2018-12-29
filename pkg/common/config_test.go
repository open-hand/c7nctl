package common

import (
	"testing"
)

func TestGetConfigMail(t *testing.T) {
	c,err := GetConfig()
	if err != nil {
		t.Fatal(err)
	}
	c.User.Mail = "test@test.com"
	c.Write()
	d, _ := GetConfig()
	mail := d.User.Mail
	if mail != "test@test.com"{
		 t.Error("get config mail error")
	}
}