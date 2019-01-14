package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type UserInfo struct {
	ID                   int    `json:"id"`
	OrganizationID       int    `json:"organizationId"`
	OrganizationName     string `json:"organizationName"`
	OrganizationCode     string `json:"organizationCode"`
	LoginName            string `json:"loginName"`
	Email                string `json:"email"`
	RealName             string `json:"realName"`
	Phone                string `json:"phone"`
	InternationalTelCode string `json:"internationalTelCode"`
	ImageURL             string `json:"imageUrl"`
	Language             string `json:"language"`
	TimeZone             string `json:"timeZone"`
	Locked               bool   `json:"locked"`
	Ldap                 bool   `json:"ldap"`
	Enabled              bool   `json:"enabled"`
	Admin                bool   `json:"admin"`
	ObjectVersionNumber  int    `json:"objectVersionNumber"`
}


type EnvInfo struct {
	Name string
	Code string
	Group string
	Status string
	Cluster string
}



type DevOpsEnvs struct {
	DevopsEnvGroupID        int `json:"devopsEnvGroupId"`
	DevopsEnvGroupName      string `json:"devopsEnvGroupName"`
	DevopsEnviromentRepDTOs []DevOpsEnv `json:"devopsEnviromentRepDTOs"`
}

type DevOpsEnv struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Code             string      `json:"code"`
	ClusterName      string      `json:"clusterName"`
	ClusterID        int         `json:"clusterId"`
	Sequence         int         `json:"sequence"`
	DevopsEnvGroupID int         `json:"devopsEnvGroupId"`
	Permission       bool 		 `json:"permission"`
	Connect          bool        `json:"connect"`
	Synchro          bool        `json:"synchro"`
	Failed           bool        `json:"failed"`
	Active           bool        `json:"active"`
}

func PrintEnvInfo(envs []EnvInfo, out io.Writer)  {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Name","Code","Status","Group","Cluster")
	for _, r := range envs {
		table.AddRow(r.Name, r.Code, r.Status,"default",  r.Cluster)
	}
	fmt.Fprintf(out,table.String())

}



func formatText(results []EnvInfo, colWidth uint) string {

	table := uitable.New()
	table.MaxColWidth = colWidth
	table.AddRow("Name","Code","Status","Group","Cluster")
	for _, r := range results {
		table.AddRow(r.Name, r.Code, r.Status,"default",  r.Cluster)
	}

	return fmt.Sprintf("%s" ,table.String())
}