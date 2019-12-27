package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Services struct {
	Pages int             `json:"pages"`
	Size  int             `json:"size"`
	Total int             `json:"total"`
	List  []DevOpsService `json:"list"`
}

type DevOpsService struct {
	ID                  int               `json:"id"`
	Name                string            `json:"name"`
	Status              string            `json:"status"`
	EnvID               int               `json:"envId"`
	EnvName             string            `json:"envName"`
	Type                string            `json:"type"`
	EnvStatus           bool              `json:"envStatus"`
	AppServiceId        int               `json:"appServiceId"`
	AppServiceProjectId int               `json:"appServiceProjectId"`
	AppServiceName      string            `json:"appServiceName"`
	Dns                 string            `json:"dns"`
	Labels              map[string]string `json:"labels"`
	CommandType         string            `json:"commandType"`
	CommandStatus       string            `json:"commandStatus"`
	Error               string            `json:"error"`
	LoadBalanceIP       string            `json:"loadBalanceIp"`
	Target              struct {
		Instances []ApplicationInstanceDTO `json:"instances"`
		Labels    map[string]string        `json:"labels"`
		EndPoints map[string][]map[string]interface{} `json:"endPoints"`
	} `json:"target"`
	Config struct {
		ExternalIps interface{} `json:"externalIps"`
		Ports       []struct {
			Name       string      `json:"name"`
			Port       int         `json:"port"`
			NodePort   interface{} `json:"nodePort"`
			Protocol   string      `json:"protocol"`
			TargetPort string      `json:"targetPort"`
		} `json:"ports"`
	} `json:"config"`
}

type DevOpsServiceInfo struct {
	Id         int
	Name       string
	Type       string
	TargetType string
	Target     string
	Status     string
}

type ServicePostInfo struct {
	Name       string                        `json:"name"`
	AppID      int                           `json:"appId"`
	Instances  []string                      `json:"instances"`
	EnvID      int                           `json:"envId"`
	ExternalIP string                        `json:"externalIp"`
	Ports      []ServicePort                 `json:"ports"`
	Type       string                        `json:"type"`
	EndPoints  map[string][]EndPointPortInfo `json:"endPoints"`
	Selectors  map[string]string             `json:"selectors"`
}

type EndPointPortInfo struct {
	Name string `json:"name"`
	Port int32  `json:"port"`
}

type ServicePort struct {
	Port       int32              `json:"port"`
	TargetPort intstr.IntOrString `json:"targetPort"`
	NodePort   int32              `json:"nodePort"`
}

func PrintServiceInfo(contents []DevOpsServiceInfo, out io.Writer) {
	table := uitable.New()
	table.MaxColWidth = 60
	table.AddRow("Id", "Name", "Type", "TargetType", "Target", "Status")
	for _, r := range contents {
		table.AddRow(r.Id, r.Name, r.Type, r.TargetType, r.Target, r.Status)
	}
	fmt.Fprintf(out, table.String())
}
