package app

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Persistence struct {
	Name string
	Path string
	Mode int
	Size string
}

func (p *Persistence) GetCapacity() v1.ResourceList {
	capacity := make(map[v1.ResourceName]resource.Quantity)
	q := resource.MustParse(p.Size)
	capacity["storage"] = q
	return capacity
}

type InstallCmd struct {
	Values  []string
	Chart   string
	RepoURL string
	Version string
}

//func (i *InstallCmd) Run() error {
//	return InstallRelease(i.Values, i.Chart, i.RepoURL, i.Version)
//}
