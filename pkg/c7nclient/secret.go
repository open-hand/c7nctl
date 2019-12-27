package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
	"strings"
	"time"
)

func (c *C7NClient) ListSecret(out io.Writer, projectId int, envId int) {
	if projectId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = 1
	paras["size"] = 10000
	paras["env_id"] = envId

	body := make(map[string]interface{})
	body["param"] = ""
	body["searchParam"] = make(map[string]string)
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/secret/page_by_options", projectId), paras, body)
	if err != nil {
		fmt.Printf("build request error")

	}
	var secrets = model.Secrets{}
	_, err = c.do(req, &secrets)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	now := time.Now()
	loc, _ := time.LoadLocation("Local")
	secretInfos := []model.SecretInfo{}
	for _, secret := range secrets.List {
		lastUpdateDate, _ := time.ParseInLocation(baseFormat, secret.LastUpdateDate, loc)
		secretInfo := model.SecretInfo{
			Id:             secret.Id,
			Name:           secret.Name,
			Key:            strings.Join(secret.Key, ","),
			Status:         c.getStatus(secret.CommandStatus),
			LastUpdateDate: c.getTime(now.Sub(lastUpdateDate).Seconds()),
		}
		secretInfos = append(secretInfos, secretInfo)
	}
	model.PrintSecretInfos(secretInfos, out)
}

func (c *C7NClient) CreateSecret(out io.Writer, projectId int, secretPostInfo *model.SecretPostInfo) {

	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/secret", projectId), nil, secretPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("Successfully created Secret %s", secretPostInfo.Name)

}
