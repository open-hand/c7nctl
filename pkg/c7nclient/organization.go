package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"github.com/pkg/errors"
	"io"
	"strconv"
)

func (c *C7NClient) ListOrganization(out io.Writer, userId int) {
	req, err := c.newRequest("GET", fmt.Sprintf("iam/v1/users/%d/organizations", userId, ), nil, nil)
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
	OrgMap.Map[strconv.Itoa(userId)] = orgs
	for _, org := range orgs {
		orgInfo := model.OrganizationInfo{
			Name: org.Name,
			Code: org.Code,
		}
		orgInfos = append(orgInfos, orgInfo)
	}
	model.PrintOrgInfo(orgInfos, out)
}

func (c *C7NClient) SetOrganization(out io.Writer, userId int) (error error){
	req, err := c.newRequest("GET", fmt.Sprintf("iam/v1/users/%d/organizations", userId, ), nil, nil)
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
	OrgMap.Map[strconv.Itoa(userId)] = orgs
	return nil
}

func (c *C7NClient) UseOrganization(out io.Writer, userId int, orgCode string) {
	values := OrgMap.Map[(strconv.Itoa(userId))]
	var index int
	for _, org := range values.([]model.Organization) {
		if org.Code == orgCode {
			c.config.Organization = orgCode
			c.config.OrganizationId = org.ID
			c.config.OrganizationCode = orgCode
			break;
		} else {
			index ++
			if index == len(values.([]model.Organization)) {
				fmt.Printf("you do not have the permission of this organization:%v", orgCode)
			}
		}
	}
}

func (c *C7NClient) GetOrganization(out io.Writer, userId int, orgCode string) (error error,organizationId int) {
	if orgCode == "" {
		return nil, c.config.OrganizationId
	} else {
		values := OrgMap.Map[(strconv.Itoa(userId))]
		var index int
		for _, org := range values.([]model.Organization) {
			if org.Code == orgCode {
				return nil, org.ID
			} else {
				index ++
				if index == len(values.([]model.Organization)) {
					fmt.Printf("you do not have the permission of the organization:%v", orgCode)
					return errors.New("you do not have the permission of the organization"),0
				}
			}
		}
		return
	}
}
