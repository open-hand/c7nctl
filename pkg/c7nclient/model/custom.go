package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Customs struct {
	Pages int      `json:"pages"`
	Size  int      `json:"size"`
	Total int      `json:"total"`
	List  []Custom `json:"list"`
}

type Custom struct {
	Id              int    `json:"id"`
	ProjectId       int    `json:"projectId"`
	EnvId           int    `json:"envId"`
	ClusterId       int    `json:"clusterId"`
	CreationDate    string `json:"creationDate"`
	LastUpdateDate  string `json:"lastUpdateDate"`
	EnvCode         string `json:"EnvCode"`
	ResourceContent string `json:"resourceContent"`
	K8sKind         string `json:"k8sKind"`
	CommandStatus   string `json:"commandStatus"`
	Name            string `json:"name"`
	Description     string `json:"description"`
}

type CustomInfo struct {
	Id             int
	LastUpdateDate string
	K8sKind        string
	Status         string
	Name           string
}

func PrintCustomInfos(customInfos []CustomInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "Name", "Type", "UpdateTime", "Status")
	for _, r := range customInfos {
		table.AddRow(r.Id, r.Name, r.K8sKind, r.LastUpdateDate, r.Status)
	}
	fmt.Fprintf(out, table.String())
}
