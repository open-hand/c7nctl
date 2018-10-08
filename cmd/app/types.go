package app

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)



type Versions struct {
	Versions []Version
}

func (v *Versions) GetLastStable() Version {
	for _, version := range v.Versions {
		if version.Status == "stable" {
			return version
		}
	}
	return Version{}
}

type Version struct {
	Version string
	Status  string
}

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
