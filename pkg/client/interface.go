package client

import "github.com/choerodon/c7nctl/pkg/config"

// Interface represents a client capable of
type Interface interface {
	GetConfig() *config.Config

	GetHelmClient() *HelmClient
}
