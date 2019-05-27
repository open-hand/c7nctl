package c7nclient

import (
	"errors"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) CreateCert(out io.Writer, projectId int, certPostInfo *model.CertificationPostInfo) {
	if projectId == 0 {
		return
	}

	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/projects/%d/certifications", projectId), nil, certPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("Successfully created Certification %s", certPostInfo.CertName)

}

func (c *C7NClient) GetCert(out io.Writer, projectId int, envId int, name string) (error error, result *model.Certification) {
	if projectId == 0 {
		return errors.New("you do not have the permission of the project!"), nil
	}
	if envId == 0 {
		return errors.New("you do not have the permission of the env!"), nil
	}
	paras := make(map[string]interface{})
	paras["env_id"] = envId
	paras["cert_name"] = name
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/projects/%d/certifications/query_by_name", projectId), paras, nil)
	if err != nil {
		fmt.Printf("request build err:%v", err)
		return err, nil
	}
	devopsCertification := model.Certification{}
	_, err = c.do(req, &devopsCertification)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, nil
	}
	return err, &devopsCertification
}
