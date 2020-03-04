package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Pvs struct {
	Pages int  `json:"pages"`
	Size  int  `json:"size"`
	Total int  `json:"total"`
	List  []Pv `json:"list"`
}

type Pv struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	AccessModes     string `json:"accessModes"`
	PvcName         string `json:"pvcName"`
	ClusterName     string `json:"clusterName"`
	Type            string `json:"type"`
	RequestResource string `json:"requestResource"`
	Status          string `json:"status"`
}

type PvInfo struct {
	Id              int
	Name            string
	PvcName         string
	ClusterName     string
	AccessModes     string
	RequestResource string
	Status          string
	Type            string
}

type PvPostInfo struct {
	Name                       string `json:"name"`
	AccessModes                string `json:"accessModes"`
	ClusterId                  int    `json:"clusterId"`
	RequestResource            string `json:"requestResource"`
	Type                       string `json:"type"`
	SkipCheckProjectPermission bool   `json:"skipCheckProjectPermission,bool"`
	ValueConfig                string `json:"valueConfig"`
	ProjectIds                 []int  `json:"projectIds"`
}

func PrintPvInfos(pvInfos []PvInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Id", "Name", "Type", "Status", "PVCName", "ClusterName", "AccessModes", "RequestResource")
	for _, r := range pvInfos {
		table.AddRow(r.Id, r.Name, r.Type, r.Status, r.PvcName, r.ClusterName, r.AccessModes, r.RequestResource)
	}
	fmt.Fprintf(out, table.String())
}
