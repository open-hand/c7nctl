package client

import (
	"github.com/vinkdong/gox/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeClient() *kubernetes.Interface {
	config, _ := getConfig()
	client, _ := getClientset(config)
	return &client
}

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

func getClientset(c *rest.Config) (kubernetes.Interface, error) {
	client, err := kubernetes.NewForConfig(c)
	if err != nil {
		log.Error(err)
	}

	return client, err
}
