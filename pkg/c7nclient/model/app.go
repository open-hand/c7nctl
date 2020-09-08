package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Apps struct {
	Pages int   `json:"pages"`
	Size  int   `json:"size"`
	Total int   `json:"total"`
	List  []App `json:"list"`
}

type App struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	Code            string      `json:"code"`
	ProjectID       int         `json:"projectId"`
	GitlabProjectId int         `json:"gitlabProjectId"`
	RepoURL         string      `json:"repoUrl"`
	PublishLevel    interface{} `json:"publishLevel"`
	Contributor     interface{} `json:"contributor"`
	Description     interface{} `json:"description"`
	SonarURL        interface{} `json:"sonarUrl"`
	Type            string      `json:"type"`
	Synchro         bool        `json:"synchro"`
	Fail            interface{} `json:"fail"`
	Active          bool        `json:"active"`
}

type AppInfo struct {
	Id      int
	Type    string
	Name    string
	Code    string
	RepoURL string
	Status  string
}

type AppPostInfo struct {
	Name                        string `json:"name"`
	Code                        string `json:"code"`
	Type                        string `json:"type"`
	TemplateAppServiceId        int    `json:"templateAppServiceId,omitempty"`
	TemplateAppServiceVersionId int    `json:"templateAppServiceVersionId,omitempty"`
}

func PrintAppInfo(appInfos []AppInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("AppId", "Type", "GetName", "Code", "RepoUrl", "Status")
	for _, r := range appInfos {
		table.AddRow(r.Id, r.Type, r.Name, r.Code, r.RepoURL, r.Status)
	}
	fmt.Fprintf(out, table.String())
}
