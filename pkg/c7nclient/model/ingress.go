package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Ingresss struct {
	Pages int             `json:"pages"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
	List  []DevopsIngress `json:"list"`
}

type DevopsIngress struct {
	ID            int           `json:"id"`
	AppServiceId  int           `json:"appServiceId"`
	Domain        string        `json:"domain"`
	Name          string        `json:"name"`
	EnvID         int           `json:"envId"`
	EnvName       string        `json:"envName"`
	EnvStatus     bool          `json:"envStatus"`
	Status        string        `json:"status"`
	CertID        int           `json:"certId"`
	CertName      string        `json:"certName"`
	CertStatus    string        `json:"certStatus"`
	PathList      []IngressPath `json:"pathList"`
	CommandType   string        `json:"commandType"`
	CommandStatus string        `json:"commandStatus"`
	Error         string        `json:"error"`
	Usable        bool          `json:"usable"`
}

type IngressPath struct {
	Path          string             `json:"path"`
	ServiceID     int                `json:"serviceId"`
	ServiceName   string             `json:"serviceName"`
	ServiceStatus string             `json:"serviceStatus"`
	ServicePort   intstr.IntOrString `json:"servicePort"`
}

type IngressPostInfo struct {
	Domain   string        `json:"domain"`
	Name     string        `json:"name"`
	EnvID    int           `json:"envId"`
	CertId   int           `json:"certId"`
	PathList []IngressPath `json:"pathList"`
}

type DevOpsIngressInfo struct {
	Id     int
	Name   string
	Status string
	Host   string
	Paths  string
}

func PrintIngressInfo(contents []DevOpsIngressInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Id", "GetName", "Status", "Host", "Paths")
	for _, r := range contents {
		table.AddRow(r.Id, r.Name, r.Status, r.Host, r.Paths)
	}
	fmt.Fprintf(out, table.String())
}
