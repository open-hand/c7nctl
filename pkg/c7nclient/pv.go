package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) ListPv(out io.Writer, projectId int) {
	if projectId == 0 {
		return
	}

	body := make(map[string]interface{})
	body["param"] = ""
	body["searchParam"] = make(map[string]string)
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/pvs/page_by_options", projectId), nil, body)
	if err != nil {
		fmt.Printf("build request error")

	}
	var pvs = model.Pvs{}
	_, err = c.do(req, &pvs)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	pvInfos := []model.PvInfo{}
	for _, pv := range pvs.List {
		pvInfo := model.PvInfo{
			Id:              pv.Id,
			Name:            pv.Name,
			PvcName:         pv.PvcName,
			ClusterName:     pv.ClusterName,
			AccessModes:     pv.AccessModes,
			RequestResource: pv.RequestResource,
			Status:          pv.Status,
			Type:            pv.Type,
		}
		pvInfos = append(pvInfos, pvInfo)
	}
	model.PrintPvInfos(pvInfos, out)
}

func (c *C7NClient) CreatePv(out io.Writer, projectId int, pvPostInfo *model.PvPostInfo) {

	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/pvs", projectId), nil, pvPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("Successfully created Secret %s", pvPostInfo.Name)

}
