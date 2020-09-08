package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Organization struct {
	ID                  int         `json:"id"`
	Name                string      `json:"name"`
	Code                string      `json:"code"`
	ObjectVersionNumber int         `json:"objectVersionNumber"`
	Enabled             bool        `json:"enabled"`
	ProjectCount        interface{} `json:"projectCount"`
	ImageURL            interface{} `json:"imageUrl"`
	OwnerLoginName      interface{} `json:"ownerLoginName"`
	OwnerRealName       interface{} `json:"ownerRealName"`
	OwnerPhone          interface{} `json:"ownerPhone"`
	OwnerEmail          interface{} `json:"ownerEmail"`
	Projects            interface{} `json:"projects"`
	Roles               interface{} `json:"roles"`
	UserID              int         `json:"userId"`
	Address             interface{} `json:"address"`
	Into                bool        `json:"into"`
}

type OrganizationInfo struct {
	Name string
	Code string
}

func PrintOrgInfo(orgInfos []OrganizationInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("GetName", "Code")
	for _, r := range orgInfos {
		table.AddRow(r.Name, r.Code)
	}
	fmt.Fprintf(out, table.String())
}
