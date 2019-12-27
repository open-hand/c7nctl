package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
	"strings"
	"time"
)

func (c *C7NClient) ListConfigMap(out io.Writer, projectId int, envId int) {
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
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/projects/%d/config_maps/page_by_options", projectId), paras, body)
	if err != nil {
		fmt.Printf("build request error")

	}
	var configMaps = model.ConfigMaps{}
	_, err = c.do(req, &configMaps)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	now := time.Now()
	loc, _ := time.LoadLocation("Local")
	configMapInfos := []model.ConfigMapInfo{}
	for _, configMap := range configMaps.List {
		lastUpdateDate, _ := time.ParseInLocation(baseFormat, configMap.LastUpdateDate,loc)
		configMapInfo := model.ConfigMapInfo{
			Id:             configMap.Id,
			Name:           configMap.Name,
			Key:            strings.Join(configMap.Key, ","),
			Status:         c.getStatus(configMap.CommandStatus),
			LastUpdateDate: c.getTime(now.Sub(lastUpdateDate).Seconds()),
		}
		configMapInfos = append(configMapInfos, configMapInfo)
	}
	model.PrintConfigMapInfos(configMapInfos, out)
}

func (c *C7NClient) CreateConfigMap(out io.Writer, projectId int, configMapPostInfo *model.ConfigMapPostInfo) {
	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/config_maps", projectId), nil, configMapPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("Successfully created ConfigMap %s", configMapPostInfo.Name)

}
