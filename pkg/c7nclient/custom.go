package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
	"net/url"
	"time"
)

func (c *C7NClient) ListCustom(out io.Writer, projectId int, envId int) {
	if projectId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = 1
	paras["size"] = 10000

	body := make(map[string]interface{})
	body["param"] = ""
	body["searchParam"] = make(map[string]string)
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/customize_resource/%d/page_by_env", projectId, envId), paras, body)
	if err != nil {
		fmt.Printf("build request error")

	}
	var customs = model.Customs{}
	_, err = c.do(req, &customs)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	now := time.Now()
	loc, _ := time.LoadLocation("Local")
	customInfos := []model.CustomInfo{}
	for _, custom := range customs.List {
		lastUpdateDate, _ := time.ParseInLocation(baseFormat, custom.LastUpdateDate, loc)
		customInfo := model.CustomInfo{
			Id:             custom.Id,
			Name:           custom.Name,
			K8sKind:        custom.K8sKind,
			Status:         c.getStatus(custom.CommandStatus),
			LastUpdateDate: c.getTime(now.Sub(lastUpdateDate).Seconds()),
		}
		customInfos = append(customInfos, customInfo)
	}
	model.PrintCustomInfos(customInfos, out)
}

func (c *C7NClient) CreateCustom(out io.Writer, projectId int, data *url.Values) {
	if projectId == 0 {
		return
	}

	req, err := c.newRequestWithFormData("POST", fmt.Sprintf("/devops/v1/projects/%d/customize_resource", projectId), nil, data)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Println("Successfully create custom resource")
}
