package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Apps struct {
	TotalPages       int   `json:"totalPages"`
	TotalElements    int   `json:"totalElements"`
	NumberOfElements int   `json:"numberOfElements"`
	Size             int   `json:"size"`
	Number           int   `json:"number"`
	Content          []App `json:"content"`
	Empty            bool  `json:"empty"`
}

type App struct {
	ID                    int         `json:"id"`
	Name                  string      `json:"name"`
	Code                  string      `json:"code"`
	ProjectID             int         `json:"projectId"`
	ApplicationTemplateID int         `json:"applicationTemplateId"`
	RepoURL               string      `json:"repoUrl"`
	PublishLevel          interface{} `json:"publishLevel"`
	Contributor           interface{} `json:"contributor"`
	Description           interface{} `json:"description"`
	SonarURL              interface{} `json:"sonarUrl"`
	Type                  string      `json:"type"`
	Permission            bool        `json:"permission"`
	Synchro               bool        `json:"synchro"`
	Fail                  interface{} `json:"fail"`
	Active                bool        `json:"active"`
}

type AppInfo struct {
	Type    string
	Name    string
	Code    string
	RepoURL string
	Status  string
}

type AppPostInfo struct {
	Name                  string `json:"name"`
	Code                  string `json:"code"`
	Type                  string `json:"type"`
	ApplicationTemplateId int    `json:"applicationTemplateId"`
	IsSkipCheckPermission bool   `json:"isSkipCheckPermission"`
}

func PrintAppInfo(appInfos []AppInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Type", "Name", "Code", "RepoUrl", "Status")
	for _, r := range appInfos {
		table.AddRow(r.Type, r.Name, r.Code, r.RepoURL, r.Status)
	}
	fmt.Fprintf(out, table.String())
}
