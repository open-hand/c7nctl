package config

import (
	"fmt"
	"k8s.io/api/core/v1"
)

type Config struct {
	Version  string
	Metadata Metadata
	Spec     Spec
}

func (c *Config) GetStorageClassName() string {
	return c.Spec.Persistence.StorageClassName
}

func (c *Config) IgnorePv() bool {
	if c.GetStorageClassName() == "" {
		return false
	}
	// from now just support nfs
	if c.Spec.Persistence.Nfs.Server == "" {
		return true
	}
	return false
}

func (c *Config) GetResource(key string) *Resource {
	if c == nil {
		return nil
	}
	if c.Spec.Resources == nil {
		return nil
	}
	if val, ok := c.Spec.Resources[key]; ok {
		return val
	}
	return nil
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
	Nfs                     `yaml:"nfs"`
	HostPath                `yaml:"hostPath"`
	StorageClassName string `yaml:"storageClassName"`
}

type Nfs struct {
	Server   string
	RootPath string `yaml:"rootPath"`
}

type HostPath struct {
	RootPath string `yaml:"rootPath"`
} 

type Resource struct {
	Host     string
	Port     int32
	Username string
	Password string
	Schema   string
	Domain   string
	External bool
	Url      string
}

func (p *Persistence) GetPersistentVolumeSource(subPath string) v1.PersistentVolumeSource {
	if p.Nfs.Server != "" {
		return p.prepareNfsPVS(subPath)
	}
	if p.HostPath.RootPath != ""{
		return p.prepareHostPathPVS(subPath)
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

func (p *Persistence) prepareHostPathPVS(subPath string) v1.PersistentVolumeSource {
	pvs := v1.PersistentVolumeSource{
		HostPath: &v1.HostPathVolumeSource{
			Path: fmt.Sprintf("%s/%s", p.HostPath.RootPath, subPath),
		},
	}
	return pvs
}
