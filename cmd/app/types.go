package app

import "k8s.io/api/core/v1"

type InstallDefine struct {
	Version string
	Metadata  Metadata
	Spec  Spec
}

type Metadata struct {
	Name string
}

type Spec struct {
	Resources v1.ResourceRequirements
}
