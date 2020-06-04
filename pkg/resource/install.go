package resource

import (
	"fmt"
	c7n_config "github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/util/maps"
)

type InstallDefinition struct {
	// api 版本
	Version string
	// Choerodon 平台版本
	PaaSVersion string
	Metadata    Metadata
	Spec        Spec
	// TODO REMOVE
	Resource           *c7n_config.Resource
	CommonLabels       map[string]string
	DefaultAccessModes []v1.PersistentVolumeAccessMode `yaml:"accessModes"`
	SkipInput          bool
	Timeout            int
	Prefix             string
	Namespace          string
	Mail               string
}

type Metadata struct {
	Name      string
	Namespace string
}

type Spec struct {
	Basic     Basic
	Resources v1.ResourceRequirements
	Release   []*Release
	Runner    *Release `json:"runner"`
	Component []*Release
}

type Basic struct {
	RepoURL string
	Slaver  slaver.Slaver
}

func (i *InstallDefinition) PrepareSlaverPvc() (string, error) {
	if context.Ctx.UserConfig == nil {
		return "", nil
	}
	pvs := context.Ctx.UserConfig.Spec.Persistence.GetPersistentVolumeSource("")

	persistence := Persistence{
		Client:       *context.Ctx.KubeClient,
		CommonLabels: i.CommonLabels,
		AccessModes:  i.DefaultAccessModes,
		Size:         "1Gi",
		Mode:         "755",
		PvcEnabled:   true,
		Name:         "slaver",
	}
	err := persistence.CheckOrCreatePv(pvs)
	if err != nil {
		return "", err
	}

	persistence.Namespace = context.Ctx.UserConfig.Metadata.Namespace

	if err := persistence.CheckOrCreatePvc(); err != nil {
		return "", err
	}
	return persistence.RefPvcName, nil
}

func (i *InstallDefinition) PrepareSlaver(stopCh <-chan struct{}) (*slaver.Slaver, error) {
	// prepare slaver to execute sql or make directory ..

	s := &i.Spec.Basic.Slaver
	s.Client = *context.Ctx.KubeClient
	// be care of use point
	s.CommonLabels = maps.CopySS(context.Ctx.CommonLabels)
	s.Namespace = context.Ctx.Namespace

	if pvcName, err := i.PrepareSlaverPvc(); err != nil {
		return s, err
	} else {
		s.PvcName = pvcName
	}

	if _, err := s.CheckInstall(); err != nil {
		return s, err
	}
	port := s.ForwardPort("http", stopCh)
	grpcPort := s.ForwardPort("grpc", stopCh)
	s.Address = fmt.Sprintf("http://127.0.0.1:%d", port)
	s.GRpcAddress = fmt.Sprintf("127.0.0.1:%d", grpcPort)
	return s, nil
}
