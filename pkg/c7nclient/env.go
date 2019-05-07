package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

func (c *C7NClient) GetEnvSyncStatus(envId int) (bool, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/envs/%d/status", c.config.ProjectId, envId), nil, nil)
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

func (c *C7NClient) ListAuthEnvs(out io.Writer, projectId int) {
	if projectId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["active"] = "true"
	req, err := c.newRequest("GET", fmt.Sprintf("/devops/v1/projects/%d/envs", projectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var authEnvs = []model.AuthEnv{}
	resp, err := c.do(req, &authEnvs)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return

	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return
	}
	envInfos := []model.EnvInfo{}
	for _, devOpsEnv := range authEnvs {
		if devOpsEnv.Failed || !devOpsEnv.Permission || !devOpsEnv.Connect {
			continue
		}
		status, err := c.GetEnvSyncStatus(devOpsEnv.ID)
		if err != nil {
			continue

		}

		envInfo := model.EnvInfo{
			Name:       devOpsEnv.Name,
			Code:       devOpsEnv.Code,
			Id:         devOpsEnv.ID,
			SyncStatus: status,
		}
		envInfos = append(envInfos, envInfo)
	}
	model.PrintAuthEnvInfo(envInfos, out)

}

func (c *C7NClient) ListEnvs(out io.Writer, projectId int) {
	if projectId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["active"] = "true"
	req, err := c.newRequest("GET", fmt.Sprintf("/devops/v1/projects/%d/envs/groups", projectId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var devOpsEnvs = []model.DevOpsEnvs{}
	_, err = c.do(req, &devOpsEnvs)

	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}

	envInfos := []model.EnvInfo{}
	for _, devOpsEnv := range devOpsEnvs[0].DevopsEnviromentRepDTOs {
		var status string
		if devOpsEnv.Failed {
			status = "Failed"
		} else if devOpsEnv.Connect {
			status = "Connected"
		} else {
			status = "Disconnected"
		}
		envInfo := model.EnvInfo{
			Name:    devOpsEnv.Name,
			Status:  status,
			Code:    devOpsEnv.Code,
			Cluster: devOpsEnv.ClusterName,
			Group:   devOpsEnvs[0].DevopsEnvGroupName,
		}
		envInfos = append(envInfos, envInfo)
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
	fmt.Printf("the env is create success!")
}
