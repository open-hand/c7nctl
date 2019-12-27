package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) ListPvc(out io.Writer, projectId int, envId int) {
	if projectId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["env_id"] = envId

	body := make(map[string]interface{})
	body["param"] = ""
	body["searchParam"] = make(map[string]string)
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/pvcs/page_by_options", projectId), paras, body)
	if err != nil {
		fmt.Printf("build request error")
		return
	}
	var pvcs = model.Pvcs{}
	_, err = c.do(req, &pvcs)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	pvcInfos := []model.PvcInfo{}
	for _, pvc := range pvcs.List {
		pvcInfo := model.PvcInfo{
			Id:              pvc.Id,
			Name:            pvc.Name,
			PvName:          pvc.PvName,
			AccessModes:     pvc.AccessModes,
			RequestResource: pvc.RequestResource,
			Status:          pvc.Status,
		}
		pvcInfos = append(pvcInfos, pvcInfo)
	}
	model.PrintPvcInfos(pvcInfos, out)
}

func (c *C7NClient) CreatePvc(out io.Writer, projectId int, pvcPostInfo *model.PvcPostInfo) {

	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/pvcs", projectId), nil, pvcPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("Successfully created Secret %s", pvcPostInfo.Name)

}
