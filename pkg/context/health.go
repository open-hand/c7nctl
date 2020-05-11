package context

import (
	syserr "errors"
	"fmt"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Health struct {
	HttpGet   []HttpGetCheck `yaml:"httpGet"`
	Socket    []SocketCheck
	PodStatus []PodCheck `yaml:"podStatus"`
}

type PodCheck struct {
	Name      string
	Status    core_v1.PodPhase
	Namespace string
	Client    kubernetes.Interface
}

type SocketCheck struct {
	Name string
	Host string
	Port int32
	Path string
}

type HttpGetCheck struct {
	Name string
	Host string
	Port int32
	Path string
}

func (p *PodCheck) MustRunning() error {
	po, err := p.Client.CoreV1().Pods(p.Namespace).Get(p.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	if status := po.Status.Phase; status != p.Status {
		return syserr.New(fmt.Sprintf("[ %s ] pod status is %s, need %s", p.Name, status, p.Status))
	}

	return nil
}
