package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type DevopsEnvInstance struct {
	DevopsEnvPreviewApp []DevopsEnvPreviewApp `json:"devopsEnvPreviewAppDTOS"`
}

type InstanceValues struct {
	Values  string `json:"yaml"`
}

type DevopsEnvPreviewApp struct {
	AppName   string `json:"appName"`
	AppCode   string `json:"appCode"`
	ProjectID int    `json:"projectId"`
	ApplicationInstanceDTOS []ApplicationInstanceDTO  `json:"applicationInstanceDTOS"`
}

type ApplicationInstanceDTO struct {
	ID                  int    `json:"id"`
	AppID               int    `json:"appId"`
	EnvID               int    `json:"envId"`
	AppVersionID        int    `json:"appVersionId"`
	Code                string `json:"code"`
	AppName             string `json:"appName"`
	AppVersion          string `json:"appVersion"`
	EnvCode             string `json:"envCode"`
	EnvName             string `json:"envName"`
	Status              string `json:"status"`
	PodCount            int    `json:"podCount"`
	PodRunningCount     int    `json:"podRunningCount"`
	CommandStatus       string `json:"commandStatus"`
	CommandType         string `json:"commandType"`
	CommandVersion      string `json:"commandVersion"`
	CommandVersionID    int    `json:"commandVersionId"`
	Error               string `json:"error"`
	ObjectVersionNumber int    `json:"objectVersionNumber"`
	ProjectID           int    `json:"projectId"`
	Connect             bool   `json:"connect"`
}

type InstancePostInfo struct {
	AppVersionId int `json:"appVersionId"`
	EnvironmentId int `json:"environmentId"`
	AppId int `json:"appId"`
	InstanceName string `json:"instanceName"`
	Values string `json:"values"`
	Type string `json:"type"`
	IsNotChange bool `json:"isNotChange"`
}

type EnvInstanceInfo struct {
	AppName string
	Id       int
	AppCode string
	InstanceCode string
	Version string
	Status  string
	PodPreviewCount string
}


func PrintEnvInstanceInfo(contents []EnvInstanceInfo, out io.Writer)  {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Id","Application","Code","Status","Pods")
	for _, r := range contents {
		app := fmt.Sprintf("%s(%s)", r.AppName, r.AppCode)
		table.AddRow(r.Id, app, r.InstanceCode, r.Status, r.PodPreviewCount)
	}
	fmt.Fprintf(out,table.String())
}



func PrintCreateInstanceInfo(contents []EnvInstanceInfo, out io.Writer)  {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Code","Status")
	for _, r := range contents {
		table.AddRow(r.InstanceCode, r.Status)
	}
	fmt.Fprintf(out,table.String())
}