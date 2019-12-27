package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Pvcs struct {
	Pages int   `json:"pages"`
	Size  int   `json:"size"`
	Total int   `json:"total"`
	List  []Pvc `json:"list"`
}

type Pvc struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	EnvId           int    `json:"envId"`
	ProjectId       int    `json:"projectId"`
	PvId            int    `json:"pvId"`
	PvName          string `json:"pvName"`
	AccessModes     string `json:"accessModes"`
	Status          string `json:"status"`
	RequestResource string `json:"requestResource"`
	EnvCode         string `json:"EnvCode"`
}

type PvcInfo struct {
	Id              int
	Name            string
	PvName          string
	AccessModes     string
	RequestResource string
	Status          string
}

type PvcPostInfo struct {
	Name            string `json:"name"`
	PvName          string `json:"pvName"`
	AccessModes     string `json:"accessModes"`
	RequestResource string `json:"requestResource"`
	EnvID           int    `json:"envId"`
	ClusterId       int    `json:"clusterId"`
}

func PrintPvcInfos(pvcInfos []PvcInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "Name", "Status", "PVName", "AccessModes", "RequestResource")
	for _, r := range pvcInfos {
		table.AddRow(r.Id, r.Name, r.Status, r.PvName, r.AccessModes, r.RequestResource)
	}
	fmt.Fprintf(out, table.String())
}
