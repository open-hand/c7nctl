package resource

import (
	"testing"
)

func TestClient_GetResource(t *testing.T) {

	rc := NewClient(nil, "http://localhost:8080/")
	rc.Username = "admin"
	rc.Password = "123456"
	rc.Business = true

	result, _ := rc.GetResource("0.23", "/assets/os/install.yml")
	t.Logf("%s", result)
}

func TestClient_Login(t *testing.T) {

	rc := NewClient(nil, "http://localhost:8080/")
	rc.Username = "admin"
	rc.Password = "123456"
	rc.Business = true

	auth, _ := rc.Login("admin", "123456", true)
	t.Logf("%s", *auth.Data.Token)
}
