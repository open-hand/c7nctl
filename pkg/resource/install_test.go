package resource

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/kube"
	"github.com/vinkdong/gox/log"
	"testing"
)

func TestRenderValue(t *testing.T) {
	infra := &Release{
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
	/*	infra := &Release{
			Values: []ChartValue{
				ChartValue{
					Name:  "abc",
					Value: "",
					Input: context.Input{},
				},
				ChartValue{
					Name:  "cde",
					Value: "cde",
					Case:  "{{ not .IgnorePv }}",
				},
			},
			Name: "test-name-1",
		}
		context.Ctx = context.Context{
			UserConfig: &config.C7n_Config{
				Spec: config.Spec{
					Persistence: config.Persistence{
						StorageClassName: "",
					},
				},
			},
		}
		infra.HelmValues()*/
}

func TestGetInfra(t *testing.T) {

	resource := make(map[string]*config.Resource)

	gitlabResource := &config.Resource{
		Host: "gitlab.example.io",
	}
	resource["gitlab"] = gitlabResource

	c := &config.C7nConfig{
		Spec: config.Spec{
			Resources: resource,
		},
	}
	context.Ctx.UserConfig = c
	*context.Ctx.KubeClient = kube.GetClient()
	context.Ctx.Namespace = ""

	preValue := PreValue{
		Name:  "GITLAB_BASE_DOMAIN",
		Value: "{{ ( .GetResource 'gitlab').Host }}",
		Check: "domain",
	}

	r := preValue.GetResource("gitlab")
	log.Info(r.Host)

}

func TestCleanJobs(t *testing.T) {
	/*	i := Install{
			Client: kube.GetClient(),
			UserConfig: &config.C7n_Config{
				Metadata: config.Metadata{
					Namespace: "test",
				},
			},
		}
		i.CleanJobs()*/
}

func TestRequestParserParams(t *testing.T) {
	param1 := ChartValue{
		Name:  "name",
		Value: "value1",
	}
	param2 := ChartValue{
		Name:  "name2",
		Value: "value5",
	}
	r := Request{
		Parameters: []ChartValue{param1, param2},
	}
	if r.parserParams() != "name=value1&name2=value5" {
		t.Error("request parames to params error")
	}
}

func TestRequestParserUrl(t *testing.T) {
	param1 := ChartValue{
		Name:  "name",
		Value: "value1",
	}
	param2 := ChartValue{
		Name:  "name2",
		Value: "value5",
	}
	r := Request{
		Parameters: []ChartValue{param1, param2},
	}
	fmt.Println(r.parserUrl())
}
