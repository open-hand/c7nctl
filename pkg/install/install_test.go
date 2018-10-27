package install

import (
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/kube"
	"github.com/vinkdong/gox/log"
	"testing"
)

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

func TestGetInfra(t *testing.T) {

	resource := make(map[string]*config.Resource)

	gitlabResource := &config.Resource{
		Host: "gitlab.example.io",
	}
	resource["gitlab"] = gitlabResource

	c := &config.Config{
		Spec: config.Spec{
			Resources: resource,
		},
	}
	Ctx.UserConfig = c
	Ctx.Client = kube.GetClient()
	Ctx.Namespace = ""

	preValue := PreValue{
		Name:  "GITLAB_BASE_DOMAIN",
		Value: "{{ ( .GetResource 'gitlab').Host }}",
		Check: "domain",
	}

	r := preValue.GetResource("gitlab")
	log.Info(r.Host)

}

func TestCleanJobs(t *testing.T) {
	i := Install{
		Client: kube.GetClient(),
		UserConfig: &config.Config{
			Metadata: config.Metadata{
				Namespace: "test",
			},
		},
	}
	i.CleanJobs()
}
