package c7nclient

import (
	"encoding/base64"
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
)

const client = "c7nclient"
const secret = "secret"
const grantType = "password"

func (c *C7NClient) Login(out io.Writer, password string, username string, server string) {

	if c.config.Token != "" {
		fmt.Println("you have login, you can use logout when you want to login of other user or other env")
		return
	}

	home, err := homedir.Dir()
	configDir := home + string(os.PathSeparator) + ".c7n.yaml"

	c.BaseURL = server

	strbytes := []byte(password)
	password = base64.StdEncoding.EncodeToString(strbytes)

	paras := make(map[string]interface{})
	paras["client_id"] = client
	paras["client_secret"] = secret
	paras["grant_type"] = grantType
	paras["password"] = password
	paras["username"] = username

	req, err := c.newRequest("POST", "oauth/oauth/token", paras, nil)

	if err != nil {
		fmt.Println("build request error")
		os.Exit(1)
	}
	var token = model.Token{}
	_, err = c.do(req, &token)
	c.config.Server = server
	c.config.Token = token.AccessToken
	err, user := c.QuerySelf(out)
	if err != nil {
		fmt.Println("query self error")
		os.Exit(1)
	}
	err = c.SetOrganization(out, user.ID)
	if err != nil {
		fmt.Println("set organization error")
		os.Exit(1)
	}
	err = c.SetProject(out, user.ID)
	if err != nil {
		fmt.Println("set project error")
		os.Exit(1)
	}
	organizations := viper.Get("orgs")
	organization := organizations.([]model.Organization)[0]
	c.config.OrganizationId = organization.ID
	c.config.OrganizationCode = organization.Code
	projects := viper.Get("pros")
	for _, pro := range projects.([]model.Project) {
		if pro.OrganizationID == organization.ID {
			c.config.ProjectId = pro.ID
			c.config.ProjectCode = pro.Code
			break
		}
	}
	bytes, err := yaml.Marshal(c.config)

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		_, err = os.Create(configDir)
		if err != nil {
			fmt.Println(err)
		}
	}
	if ioutil.WriteFile(configDir, bytes, 0644) != nil {
		fmt.Println("modify config file failed")
	}
	fmt.Println("Login Success")
}

func (c *C7NClient) Logout(out io.Writer) {

	emptyContext := C7NPlatformContext{}
	bytes, _ := yaml.Marshal(emptyContext)
	if ioutil.WriteFile(viper.ConfigFileUsed(), bytes, 0644) != nil {
		fmt.Println("modify config file failed")
	}
	fmt.Println("Login Out")

}
