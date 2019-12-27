package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type ConfigMaps struct {
	Pages int         `json:"pages"`
	Size  int         `json:"size"`
	Total int         `json:"total"`
	List  []ConfigMap `json:"list"`
}

type ConfigMap struct {
	Id             int      `json:"id"`
	EnvId          int      `json:"envId"`
	CommandStatus  string   `json:"commandStatus"`
	EnvCode        string   `json:"envCode"`
	EnvStatus      string   `json:"envStatus"`
	Error          string   `json:"error"`
	Name           string   `json:"name"`
	Key            []string `json:"key"`
	Value          []string `json:"value"`
	Description    string   `json:"description"`
	CreationDate   string   `json:"creationDate"`
	LastUpdateDate string   `json:"lastUpdateDate"`
}

type ConfigMapInfo struct {
	Id             int
	Name           string
	Key            string
	Status         string
	LastUpdateDate string
}

type ConfigMapPostInfo struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	EnvID       int               `json:"envId"`
	Type        string            `json:"type"`
	Value       map[string]string `json:"value"`
}

func PrintConfigMapInfos(configMapInfos []ConfigMapInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "Name", "Key", "UpdateTime", "Status")
	for _, r := range configMapInfos {
		table.AddRow(r.Id, r.Name, r.Key, r.LastUpdateDate, r.Status)
	}
	fmt.Fprintf(out, table.String())
}
