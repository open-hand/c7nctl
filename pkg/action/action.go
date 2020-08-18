package action

import (
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
)

// C7nConfiguration injects the dependencies that all actions shares.
type C7nConfiguration struct {

	// TODO c7n api client

	// kubeClient is a kubernetes API client
	// TODO refactor kubeClient
	KubeClient *c7nclient.K8sClient

	// HelmInstall is a client for working with helm
	// helm3 的都是依赖于这个
	//
	HelmClient *c7nclient.Helm3Client
}
