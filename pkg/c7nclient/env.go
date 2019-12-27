package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/pkg/errors"
	"io"
)

func (c *C7NClient) GetEnvSyncStatus(envId int) (bool, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/envs/%d/status", c.currentContext.User.ProjectId, envId), nil, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var envSyncStatus = model.EnvSyncStatus{}
	_, err = c.do(req, &envSyncStatus)
	if err != nil {
		return false, err

	}
	if envSyncStatus.AgentSyncCommit == envSyncStatus.DevopsSyncCommit && envSyncStatus.AgentSyncCommit == envSyncStatus.SagaSyncCommit {
		return true, nil
	} else {
		return false, nil
	}

}

func (c *C7NClient) ListEnvs(out io.Writer, projectId int) {
	if projectId == 0 {
		return
	}
	req, err := c.newRequest("GET", fmt.Sprintf("/devops/v1/projects/%d/envs/env_tree_menu", projectId), nil, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var devOpsEnvsInGroups []model.DevOpsEnvs
	_, err = c.do(req, &devOpsEnvsInGroups)

	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}

	var envInfos []model.EnvInfo
	for _, devopsEnvInGroup := range devOpsEnvsInGroups {
		for _, devOpsEnv := range devopsEnvInGroup.DevopsEnvironmentRepDTOs {
			var status string
			if devOpsEnv.Failed {
				status = "Failed"
			} else if devOpsEnv.Connect {
				status = "Connected"
			} else {
				status = "Disconnected"
			}
			envInfo := model.EnvInfo{
				Id:      devOpsEnv.ID,
				Name:    devOpsEnv.Name,
				Status:  status,
				Code:    devOpsEnv.Code,
				Cluster: devOpsEnv.ClusterName,
				Group:   devopsEnvInGroup.DevopsEnvGroupName,
			}
			envInfos = append(envInfos, envInfo)
		}
	}
	model.PrintEnvInfo(envInfos, out)

}

func (c *C7NClient) GetEnv(out io.Writer, projectId int, code string) (error error, result *model.DevOpsEnv) {
	if projectId == 0 {
		return errors.New("the project you choose is not found!"), nil
	}
	if code == "" {
		return errors.New("the env code is empty!"), nil
	}
	paras := make(map[string]interface{})
	paras["code"] = code
	req, err := c.newRequest("GET", fmt.Sprintf("/devops/v1/projects/%d/envs/query_by_code", projectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var devopsEnv = model.DevOpsEnv{}
	_, err = c.do(req, &devopsEnv)

	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, nil
	}
	return nil, &devopsEnv
}

func (c *C7NClient) CreateEnv(out io.Writer, projectId int, envPostInfo *model.EnvPostInfo) {
	if projectId == 0 {
		return
	}
	if envPostInfo.ClusterId == 0 {
		fmt.Printf("the cluster you choose is not exist!")
		return
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/envs", projectId), nil, envPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}

	_, err = c.do(req, nil)

	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Sprintf("Successfully created env %s", envPostInfo.Name)
}
