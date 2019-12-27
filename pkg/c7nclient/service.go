package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/pkg/errors"
	"io"
	"strings"
)

func (c *C7NClient) ListService(out io.Writer, envId int) {

	if c.currentContext.User.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = 0
	paras["size"] = 10000
	paras["env_id"] = envId
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/service/page_by_options", c.currentContext.User.ProjectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var resp = model.Services{}
	_, err = c.do(req, &resp)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return

	}
	devOpsServiceInfos := []model.DevOpsServiceInfo{}
	for _, service := range resp.List {
		serviceInfo := model.DevOpsServiceInfo{
			Id:     service.ID,
			Name:   service.Name,
			Type:   service.Type,
			Status: service.Status,
		}
		var targetContent string
		if len(service.Target.Instances) != 0 {
			targetInstances := []string{}
			for _, ins := range service.Target.Instances {
				targetInstances = append(targetInstances, ins.Code)
			}
			targetContent = strings.Join(targetInstances, ",")
			serviceInfo.TargetType = "instance"
		} else if service.Target.Labels != nil {
			targetLabels := []string{}
			for k, v := range service.Target.Labels {
				targetLabels = append(targetLabels, fmt.Sprintf("%s=%s", k, v))
			}
			targetContent = strings.Join(targetLabels, ",")
			serviceInfo.TargetType = "label"
		} else {
			for k, _ := range service.Target.EndPoints {
				targetContent = k
			}
			serviceInfo.TargetType = "endpoint"
		}
		serviceInfo.Target = targetContent
		devOpsServiceInfos = append(devOpsServiceInfos, serviceInfo)
	}
	model.PrintServiceInfo(devOpsServiceInfos, out)
}

func (c *C7NClient) GetService(out io.Writer, projectId int, envId int, name string) (error error, result *model.DevOpsService) {
	if projectId == 0 {
		return errors.New("you do not have the permission of the project!"), nil
	}
	if envId == 0 {
		return errors.New("you do not have the permission of the env!"), nil
	}
	paras := make(map[string]interface{})
	paras["env_id"] = envId
	paras["name"] = name
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/service/query_by_name", projectId), paras, nil)
	if err != nil {
		fmt.Printf("request build err:%v", err)
		return err, nil
	}
	devopsService := model.DevOpsService{}
	_, err = c.do(req, &devopsService)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, nil
	}
	return err, &devopsService
}

func (c *C7NClient) CreateService(out io.Writer, projectId int, servicePostInfo *model.ServicePostInfo) {
	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/service", projectId), nil, servicePostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("Successfully created Service %s", servicePostInfo.Name)

}
