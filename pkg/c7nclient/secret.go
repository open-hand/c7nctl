package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) CreateSecret(out io.Writer, projectId int, secretPostInfo *model.SecretPostInfo) {

	if projectId == 0 {
		return
	}

	req, err := c.newRequest("PUT", fmt.Sprintf("devops/v1/projects/%d/secret", projectId), nil, secretPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var result string
	_, err = c.doHandleString(req, &result)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Printf("create Secret %s success!", secretPostInfo.Name)

}
