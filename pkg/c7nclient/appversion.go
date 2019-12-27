package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/pkg/errors"
	"io"
	"strconv"
)

func (c *C7NClient) ListAppVersions(out io.Writer, appCode *string, projectId int) {
	if projectId == 0 {
		return
	}
	paras := make(map[string]interface{})
	if *appCode != "" {
		_, app := c.GetApp(*appCode, projectId)
		if app == nil {
			fmt.Println("the project do not hava the application!")
			return
		}
		paras["app_service_id"] = strconv.Itoa(app.ID)
	}
	paras["do_page"] = "false"
	paras["deploy_only"] = "true"
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/app_service_versions/page_by_options", projectId), paras, nil)
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
	for _, appVersion := range appVersions.List {
		appVersionInfo := model.AppVersionInfo{
			Id:           appVersion.ID,
			Version:      appVersion.Version,
			CreationDate: appVersion.CreationDate,
		}
		appVersionInfos = append(appVersionInfos, appVersionInfo)
	}
	model.PrintAppVersionInfo(appVersionInfos, out)
}

func (c *C7NClient) GetAppVersion(out io.Writer, projectId int, version string, appId int) (error error, result *model.AppVersion) {
	if projectId == 0 {
		return errors.New("the project you choose is not found"), nil
	}
	if version == "" {
		return errors.New("the app version you choose is not found"), nil
	}
	if appId == 0 {
		return errors.New("the app you choose is not found"), nil
	}
	paras := make(map[string]interface{})
	paras["version"] = version
	paras["appId"] = strconv.Itoa(appId)
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/app_versions/query_by_version", c.currentContext.User.ProjectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
		return err, nil
	}
	var appVersion = model.AppVersion{}
	_, err = c.do(req, &appVersion)
	return nil, &appVersion
}
