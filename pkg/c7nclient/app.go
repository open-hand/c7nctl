package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) ListApps(out io.Writer, ) {
	if c.config.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = "0"
	paras["size"] = "10"
	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/apps/list_by_options", c.config.ProjectId, ), paras, nil)
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

func (c *C7NClient) getApp(appCode *string) (*model.App, error) {
	if c.config.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return nil, nil
	}
	paras := make(map[string]interface{})
	paras["code"] = *appCode
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/apps/query_by_code", c.config.ProjectId, ), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
		return nil, err
	}
	var app = model.App{}
	_, err = c.do(req, &app)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return nil, err
	}
	return &app, nil
}
