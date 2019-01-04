package authorize

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/utils"
	"github.com/vinkdong/gox/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	loginPath = "/oauth/oauth/token"
)

//curl -XPOST https://api.choerodon.com.cn/oauth/oauth/token\?grant_type\=password\&password\=UXdlcjEyMzQ\=\&username\=8377

func (a *Authorization) Login() error {
	var (
		username string
		password string
		err      error
	)
	userNameInput := utils.Input{
		Regex:    ".+",
		Password: false,
		Tip:      "请输入用户名: ",
	}
	passwordInput := utils.Input{
		Regex:    ".+",
		Password: true,
		Tip:      "请输入密码: ",
	}

input:
	if username, err = utils.AcceptUserInput(userNameInput); err != nil {
		return err
	}
	if password, err = utils.AcceptUserInput(passwordInput); err != nil {
		return err
	}

	params := url.Values{}
	params.Set("grant_type", "password")
	params.Set("password", base64Encoding(password))
	params.Set("client_id", "client")
	params.Set("client_secret", "secret")
	params.Set("username", username)
	reqUrl := fmt.Sprintf("%s%s?%s", a.ServerUrl, loginPath, params.Encode())
	log.Debugf("request to %s", reqUrl)
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		fmt.Println("登录失败，请重试")
		goto input
	}

	data, err := ioutil.ReadAll(resp.Body)
	loginResp := &LoginResp{}
	json.Unmarshal(data, loginResp)
	a.Token = loginResp.AccessToken
	a.Username = username
	return a.Write()
}

func base64Encoding(string string) string {
	return base64.StdEncoding.EncodeToString([]byte(string))
}
