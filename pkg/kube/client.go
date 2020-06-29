package kube

import (
	"github.com/vinkdong/gox/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// TODO remove
func getConfig() (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)

	config, err := clientConfig.ClientConfig()

	if err != nil {
		log.Error(err)
	}

	return config, err
}

func GetConfig() (*rest.Config, error) {
	return getConfig()
}

func GetClient() kubernetes.Interface {
	config, _ := getConfig()
	_, client, _ := getClientset(config)
	return client
}
