package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"io"
	"strings"
)

// /devops/v1/projects/42/ingress/54/listByEnv
func (c *C7NClient) ListIngress(out io.Writer, envId int) {

	if c.config.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = 0
	paras["size"] = 99
	paras["sort"] = "id,desc"
	body := make(map[string]interface{})
	body["param"] = ""
	body["searchParam"] = make(map[string]string)
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/ingress/%d/listByEnv", c.config.ProjectId, envId), paras, body)
	if err != nil {
		fmt.Printf("build request error")
	}
	var resp = model.DevopsIngressPage{}
	_, err = c.do(req, &resp)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return

	}
	ingressInfoList := []model.DevOpsIngressInfo{}
	for _, ingress := range resp.Content {
		ingressInfo := model.DevOpsIngressInfo{
			Id:     ingress.ID,
			Name:   ingress.Name,
			Host:   ingress.Domain,
			Status: ingress.Status,
		}
		var paths = []string{}
		if len(ingress.PathList) != 0 {
			for _, ingressPath := range ingress.PathList {
				paths = append(paths, fmt.Sprintf("%s -> %s", ingressPath.Path, ingressPath.ServiceName))
			}
			ingressInfo.Paths = strings.Join(paths, ",")
		}

		ingressInfoList = append(ingressInfoList, ingressInfo)
	}
	model.PrintIngressInfo(ingressInfoList, out)

}

func (c *C7NClient) CreateIngress(out io.Writer, projectId int, ingressPostInfo *model.IngressPostInfo) {
	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/ingress", projectId, ), nil, ingressPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("create ingress %s success!", ingressPostInfo.Name)

}
