package install

import "testing"

func TestRenderValue(t *testing.T) {
	infra := &InfraResource{
		Persistence: []*Persistence{
			&Persistence{
				RefPvcName: "test-pvc-1",
			},
		},
		Name: "test-name-1",
	}

	tpl := "{{ (index .Persistence 0).RefPvcName }}"

	if val := infra.renderValue(tpl); val != "test-pvc-1" {
		t.Errorf("render template failed got %s", val)
	}
}

func TestHelmValues(t *testing.T) {
	infra := &InfraResource{
		Values: []ChartValue{
			ChartValue{
				Name:  "abc",
				Value: "",
				Input: Input{},
			},
		},
		Name: "test-name-1",
	}
	infra.HelmValues()
}
