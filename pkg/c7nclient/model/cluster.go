package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Clusters struct {
	Pages int       `json:"pages"`
	Size  int       `json:"size"`
	Total int       `json:"total"`
	List  []Cluster `json:"list"`
}

type Cluster struct {
	ID                         int         `json:"id"`
	Name                       string      `json:"name"`
	Code                       string      `json:"code"`
	Connect                    bool        `json:"connect"`
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
	Pages int       `json:"pages"`
	Size  int       `json:"size"`
	Total int       `json:"total"`
	List  []Node `json:"list"`
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
	Status      string
}

type ClusterPostInfo struct {
	Name                       string `json:"name"`
	Code                       string `json:"code"`
	Description                string `json:"description"`
	SkipCheckProjectPermission bool   `json:"skipCheckProjectPermission"`
}

func PrintClusterInfo(clusterInfos []ClusterInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Code", "Status")
	for _, r := range clusterInfos {
		table.AddRow(r.Name, r.Code, r.Status)
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
