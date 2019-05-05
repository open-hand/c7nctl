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

type DevOpsEnvs struct {
	DevopsEnvGroupID        int         `json:"devopsEnvGroupId"`
	DevopsEnvGroupName      string      `json:"devopsEnvGroupName"`
	DevopsEnviromentRepDTOs []DevOpsEnv `json:"devopsEnviromentRepDTOs"`
}

type DevOpsEnv struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Code               string `json:"code"`
	ClusterName        string `json:"clusterName"`
	ClusterID          int    `json:"clusterId"`
	Sequence           int    `json:"sequence"`
	DevopsEnvGroupID   int    `json:"devopsEnvGroupId"`
	GitlabEnvProjectId int    `json:"gitlabEnvProjectId"`
	Permission         bool   `json:"permission"`
	Connect            bool   `json:"connect"`
	Synchro            bool   `json:"synchro"`
	Failed             bool   `json:"failed"`
	Active             bool   `json:"active"`
}

type AuthEnv struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	ClusterID   int    `json:"clusterId"`
	Sequence    int    `json:"sequence"`
	Permission  bool   `json:"permission"`
	Active      bool   `json:"active"`
	Synchro     bool   `json:"synchro"`
	Failed      bool   `json:"failed"`
	Connect     bool   `json:"connect"`
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
	table.AddRow("Name", "Code", "Status", "Group", "Cluster")
	for _, r := range envs {
		table.AddRow(r.Name, r.Code, r.Status, "default", r.Cluster)
	}
	fmt.Fprintf(out, table.String())
}

func PrintAuthEnvInfo(envs []EnvInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Id", "Name", "Code", "SyncStatus")
	for _, r := range envs {
		table.AddRow(r.Id, r.Name, r.Code, r.SyncStatus)
	}
	fmt.Fprintf(out, table.String())
}
