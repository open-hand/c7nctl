package resource

import (
	"testing"
)

func TestInfraResource_GetRequirement(t *testing.T) {
	/*	r := make(map[string]*config.Resource)

		r["mysql"] = &config.Resource{
			Password: "abc123",
		}
		c := &config.Config{
			Spec: config.Spec{
				Resources: r,
			},
		}
		context.Ctx = context.Context{
			UserConfig: c,
		}
		infra := Release{
			Requirements: []string{"mysql"},
			Values: []ChartValue{
				ChartValue{
					Name:  "abc",
					Value: `{{ .GetRequirement "mysql" "GITLAB_BASE_DOMAIN" }}`,
					Input: context.Input{},
				},
			},
		}
		result := infra.GetResource("mysql")
		log.Info(result.Password)

		client := kube.GetClient()
		context.Ctx.KubeClient = client
		context.Ctx.Namespace = "resource"
		result2 := infra.GetResource("mysql4")
		log.Info(result2.Password)*/
}
