package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type Instances struct {
	Pages int                      `json:"pages"`
	Size  int                      `json:"size"`
	Total int                      `json:"total"`
	List  []ApplicationInstanceDTO `json:"list"`
}

type InstanceValues struct {
	Values string `json:"yaml"`
}

type ApplicationInstanceDTO struct {
	ID                  int    `json:"id"`
	Code                string `json:"code"`
	Status              string `json:"status"`
	PodCount            int    `json:"podCount"`
	PodRunningCount     int    `json:"podRunningCount"`
	AppServiceId        int    `json:"appServiceId"`
	AppServiceName      string `json:"appServiceName"`
	AppServiceVersionId int    `json:"AppServiceVersionId"`
	VersionName         string `json:"VersionName"`
	ObjectVersionNumber int    `json:"objectVersionNumber"`
	Connect             bool   `json:"connect"`
	CommandVersion      string `json:"commandVersion"`
	CommandVersionID    int    `json:"commandVersionId"`
	CommandType         string `json:"commandType"`
	CommandStatus       string `json:"commandStatus"`
	Error               string `json:"error"`
	ProjectID           int    `json:"projectId"`
	ClusterId           int    `json:"clusterId"`
}

type InstancePostInfo struct {
	AppServiceId        int    `json:"appServiceId"`
	AppServiceVersionId int    `json:"appServiceVersionId"`
	EnvironmentId       int    `json:"environmentId"`
	InstanceName        string `json:"instanceName"`
	Type                string `json:"type"`
	Values              string `json:"values"`
}

type InstanceInfo struct {
	Id             int
	Code           string
	VersionName    string
	AppServiceName string
	Status         string
	Pod            string
}

func PrintEnvInstanceInfo(contents []InstanceInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Id", "ApplicationName", "InstanceCode", "Status", "Pods")
	for _, r := range contents {
		table.AddRow(r.Id, r.AppServiceName, r.Code, r.Status, r.Pod)
	}
	fmt.Fprintf(out, table.String())
}

func PrintCreateInstanceInfo(contents []InstanceInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Code", "Status")
	for _, r := range contents {
		table.AddRow(r.Code, r.Status)
	}
	fmt.Fprintf(out, table.String())
}
