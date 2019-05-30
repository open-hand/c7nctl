package c7nclient

import (
	"encoding/base64"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/choerodon/c7nctl/pkg/utils"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const client = "c7nclient"
const secret = "secret"
const grantType = "password"

func (c *C7NClient) Login(out io.Writer) {

	if c.config.User.Token != "" {
		fmt.Println("you have login, you can use logout when you want to login of other user or other env")
		return
	}

	var (
		username string
		password string
		err      error
	)

	username, err = utils.AcceptUserInput(utils.Input{
		Password: false,
		Tip:      "请输入用户名: ",
		Regex:    ".+",
	})
	password, err = utils.AcceptUserInput(utils.Input{
		Password: true,
		Tip:      "请输入密码: ",
		Regex:    ".+",
	})

	home, err := homedir.Dir()
	configDir := fmt.Sprintf("%s%c.c7n%c%s", home, os.PathSeparator, os.PathSeparator, "config.yaml")

	c.BaseURL = c.config.Server

	strbytes := []byte(password)
	password = base64.StdEncoding.EncodeToString(strbytes)

	paras := make(map[string]interface{})
	paras["client_id"] = client
	paras["client_secret"] = secret
	paras["grant_type"] = grantType
	paras["password"] = password
	paras["username"] = strings.TrimSpace(username)

	req, err := c.newRequest("POST", "oauth/oauth/token", paras, nil)
	if err != nil {
		fmt.Println("build request error")
		os.Exit(1)
	}
	var token = model.Token{}
	_, err = c.do(req, &token)
	if err != nil {
		fmt.Println("username or password is error!")
		os.Exit(1)
	}
	c.config.User.Token = token.AccessToken
	c.config.User.UserName = username
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
	c.config.User.OrganizationId = organization.ID
	c.config.User.OrganizationCode = organization.Code
	projects := viper.Get("pros")
	for _, pro := range projects.([]model.Project) {
		if pro.OrganizationID == organization.ID {
			c.config.User.ProjectId = pro.ID
			c.config.User.ProjectCode = pro.Code
			break
		}
	}

	var allConfig C7NConfig
	viper.Unmarshal(&allConfig)

	for i, context := range allConfig.Contexts {
		if context.Name == allConfig.CurrentContext {
			allConfig.Contexts[i] = *c.config
		}
	}

	bytes, err := yaml.Marshal(allConfig)

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
	fmt.Println("Login Succeeded!")
}

func (c *C7NClient) Logout(out io.Writer) {

	var allConfig C7NConfig
	viper.Unmarshal(&allConfig)

	for i, context := range allConfig.Contexts {
		if context.Name == allConfig.CurrentContext {
			allConfig.Contexts[i].User = C7NUser{}
		}
	}

	bytes, _ := yaml.Marshal(allConfig)
	if ioutil.WriteFile(viper.ConfigFileUsed(), bytes, 0644) != nil {
		fmt.Println("modify config file failed")
	}
	fmt.Println("Login Out")

}

func (c *C7NClient) SwitchContext(out io.Writer, name string) {

	var allConfig C7NConfig
	viper.Unmarshal(&allConfig)

	var index int
	for _, context := range allConfig.Contexts {
		if context.Name == name {
			allConfig.CurrentContext = name
		} else {
			index++
		}
	}
	if index == len(allConfig.Contexts) {
		fmt.Println("The context is not exist in the config.yaml")
		return
	}

	bytes, _ := yaml.Marshal(allConfig)
	if ioutil.WriteFile(viper.ConfigFileUsed(), bytes, 0644) != nil {
		fmt.Println("modify config file failed")
		return
	}

}
