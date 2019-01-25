package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"io"
	"strconv"
)

func (c *C7NClient) ListAppVersions(out io.Writer, appCode *string) {
	if c.config.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	if *appCode != "" {
		app, _ := c.getApp(appCode)
		paras["appId"] = strconv.Itoa(app.ID)
	}
	paras["page"] = "0"
	paras["size"] = "10"
	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/app_versions/list_by_options", c.config.ProjectId, ), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var appVersions = model.AppVersions{}
	_, err = c.do(req, &appVersions)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	appVersionInfos := []model.AppVersionInfo{}
	for _, appVersion := range appVersions.Content {
		appVersionInfo := model.AppVersionInfo{
			AppName:      appVersion.AppName,
			AppCode:      appVersion.AppCode,
			Version:      appVersion.Version,
			CreationDate: appVersion.CreationDate,
		}
		appVersionInfos = append(appVersionInfos, appVersionInfo)
	}
	model.PrintAppVersionInfo(appVersionInfos, out)
}
