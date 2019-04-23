package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DevOpsService struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	EnvID        int    `json:"envId"`
	EnvName      string `json:"envName"`
	Type         string `json:"type"`
	EnvStatus    bool   `json:"envStatus"`
	AppID        int    `json:"appId"`
	AppProjectID int    `json:"appProjectId"`
	AppName      string `json:"appName"`
	Target       struct {
		AppInstance []struct {
			ID             string `json:"id"`
			Code           string `json:"code"`
			InstanceStatus string `json:"instanceStatus"`
		} `json:"appInstance"`
		Labels map[string]string `json:"labels"`
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
	CommandType   string `json:"commandType"`
	CommandStatus string `json:"commandStatus"`
	Error         string `json:"error"`
	LoadBalanceIP string `json:"loadBalanceIp"`
}

type DevOpsServicePage struct {
	TotalPages       int             `json:"totalPages"`
	TotalElements    int             `json:"totalElements"`
	NumberOfElements int             `json:"numberOfElements"`
	Size             int             `json:"size"`
	Number           int             `json:"number"`
	Content          []DevOpsService `json:"content"`
	Empty            bool            `json:"empty"`
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
	Name        string                        `json:"name"`
	AppID       int                           `json:"appId"`
	AppInstance []string                      `json:"appInstance"`
	EnvID       int                           `json:"envId"`
	ExternalIP  string                        `json:"externalIp"`
	Ports       []ServicePort                 `json:"ports"`
	Label       map[string]string             `json:"label"`
	Type        string                        `json:"type"`
	EndPoints   map[string][]EndPointPortInfo `json:"endPoints"`
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
