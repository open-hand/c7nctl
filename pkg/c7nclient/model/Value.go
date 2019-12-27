package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Values struct {
	Pages int     `json:"pages"`
	Size  int     `json:"size"`
	Total int     `json:"total"`
	List  []Value `json:"list"`
}

type Value struct {
	Id             int    `json:"id"`
	Value          string `json:"value"`
	ProjectId      int    `json:"projectId"`
	EnvId          int    `json:"envId"`
	AppServiceId   int    `json:"appServiceId"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	EnvName        string `json:"envName"`
	AppServiceName string `json:"AppServiceName"`
}

type ValueInfo struct {
	Id             int
	Name           string
	Description    string
	AppServiceName string
	Environment    string
}

func PrintValueInfo(valueInfos []ValueInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "Name", "Description", "AppServiceName", "Environment")
	for _, r := range valueInfos {
		table.AddRow(r.Id, r.Name, r.Description, r.AppServiceName, r.Environment)
	}
	fmt.Fprintf(out, table.String())
}
