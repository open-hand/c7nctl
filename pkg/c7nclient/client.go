package c7nclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/choerodon/c7n/pkg/c7nclient/model"
	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

var Client C7NClient

type C7NClient struct {
	BaseURL    string
	httpClient *http.Client
	token      string
	config     *C7NPlatformContext
}

func InitClient(config *C7NPlatformContext) {
	Client = C7NClient{
		BaseURL: config.Server,
		httpClient: &http.Client{
		},
		token:  config.Token,
		config: config,
	}
}

func (c *C7NClient) QuerySelf(out io.Writer, ) {

	req, err := c.newRequest("GET", "/iam/v1/users/self", nil, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var userInfo *model.UserInfo
	_, err = c.do(req, userInfo)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
}

func (c *C7NClient) newRequest(method, path string, paras map[string]interface{}, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	base, _ := url.Parse(c.BaseURL)
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
		for key, value := range paras {
			q.Add(key, fmt.Sprintf("%v", value))
		}
	}
	req.URL.RawQuery = q.Encode()
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization", "bearer "+c.token)
	return req, nil
}

func (c *C7NClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	err = c.handleRep(resp)
	if err != nil {
		return resp, err
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func (c *C7NClient) handleRep(resp *http.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}
	if resp.StatusCode == 403 {
		return errors.New("You do not have the permissions!")
	} else {
		return errors.New(resp.Status)
	}
	return nil
}

func (c *C7NClient) printContextInfo() {
	fmt.Printf("organization: %s(%s) project: %s(%s)", c.config.Organization, c.config.OrganizationCode,
		c.config.Project, c.config.ProjectCode)
}

func PrintConfigInfo(config C7NPlatformContext, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Name", "Server", "Organization", "Project", "Token")
	table.AddRow(config.Name, config.Server,
		fmt.Sprintf("%s(%s)", config.Organization, config.OrganizationCode),
		fmt.Sprintf("%s(%s)", config.Project, config.ProjectCode), config.Token)
	fmt.Fprintf(out, table.String())
}
