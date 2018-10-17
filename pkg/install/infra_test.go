package install

import (
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/kube"
	"github.com/vinkdong/gox/log"
	"testing"
)

func TestInfraResource_GetRequirement(t *testing.T) {
	r := make(map[string]*config.Resource)

	r["mysql"] = &config.Resource{
		Password: "abc123",
	}
	c := &config.Config{
		Spec: config.Spec{
			Resources: r,
		},
	}
	Ctx = Context{
		UserConfig: c,
	}
	infra := InfraResource{
		Requirements: []string{"mysql"},
		Values: []ChartValue{
			ChartValue{
				Name:  "abc",
				Value: `{{ .GetRequirement "mysql" "GITLAB_BASE_DOMAIN" }}`,
				Input: Input{},
			},
		},
	}
	result := infra.GetRequireResource("mysql")
	log.Info(result.Password)

	client := kube.GetClient()
	Ctx.Client = client
	Ctx.Namespace = "install"
	result2 := infra.GetRequireResource("mysql4")
	log.Info(result2.Password)
}
