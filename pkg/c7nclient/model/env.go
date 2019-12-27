package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type EnvInfo struct {
	Name       string
	Code       string
	Group      string
	Status     string
	Cluster    string
	SyncStatus bool
	Id         int
}

type DevopsGroupInfo struct {
	Id        int    `json:"id"`
	ProjectId int    `json:"projectId"`
	Name      string `json:"name"`
}

type DevOpsEnvs struct {
	DevopsEnvGroupID         int         `json:"devopsEnvGroupId"`
	DevopsEnvGroupName       string      `json:"devopsEnvGroupName"`
	DevopsEnvironmentRepDTOs []DevOpsEnv `json:"devopsEnvironmentRepDTOs"`
}

type DevOpsEnv struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Code               string `json:"code"`
	GitlabEnvProjectId int    `json:"gitlabEnvProjectId"`
	ClusterName        string `json:"clusterName"`
	ClusterID          int    `json:"clusterId"`
	DevopsEnvGroupID   int    `json:"devopsEnvGroupId"`
	Permission         bool   `json:"permission"`
	Connect            bool   `json:"connect"`
	Synchro            bool   `json:"synchro"`
	Failed             bool   `json:"failed"`
	Active             bool   `json:"active"`
}

type EnvPostInfo struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	ClusterId   int    `json:"clusterId"`
}

type EnvSyncStatus struct {
	DevopsSyncCommit string `json:"devopsSyncCommit"`
	AgentSyncCommit  string `json:"agentSyncCommit"`
	SagaSyncCommit   string `json:"sagaSyncCommit"`
	CommitURL        string `json:"commitUrl"`
}

func PrintEnvInfo(envs []EnvInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Id", "Name", "Code", "Status", "Group", "Cluster")
	for _, r := range envs {
		if r.Group == "" {
			r.Group = "default"
		}
		table.AddRow(r.Id, r.Name, r.Code, r.Status, r.Group, r.Cluster)
	}
	fmt.Fprintf(out, table.String())
}

func PrintEnvGroupInfo(groups []DevopsGroupInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Id", "Name")
	for _, r := range groups {
		table.AddRow(r.Id, r.Name)
	}
	fmt.Fprintf(out, table.String())
}
