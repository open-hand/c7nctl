package client

import "testing"

func TestVals(t *testing.T) {
	valuesFile := `env:
  open:
    STORAGE: local
    DISABLE_API: false
    DEPTH: 2
persistence:
  enabled: true
  storageClass: test
ingress:
  enabled: true
  hosts:
  - path: /
`
	var values []string
	values = append(values, "ingress.hosts[0].name=admin")
	if v, err := vals(values, valuesFile); err != nil {
		t.Error(err)
	} else {
		t.Log(string(v))
	}
}

func TestRunHelmInstall(t *testing.T) {

}
