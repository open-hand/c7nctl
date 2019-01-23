package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"io"
	"io/ioutil"
	"net/http"
)


func (c *C7NClient) GetEnvSyncStatus(envId int) ( bool,error) {
	req,err := c.newRequest("GET",fmt.Sprintf("devops/v1/projects/%d/envs/%d/status",c.config.ProjectId,envId),nil,nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var envSyncStatus = model.EnvSyncStatus{}
	_,err = c.do(req,&envSyncStatus)
	if err != nil {
		return false,err

	}
	if envSyncStatus.AgentSyncCommit == envSyncStatus.DevopsSyncCommit && envSyncStatus.AgentSyncCommit  == envSyncStatus.SagaSyncCommit {
		return true, nil
	} else {
		return false, nil
	}

}


func (c *C7NClient) ListAuthEnvs(out io.Writer,) {
	if c.config.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	paras["active"]="true"
	req,err := c.newRequest("GET",fmt.Sprintf("/devops/v1/projects/%d/envs",c.config.ProjectId,),paras,nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var authEnvs = []model.AuthEnv{}
	resp,err := c.do(req,&authEnvs)
	if err != nil {
		fmt.Printf("request err:%v",err)
		return

	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return
	}
	envInfos := []model.EnvInfo{}
	for _,devOpsEnv := range authEnvs {
		if devOpsEnv.Failed || !devOpsEnv.Permission || !devOpsEnv.Connect{
			continue
		}
		status,err := c.GetEnvSyncStatus(devOpsEnv.ID)
		if err != nil {
			continue

		}

		envInfo := model.EnvInfo{
			Name: devOpsEnv.Name,
			Code: devOpsEnv.Code,
			Id:   devOpsEnv.ID,
			SyncStatus: status,
		}
		envInfos = append(envInfos, envInfo)
	}
	model.PrintAuthEnvInfo(envInfos,out)

}