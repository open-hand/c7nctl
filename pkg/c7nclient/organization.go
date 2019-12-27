package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
)

func (c *C7NClient) ListOrganization(out io.Writer, userId int) {
	req, err := c.newRequest("GET", fmt.Sprintf("base/v1/users/%d/organizations", userId), nil, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var orgs = []model.Organization{}
	_, err = c.do(req, &orgs)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	orgInfos := []model.OrganizationInfo{}
	for _, org := range orgs {
		orgInfo := model.OrganizationInfo{
			Name: org.Name,
			Code: org.Code,
		}
		orgInfos = append(orgInfos, orgInfo)
	}
	model.PrintOrgInfo(orgInfos, out)
}

func (c *C7NClient) SetOrganization(out io.Writer, userId int) (error error) {
	req, err := c.newRequest("GET", fmt.Sprintf("base/v1/users/%d/organizations", userId), nil, nil)
	if err != nil {
		fmt.Printf("build request error")
		return err
	}
	var orgs = []model.Organization{}
	_, err = c.do(req, &orgs)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err
	}
	viper.Set("orgs", orgs)
	return nil
}

func (c *C7NClient) UseOrganization(out io.Writer, orgCode string) {
	orgs := viper.Get("orgs")

	var index int
	for _, org := range orgs.([]model.Organization) {
		if org.Code == orgCode {

			for index, context := range c.platformConfig.Contexts {
				if context.Name == c.platformConfig.CurrentContext {
					c.currentContext.User.OrganizationId = org.ID
					c.currentContext.User.OrganizationCode = orgCode
					c.platformConfig.Contexts[index].User.OrganizationId=org.ID
					c.platformConfig.Contexts[index].User.OrganizationCode=orgCode
				}
			}

			bytes, _ := yaml.Marshal(c.platformConfig)
			if ioutil.WriteFile(viper.ConfigFileUsed(), bytes, 0644) != nil {
				fmt.Println("modify config file failed")
			}
			break
		} else {
			index++
			if index == len(orgs.([]model.Organization)) {
				fmt.Printf("you do not have the permission of this organization:%v", orgCode)
			}
		}
	}
}

func (c *C7NClient) GetOrganization(out io.Writer, userId int, orgCode string) (error error, organizationId int) {
	if orgCode == "" {
		return nil, c.currentContext.User.OrganizationId
	} else {
		orgs := viper.Get("orgs")
		var index int
		for _, org := range orgs.([]model.Organization) {
			if org.Code == orgCode {
				return nil, org.ID
			} else {
				index++
				if index == len(orgs.([]model.Organization)) {
					fmt.Printf("you do not have the permission of the organization:%v", orgCode)
					return errors.New("you do not have the permission of the organization"), 0
				}
			}
		}
		return
	}
}
