package resource

import (
	"testing"
)

func TestHelmValues(t *testing.T) {
	/*	infra := &Release{
			Values: []ChartValue{
				ChartValue{
					GetName:  "abc",
					Value: "",
					Input: context.Input{},
				},
				ChartValue{
					GetName:  "cde",
					Value: "cde",
					Case:  "{{ not .IgnorePv }}",
				},
			},
			GetName: "test-name-1",
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

func TestCleanJobs(t *testing.T) {
	/*	i := InstallChoerodon{
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
	/*
		param1 := ChartValue{
			GetName:  "name",
			Value: "value1",
		}
		param2 := ChartValue{
			GetName:  "name2",
			Value: "value5",
		}
		r := Request{
			Parameters: []ChartValue{param1, param2},
		}
		if r.parserParams() != "name=value1&name2=value5" {
			t.Error("request parames to params error")
		}

	*/
}

func TestRequestParserUrl(t *testing.T) {
	/*
		param1 := ChartValue{
			GetName:  "name",
			Value: "value1",
		}
		param2 := ChartValue{
			GetName:  "name2",
			Value: "value5",
		}
		r := Request{
			Parameters: []ChartValue{param1, param2},
		}
		fmt.Println(r.parserUrl())

	*/
}
