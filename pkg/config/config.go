package config

import (
	"k8s.io/api/core/v1"
	"fmt"
)

type Config struct {
	Version  string
	Metadata Metadata
	Spec     Spec
}

type Metadata struct {
	Name      string
	Namespace string
}

type Spec struct {
	Persistence Persistence
	Resources   map[string]*Resource
}

type Persistence struct {
	Nfs
}

type Nfs struct {
	Server   string
	RootPath string `yaml:"rootPath"`
}

// todo: remove
type Resources struct {
	Mysql  Resource
	Gitlab Resource
}

type Resource struct {
	Host     string
	Port     int
	Username string
	Password string
	Schema   string
	External bool
}

func (p *Persistence) GetPersistentVolumeSource(subPath string) v1.PersistentVolumeSource {
	if p.Nfs.Server != "" {
		return p.prepareNfsPVS(subPath)
	}
	return v1.PersistentVolumeSource{}
}

func (p *Persistence) prepareNfsPVS(subPath string) v1.PersistentVolumeSource {
	pvs := v1.PersistentVolumeSource{
		NFS: &v1.NFSVolumeSource{
			Server:   p.Server,
			Path:     fmt.Sprintf("%s/%s", p.Nfs.RootPath, subPath),
			ReadOnly: false,
		},
	}
	return pvs
}
