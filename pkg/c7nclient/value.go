package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
	"io/ioutil"
)

func (c *C7NClient) ListValue(out io.Writer, envId int, valueDir string) {
	if c.currentContext.User.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	paras["env_id"] = envId
	paras["page"] = 1
	paras["size"] = 100000
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/deploy_value/page_by_options", c.currentContext.User.ProjectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var values = model.Values{}
	_, err = c.do(req, &values)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return

	}
	valuesInfos := []model.ValueInfo{}

	for _, value := range values.List {
		valueInfo := model.ValueInfo{
			Id:             value.Id,
			Name:           value.Name,
			Description:    value.Description,
			AppServiceName: value.Description,
			Environment:    value.EnvName,
		}

		err := ioutil.WriteFile(valueDir+value.Name, []byte(value.Value), 0644)
		if err != nil {
			fmt.Println(err)
		}
		valuesInfos = append(valuesInfos, valueInfo)
	}
	model.PrintValueInfo(valuesInfos, out)

}
