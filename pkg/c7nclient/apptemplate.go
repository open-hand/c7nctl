package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/pkg/errors"
	"io"
)

func (c *C7NClient) ListAppTemplates(out io.Writer, organizationId int) {
	if organizationId ==0 {
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = "0"
	paras["size"] = "20"
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/organizations/%d/app_templates/list_by_options", organizationId, ), paras, nil)
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

func (c *C7NClient) GetAppTemplate(out io.Writer, organizationId int, appTemplateCode string) (error error, result model.AppTemplate) {
	if organizationId ==0 {
		return errors.New("the organization is not found!"), model.AppTemplate{}
	}
	paras := make(map[string]interface{})
	paras["code"] = appTemplateCode
	req, err := c.newRequest("GET", fmt.Sprintf("/devops/v1/organizations/%d/app_templates/query_by_code", organizationId, ), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
		return errors.New("build request error"), model.AppTemplate{}
	}
	var appTemplate = model.AppTemplate{}
	_, err = c.do(req, &appTemplate)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, appTemplate
	}
	return nil, appTemplate
}

func (c *C7NClient) CreateAppTemplate(out io.Writer, organizationId int, appTemplatePostInfo *model.AppTemplatePostInfo) {
	if organizationId ==0 {
		return
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/organizations/%d/app_templates", organizationId, ), nil, appTemplatePostInfo)
	if err != nil {
		fmt.Printf("build request error")

	}
	var appTemplate = model.AppTemplate{}
	_, err = c.do(req, &appTemplate)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	appTemplateInfos := []model.AppTemplateInfo{}
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
	model.PrintAppTemplateInfo(appTemplateInfos, out)

}
