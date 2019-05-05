package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"io"
)

func (c *C7NClient) QuerySelf(out io.Writer) (error error, info model.UserInfo) {

	req, err := c.newRequest("GET", "/iam/v1/users/self", nil, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var userInfo = model.UserInfo{}
	_, err = c.do(req, &userInfo)
	if err != nil {
		fmt.Println(err)
		return err, userInfo
	}
	return nil, userInfo
}

func (c *C7NClient) QueryGitlabUserId(out io.Writer, userId int, projectId int) (error error, gitlabUserId int) {

	req, err := c.newRequest("GET", fmt.Sprintf("/v1/projects/%d/users/%d", projectId, userId, ), nil, nil)

	if err != nil {
		fmt.Printf("build request error")
	}
	var userAttrInfo = model.UserAttrInfo{}
	_, err = c.do(req, &userAttrInfo)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, 0
	}
	return nil, userAttrInfo.GitlabUserId

}
