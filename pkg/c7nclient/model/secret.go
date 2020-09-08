package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Secrets struct {
	Pages int      `json:"pages"`
	Size  int      `json:"size"`
	Total int      `json:"total"`
	List  []Secret `json:"list"`
}

type Secret struct {
	Id             int               `json:"id"`
	EnvId          int               `json:"envId"`
	CommandStatus  string            `json:"commandStatus"`
	Error          string            `json:"error"`
	Name           string            `json:"name"`
	Key            []string          `json:"key"`
	Value          map[string]string `json:"value"`
	Description    string            `json:"description"`
	CreationDate   string            `json:"creationDate"`
	LastUpdateDate string            `json:"lastUpdateDate"`
}

type SecretInfo struct {
	Id             int
	Status         string
	Name           string
	Key            string
	LastUpdateDate string
}

type SecretPostInfo struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	EnvID       int               `json:"envId"`
	Type        string            `json:"type"`
	Value       map[string]string `json:"value"`
}

func PrintSecretInfos(secretInfos []SecretInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "GetName", "Key", "UpdateTime", "Status")
	for _, r := range secretInfos {
		table.AddRow(r.Id, r.Name, r.Key, r.LastUpdateDate, r.Status)
	}
	fmt.Fprintf(out, table.String())
}
