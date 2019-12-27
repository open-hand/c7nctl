package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) ListEnvsInstance(out io.Writer, envId int) {
	if c.currentContext.User.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	paras["env_id"] = envId
	paras["page"] = 1
	paras["size"] = 100000
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/app_service_instances/info/page_by_options", c.currentContext.User.ProjectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var instances = model.Instances{}
	_, err = c.do(req, &instances)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return

	}
	instanceInfos := []model.InstanceInfo{}
	for _, instance := range instances.List {
		instanceInfo := model.InstanceInfo{
			Id:             instance.ID,
			Code:           instance.Code,
			VersionName:    instance.VersionName,
			AppServiceName: instance.AppServiceName,
			Status:         instance.Status,
			Pod:            fmt.Sprintf("%d/%d", instance.PodRunningCount, instance.PodCount),
		}
		instanceInfos = append(instanceInfos, instanceInfo)
	}
	model.PrintEnvInstanceInfo(instanceInfos, out)

}

func (c *C7NClient) InstanceResources(out io.Writer, instancesId int) {
	if c.currentContext.User.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/app_instances/%d/value", c.currentContext.User.ProjectId, instancesId), nil, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var resp = model.InstanceValues{}
	_, err = c.do(req, &resp)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return

	}
	fmt.Printf("The values of the instance:\n")
	fmt.Printf(resp.Values)

}

func (c *C7NClient) CreateInstance(out io.Writer, projectId int, instancePostInfo *model.InstancePostInfo) {
	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/app_service_instances", projectId), nil, instancePostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	applicationInstance := model.ApplicationInstanceDTO{}
	_, err = c.do(req, &applicationInstance)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}

	instanceInfos := []model.InstanceInfo{}

	instance := model.InstanceInfo{
		Code:   applicationInstance.Code,
		Status: applicationInstance.Status,
	}
	instanceInfos = append(instanceInfos, instance)

	model.PrintCreateInstanceInfo(instanceInfos, out)

}
