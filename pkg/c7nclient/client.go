package c7nclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"github.com/gosuri/uitable"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var Client C7NClient

type C7NClient struct {
	BaseURL   string
	httpClient *http.Client
	Token     string
	config    *C7NPlatformContext
}


func InitClient(config *C7NPlatformContext) {
	Client = C7NClient{
		BaseURL: config.Server,
		httpClient: &http.Client{

		},
		Token: config.Token,
		config: config,
	}
}

func (c *C7NClient) ListEnvs(out io.Writer,) {
	if c.config.ProjectId == -1 {
		fmt.Printf("Set project Id")
		return
	}
	paras := make(map[string]interface{})
	paras["active"]="true"
	req,err := c.newRequest("GET",fmt.Sprintf("/devops/v1/projects/%d/envs/groups",c.config.ProjectId,),paras,nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var devOpsEnvs = []model.DevOpsEnvs{}
	_,err = c.do(req,&devOpsEnvs)
	if err != nil {
		fmt.Printf("request err:%v",err)
		return

	}

	envInfos := []model.EnvInfo{}
	for _,devOpsEnv := range devOpsEnvs[0].DevopsEnviromentRepDTOs {
		var status string
		if devOpsEnv.Failed {
			status = "Failed"
		} else if devOpsEnv.Connect  {
			status = "Connected"
		} else {
			status = "Disconnected"
		}
		envInfo := model.EnvInfo{
			Name: devOpsEnv.Name,
			Status: status,
			Code: devOpsEnv.Code,
			Cluster: devOpsEnv.ClusterName,
			Group: devOpsEnvs[0].DevopsEnvGroupName,
		}
		envInfos = append(envInfos, envInfo)
	}
	model.PrintEnvInfo(envInfos,out)

}

func (c *C7NClient) QuerySelf(out io.Writer,) {

	req,err := c.newRequest("GET","/iam/v1/users/self",nil,nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var userInfo *model.UserInfo
	resp,err := c.do(req,userInfo)
	if err != nil {
		fmt.Printf("request err:%v",err)
		return

	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return
	}


}


func (c *C7NClient) newRequest(method, path string,paras map[string]interface{}, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	base,_ := url.Parse(c.BaseURL)
	u := base.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	if paras != nil {
		for key,value := range paras {
			q.Add(key, fmt.Sprintf("%v",value))
		}
	}
	req.URL.RawQuery = q.Encode()
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "bearer "+c.Token)
	return req, nil
}

func (c *C7NClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("request error with: %s", bodyString)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func (c *C7NClient) printContextInfo()  {
	fmt.Printf("organization: %s(%s) project: %s(%s)", c.config.Organization, c.config.OrganizationCode,
		c.config.Project, c.config.ProjectCode)
}



func PrintConfigInfo(config C7NPlatformContext, out io.Writer)  {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Name","Server","Organization","Project","Token")
	table.AddRow(config.Name, config.Server,
		fmt.Sprintf("%s(%s)",config.Organization,config.OrganizationCode),
		fmt.Sprintf("%s(%s)",config.Project,config.ProjectCode),  config.Token)
	fmt.Fprintf(out,table.String())

}