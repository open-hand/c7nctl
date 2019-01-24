package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) ListAppTemplates(out io.Writer, ) {
	if c.config.OrganizationId == -1 {
		fmt.Printf("Set organization Id")
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = "0"
	paras["size"] = "10"
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/organizations/%d/app_templates/list_by_options", c.config.OrganizationId, ), paras, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var appTemplates = model.AppTemplates{}
	_, err = c.do(req, &appTemplates)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	appTemplateInfos := []model.AppTemplateInfo{}
	for _, appTemplate := range appTemplates.Content {
		var available string
		if appTemplate.Synchro {
			available = "可用"
		} else {
			available = "不可用"
		}
		appTemplateInfo := model.AppTemplateInfo{
			Name:      appTemplate.Name,
			Code:      appTemplate.Code,
			RepoUrl:   appTemplate.RepoURL,
			Available: available,
		}
		appTemplateInfos = append(appTemplateInfos, appTemplateInfo)
	}
	model.PrintAppTemplateInfo(appTemplateInfos, out)

}
