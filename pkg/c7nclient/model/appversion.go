package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
)

type AppVersions struct {
	Pages int          `json:"pages"`
	Size  int          `json:"size"`
	Total int          `json:"total"`
	List  []AppVersion `json:"list"`
}

type AppVersion struct {
	ID           int    `json:"id"`
	Version      string `json:"version"`
	AppServiceId int    `json:"appServiceId"`
	CreationDate string `json:"creationDate"`
}

type AppVersionInfo struct {
	Id           int
	Version      string
	CreationDate string
}

func PrintAppVersionInfo(appVersionInfos []AppVersionInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("VersionId", "Version", "CreationDate")
	for _, r := range appVersionInfos {
		table.AddRow(r.Id, r.Version, r.CreationDate)
	}
	fmt.Fprintf(out, table.String())
}
