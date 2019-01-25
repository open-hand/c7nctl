package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Clusters struct {
	TotalPages       int       `json:"totalPages"`
	TotalElements    int       `json:"totalElements"`
	NumberOfElements int       `json:"numberOfElements"`
	Size             int       `json:"size"`
	Number           int       `json:"number"`
	Content          []Cluster `json:"content"`
	Empty            bool      `json:"empty"`
}

type Cluster struct {
	ID                         int         `json:"id"`
	Name                       string      `json:"name"`
	SkipCheckProjectPermission bool        `json:"skipCheckProjectPermission"`
	Code                       string      `json:"code"`
	Connect                    bool        `json:"connect"`
	Upgrade                    bool        `json:"upgrade"`
	UpgradeMessage             interface{} `json:"upgradeMessage"`
	Description                string      `json:"description"`
}

type Node struct {
	Type                    string `json:"type"`
	NodeName                string `json:"nodeName"`
	Status                  string `json:"status"`
	CreateTime              string `json:"createTime"`
	CPUTotal                string `json:"cpuTotal"`
	CPURequest              string `json:"cpuRequest"`
	CPULimit                string `json:"cpuLimit"`
	CPURequestPercentage    string `json:"cpuRequestPercentage"`
	CPULimitPercentage      string `json:"cpuLimitPercentage"`
	MemoryTotal             string `json:"memoryTotal"`
	MemoryRequest           string `json:"memoryRequest"`
	MemoryLimit             string `json:"memoryLimit"`
	MemoryRequestPercentage string `json:"memoryRequestPercentage"`
	MemoryLimitPercentage   string `json:"memoryLimitPercentage"`
	PodTotal                int    `json:"podTotal"`
	PodCount                int    `json:"podCount"`
	PodPercentage           string `json:"podPercentage"`
}

type Nodes struct {
	TotalPages       int    `json:"totalPages"`
	TotalElements    int    `json:"totalElements"`
	NumberOfElements int    `json:"numberOfElements"`
	Size             int    `json:"size"`
	Number           int    `json:"number"`
	Content          []Node `json:"content"`
	Empty            bool   `json:"empty"`
}

type NodeInfo struct {
	Status        string
	NodeName      string
	NodeType      string
	CpuLimit      string
	CpuRequest    string
	MemoryLimit   string
	MemoryRequest string
	CreationDate  string
}

type ClusterInfo struct {
	Name        string
	Code        string
	Description string
	Status      string
}

func PrintClusterInfo(clusterInfos []ClusterInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Code", "Description", "Status")
	for _, r := range clusterInfos {
		table.AddRow(r.Name, r.Code, r.Description, r.Status)
	}
	fmt.Fprintf(out, table.String())
}

func PrintNodeInfo(nodeInfos []NodeInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Status", "Name", "Type", "Cpu Request", "Cpu Limit", "Memory Request", "Memory Limit", "CreationDate")
	for _, r := range nodeInfos {
		table.AddRow(r.Status, r.NodeName, r.NodeType, r.CpuRequest, r.CpuLimit, r.MemoryRequest, r.MemoryLimit, r.CreationDate)
	}
	fmt.Fprintf(out, table.String())
}
