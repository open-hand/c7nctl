package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)


type AppVersions struct {
	TotalPages       int `json:"totalPages"`
	TotalElements    int `json:"totalElements"`
	NumberOfElements int `json:"numberOfElements"`
	Size             int `json:"size"`
	Number           int `json:"number"`
	Content          []AppVersion  `json:"content"`
	Empty bool `json:"empty"`
}

type AppVersion struct {
	ID           int         `json:"id"`
	Version      string      `json:"version"`
	Commit       string      `json:"commit"`
	AppName      string      `json:"appName"`
	AppCode      string      `json:"appCode"`
	AppID        int         `json:"appId"`
	AppStatus    bool        `json:"appStatus"`
	CreationDate string      `json:"creationDate"`
	Permission   interface{} `json:"permission"`
}

type AppVersionInfo struct {
	Version string
	AppName string
	AppCode string
	CreationDate string
}

func PrintAppVersionInfo(appVersionInfos []AppVersionInfo, out io.Writer)  {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Version","AppName","AppCode","CreationDate")
	for _, r := range appVersionInfos {
		table.AddRow(r.Version, r.AppName, r.AppCode, r.CreationDate)
	}
	fmt.Fprintf(out,table.String())
}