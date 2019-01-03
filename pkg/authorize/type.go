package authorize

import (
	"github.com/choerodon/c7nctl/pkg/utils"
	"fmt"
	"github.com/vinkdong/gox/log"
	"net/http"
)

type Authorization struct {
	Token       string
	ServerUrl   string
	Username    string
	Config      *utils.Config
	ClusterName string
	TokenType   string
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	CreateTime   int64  `json:"createTime"`
	SessionId    string `json:"sessionId"`
}

func (a *Authorization) Write() error {
	var (
		user    *utils.User
		cluster *utils.Cluster
	)
	user = a.Config.CurrentUser()
	userName := fmt.Sprintf("%s-%s", a.ClusterName, a.Username)
	if user == nil {
		user = &utils.User{}
		a.Config.Users = append(a.Config.Users, &utils.NamedUser{
			Name: userName,
			User: user,
		})
	}
	user.Token = a.Token
	user.Name = a.Username
	cluster = a.Config.CurrentCluster()
	if cluster == nil {
		cluster = &utils.Cluster{
		}
		a.Config.Clusters = append(a.Config.Clusters, &utils.NamedCluster{
			Name:    a.ClusterName,
			Cluster: cluster,
		})
	}
	cluster.Server = a.ServerUrl
	cluster.SelectedUser = userName
	a.Config.SelectedCluster = a.ClusterName

	return a.Config.Write()
}

func (a *Authorization) IsAuthorized() bool {
	auth := DefaultAuthorization()
	//todo: check from server
	if auth.Token == "" {
		return false
	}
	return true
}

func (a *Authorization) Request(req *http.Request) (*http.Response, error) {
	client := http.Client{}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", a.TokenType, a.Token))
	return client.Do(req)
}

func DefaultAuthorization(configs ...*utils.Config) *Authorization {
	var (
		err error
		cfg *utils.Config
		)
	if len(configs)==1 {
		cfg = configs[0]
	}else{
		cfg, err = utils.GetConfig()
	}
	if err != nil {
		log.Error(err)
	}

	auth := &Authorization{
		ClusterName: cfg.SelectedCluster,
		ServerUrl:   cfg.CurrentServer(),
		Config:      cfg,
		TokenType:   "Bearer",
	}

	user:= cfg.CurrentUser()
	if user != nil {
		auth.Username = user.Name
		auth.Token = user.Token
	}
	return auth
}