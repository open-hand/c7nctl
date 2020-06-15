package action

import (
	"github.com/choerodon/c7nctl/pkg/client"
	// "github.com/choerodon/c7nctl/pkg/helm"
	"k8s.io/client-go/kubernetes"
)

// Configuration injects the dependencies that all actions shares.
type Configuration struct {

	// Release stores records of c7n component release
	release string

	// kubeClient is a kubernetes API client
	KubeClient *kubernetes.Interface

	// HelmClient is a client for working with helm
	HelmClient *client.HelmClient
}

func NewCfg() *Configuration {
	return &Configuration{
		release:    "",
		KubeClient: new(kubernetes.Interface),
		HelmClient: new(client.HelmClient),
	}
}

func (c *Configuration) InitCfg() {
	c.KubeClient = client.GetKubeClient()
	c.HelmClient = client.GetHelmClient(c.HelmClient)
}
