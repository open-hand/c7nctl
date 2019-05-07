package model

import (
	"fmt"
	"github.com/gosuri/uitable"
	"io"
	"reflect"
	"time"
)

type StatefulSetDTOS struct {
	Name            string    `json:"name"`
	ReadyReplicas   int       `json:"readyReplicas"`
	CurrentReplicas int       `json:"currentReplicas"`
	DesiredReplicas int       `json:"desiredReplicas"`
	Age             time.Time `json:"age"`
}

type IngressDTOS struct {
	Name    string    `json:"name"`
	Hosts   string    `json:"hosts"`
	Address string    `json:"address"`
	Ports   string    `json:"ports"`
	Age     time.Time `json:"age"`
}

type DeploymentDTOS struct {
	Name      string    `json:"name"`
	Current   int       `json:"current"`
	UpToDate  int       `json:"upToDate"`
	Desired   int       `json:"desired"`
	Available int       `json:"available"`
	Age       time.Time `json:"age"`
}

type PersistentVolumeClaimDTOS struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	AccessModes string    `json:"accessModes"`
	Capacity    string    `json:"capacity"`
	Age         time.Time `json:"age"`
}

type PodDTOS struct {
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	Ready    int       `json:"ready"`
	Desire   int       `json:"desire"`
	Restarts int       `json:"restarts"`
	Age      time.Time `json:"age"`
}

type ReplicaSetDTOS struct {
	Name    string    `json:"name"`
	Ready   int       `json:"ready"`
	Current int       `json:"current"`
	Desired int       `json:"desired"`
	Age     time.Time `json:"age"`
}

type DaemonSetDTOS struct {
	Name             string    `json:"name"`
	CurrentScheduled int       `json:"currentScheduled"`
	DesiredScheduled int       `json:"desiredScheduled"`
	NumberAvailable  int       `json:"numberAvailable"`
	Age              time.Time `json:"age"`
}

type ServiceDTOS struct {
	Name       string    `json:"name"`
	ClusterIP  string    `json:"clusterIp"`
	Type       string    `json:"type"`
	ExternalIP string    `json:"externalIp"`
	Port       string    `json:"port"`
	TargetPort string    `json:"targetPort"`
	Age        time.Time `json:"age"`
}

type InstanceResources struct {
	StatefulSets           []StatefulSetDTOS           `json:"statefulSetDTOS"`
	Ingresses              []IngressDTOS               `json:"ingressDTOS"`
	Deployments            []DeploymentDTOS            `json:"deploymentDTOS"`
	PersistentVolumeClaims []PersistentVolumeClaimDTOS `json:"persistentVolumeClaimDTOS"`
	Pods                   []PodDTOS                   `json:"podDTOS"`
	ReplicaSets            []ReplicaSetDTOS            `json:"replicaSetDTOS"`
	DaemonSets             []DaemonSetDTOS             `json:"daemonSetDTOS"`
	Services               []ServiceDTOS               `json:"serviceDTOS"`
}

func PrintInstanceResources(resources InstanceResources, out io.Writer) {
	var interfaceSlice = make([]interface{}, len(resources.Pods))
	for i, d := range resources.Pods {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "Pods")

	interfaceSlice = make([]interface{}, len(resources.Services))
	for i, d := range resources.Services {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "Services")

	interfaceSlice = make([]interface{}, len(resources.DaemonSets))
	for i, d := range resources.DaemonSets {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "DaemonSets")

	interfaceSlice = make([]interface{}, len(resources.StatefulSets))
	for i, d := range resources.StatefulSets {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "StatefulSets")

	interfaceSlice = make([]interface{}, len(resources.Deployments))
	for i, d := range resources.Deployments {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "Deployments")

	interfaceSlice = make([]interface{}, len(resources.Ingresses))
	for i, d := range resources.Ingresses {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "Ingresses")

	interfaceSlice = make([]interface{}, len(resources.ReplicaSets))
	for i, d := range resources.ReplicaSets {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "ReplicaSets")

	interfaceSlice = make([]interface{}, len(resources.PersistentVolumeClaims))
	for i, d := range resources.PersistentVolumeClaims {
		interfaceSlice[i] = d
	}
	printObj(interfaceSlice, out, "PersistentVolumeClaims")

}

func printObj(obj []interface{}, out io.Writer, title string) {
	if len(obj) == 0 {
		return
	}
	fmt.Printf("==>%s\n", title)
	table := uitable.New()
	table.MaxColWidth = 60
	printObjHeader(obj[0], table)
	for _, r := range obj {
		t := reflect.ValueOf(r)
		fieldCount := t.NumField()
		tableValue := make([]interface{}, 0, fieldCount)
		for i := 0; i < fieldCount; i++ {
			if t.Field(i).Type().Name() == "Time" {
				createTime, _ := t.Field(i).Interface().(time.Time)
				tableValue = append(tableValue, fmt.Sprintf("%s", duration(createTime)))

			} else {
				tableValue = append(tableValue, fmt.Sprintf("%v", t.Field(i).Interface()))
			}

		}
		table.AddRow(tableValue[0:]...)
	}
	fmt.Fprintf(out, table.String())
	fmt.Printf("\n\n\n")
}

func printObjHeader(obj interface{}, table *uitable.Table) {
	t := reflect.TypeOf(obj)
	fieldCount := t.NumField()
	tableHeader := make([]interface{}, 0, fieldCount)
	for i := 0; i < fieldCount; i++ {
		tableHeader = append(tableHeader, t.Field(i).Name)
	}
	table.AddRow(tableHeader[0:]...)
}

func duration(t time.Time) string {
	now := time.Now()
	d := now.Sub(t)
	if d.Hours() > 24 {
		return fmt.Sprintf("%dd", int(d.Hours())%24)
	} else if d.Hours() > 0 {
		return fmt.Sprintf("%dh", int(d.Hours()))
	} else if d.Minutes() > 0 {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	} else {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
}
