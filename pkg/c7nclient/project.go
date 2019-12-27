package c7nclient

import (
	"errors"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
)

func (c *C7NClient) ListProject(out io.Writer, userId int) {
	req, err := c.newRequest("GET", fmt.Sprintf("base/v1/users/%d/projects", userId), nil, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var pros = []model.Project{}
	_, err = c.do(req, &pros)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	proInfos := []model.ProjectInfo{}
	for _, pro := range pros {
		if pro.OrganizationID == c.currentContext.User.OrganizationId {
			proInfo := model.ProjectInfo{
				Name: pro.Name,
				Code: pro.Code,
			}
			proInfos = append(proInfos, proInfo)
		}
	}
	model.PrintProInfo(proInfos, out)
}

func (c *C7NClient) SetProject(out io.Writer, userId int) (error error) {
	req, err := c.newRequest("GET", fmt.Sprintf("base/v1/users/%d/projects", userId), nil, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var pros = []model.Project{}
	_, err = c.do(req, &pros)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err
	}
	viper.Set("pros", pros)
	return nil
}

func (c *C7NClient) UseProject(out io.Writer, proCode string) {
	pros := viper.Get("pros")
	var index int
	for _, pro := range pros.([]model.Project) {
		if pro.Code == proCode {
			for index, context := range c.platformConfig.Contexts {
				if pro.OrganizationID == context.User.OrganizationId && context.Name == c.platformConfig.CurrentContext {
					c.currentContext.User.ProjectId = pro.ID
					c.currentContext.User.ProjectCode = proCode
					c.platformConfig.Contexts[index].User.ProjectId = pro.ID
					c.platformConfig.Contexts[index].User.ProjectCode = proCode
				}
			}

			bytes, _ := yaml.Marshal(c.platformConfig)
			if ioutil.WriteFile(viper.ConfigFileUsed(), bytes, 0644) != nil {
				fmt.Println("modify config file failed")
			}
			break
		} else {
			index++
			if index == len(pros.([]model.Project)) {
				fmt.Printf("you do not have the permission of this project:%v", proCode)
			}
		}
	}
}

func (c *C7NClient) GetProject(out io.Writer, userId int, proCode string) (error error, project model.Project) {
	if proCode == "" {
		pro := model.Project{}
		pro.ID = c.currentContext.User.ProjectId
		pro.OrganizationID = c.currentContext.User.OrganizationId
		return nil, pro
	} else {
		pros := viper.Get("pros")
		var index int
		for _, pro := range pros.([]model.Project) {
			if pro.Code == proCode {
				return nil, pro
			} else {
				index++
				if index == len(pros.([]model.Project)) {
					fmt.Printf("you do not have the permission of this project:%v", proCode)
					return errors.New("you do not have the permission of this project"), model.Project{}
				}
			}
		}
		return
	}
}
