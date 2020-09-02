package resource

import (
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
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

func TestRelease_String(t *testing.T) {
	releaseTest := []struct {
		Release
		result string
	}{
		{
			Release{
				Name: "test1",
				PreInstall: []ReleaseJob{
					{
						Name:     "test-job1",
						InfraRef: "haha",
						Database: "",
						Commands: nil,
						Mysql:    []string{"select * from test"},
						Psql:     nil,
						Opens:    nil,
						Request:  nil,
					},
				},
			},
			"{Name: \"test1\"}",
		},
	}

	for _, r := range releaseTest {
		c7nutils.PrettyPrint(r)
	}
}
