package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Project  struct {
	ID                  int         `json:"id"`
	Name                string      `json:"name"`
	OrganizationID      int         `json:"organizationId"`
	Code                string      `json:"code"`
	Enabled             bool        `json:"enabled"`
	ObjectVersionNumber int         `json:"objectVersionNumber"`
	TypeName            interface{} `json:"typeName"`
	Type                interface{} `json:"type"`
	ImageURL            interface{} `json:"imageUrl"`
}

type ProjectInfo struct {
	Name string
	Code string
}

func PrintProInfo(proInfos []ProjectInfo, out io.Writer)  {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Code")
	for _, r := range proInfos {
		table.AddRow(r.Name, r.Code)
	}
	fmt.Fprintf(out, table.String())
}
