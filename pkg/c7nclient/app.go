package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) ListApps(out io.Writer, projectId int) {
	if projectId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = "0"
	paras["size"] = "10"
	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/apps/list_by_options", projectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var apps = model.Apps{}
	_, err = c.do(req, &apps)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	appInfos := []model.AppInfo{}
	for _, app := range apps.Content {
		var status string
		var types string
		if app.Synchro {
			status = "可用"
		} else {
			status = "不可用"
		}
		if app.Type == "normal" {
			types = "普通应用"
		} else {
			types = "测试应用"
		}
		appInfo := model.AppInfo{
			Name:    app.Name,
			Code:    app.Code,
			RepoURL: app.RepoURL,
			Status:  status,
			Type:    types,
		}
		appInfos = append(appInfos, appInfo)
	}
	model.PrintAppInfo(appInfos, out)
}

func (c *C7NClient) GetApp(appCode string, projectId int) (error error, result *model.App) {

	if projectId == 0 {
		return nil, nil
	}
	paras := make(map[string]interface{})
	paras["code"] = appCode
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/apps/query_by_code", projectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
		return err, nil
	}
	var app = model.App{}
	_, err = c.do(req, &app)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, nil
	}
	return nil, &app
}

func (c *C7NClient) CreateApp(out io.Writer, projectId int, appPostInfo *model.AppPostInfo) {
	if projectId == 0 {
		return
	}
	if appPostInfo.ApplicationTemplateId == 0 {
		fmt.Printf("the app template you hava choose not exist!")
		return
	}
	if appPostInfo.Type != "normal" && appPostInfo.Type != "test" {
		fmt.Printf("the app type value should be normal or test!")
		return
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/apps", projectId), nil, appPostInfo)
	if err != nil {
		fmt.Printf("build request error")

	}
	var app = model.App{}
	_, err = c.do(req, &app)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	var status string
	var types string
	appInfos := []model.AppInfo{}
	if app.Synchro {
		status = "可用"
	} else {
		status = "不可用"
	}
	if app.Type == "normal" {
		types = "普通应用"
	} else {
		types = "测试应用"
	}
	appInfo := model.AppInfo{
		Name:    app.Name,
		Code:    app.Code,
		RepoURL: app.RepoURL,
		Status:  status,
		Type:    types,
	}
	appInfos = append(appInfos, appInfo)
	model.PrintAppInfo(appInfos, out)

}
