package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)


type AppTemplates struct {
	TotalPages       int `json:"totalPages"`
	TotalElements    int `json:"totalElements"`
	NumberOfElements int `json:"numberOfElements"`
	Size             int `json:"size"`
	Number           int `json:"number"`
	Content []AppTemplate `json:"content"`
	Empty bool `json:"empty"`
}

type AppTemplate struct {
	ID                  int         `json:"id"`
	OrganizationID      interface{} `json:"organizationId"`
	Name                string      `json:"name"`
	Description         string      `json:"description"`
	Code                string      `json:"code"`
	CopyFrom            interface{} `json:"copyFrom"`
	RepoURL             string      `json:"repoUrl"`
	Type                bool        `json:"type"`
	ObjectVersionNumber interface{} `json:"objectVersionNumber"`
	Failed              bool        `json:"failed"`
	Synchro             bool        `json:"synchro"`
}

type AppTemplateInfo struct {
	Name string
	Code string
	RepoUrl string
	Available string
}


func PrintAppTemplateInfo(appTemplates []AppTemplateInfo, out io.Writer)  {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name","Code","RepoUrl","Available")
	for _, r := range appTemplates {
		table.AddRow(r.Name, r.Code, r.RepoUrl,r.Available)
	}
	fmt.Fprintf(out,table.String())
}
