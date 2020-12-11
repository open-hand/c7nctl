package action

import (
	"context"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/cli"
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/resource"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/api/errors"
	"time"
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

func (c *C7nConfiguration) Init(s *cli.EnvSettings) {
	cfg := c7nclient.InitConfiguration(s.KubeConfig, s.Namespace)
	// 初始化 helm3Client
	c.HelmClient = c7nclient.NewHelm3Client(cfg)
	// 初始化 kubeClient
	kubeclient, _ := c7nclient.GetKubeClient(s.KubeConfig)
	c.KubeClient = c7nclient.NewK8sClient(kubeclient, s.Namespace)
}

// 基础组件——比如 gitlab-ha ——有 app 标签，c7n 有 choerodon.io/release 标签
// TODO 去掉 app label
func (c *C7nConfiguration) CheckReleasePodRunning(rls, namespace string) {
	clientset := c.KubeClient.GetClientSet()

	labels := []string{
		fmt.Sprintf("choerodon.io/release=%s", rls),
		fmt.Sprintf("app=%s", rls),
	}

	log.Infof("Waiting %s running", rls)
	for {
		for _, label := range labels {
			deploy, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), meta_v1.ListOptions{LabelSelector: label})
			if k8serrors.IsNotFound(err) {
				log.Debugf("Deployment %s in namespace %s not found\n", label, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				log.Debugf("Error getting deployment %s in namespace %s: %v\n",
					label, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				for _, d := range deploy.Items {
					if *d.Spec.Replicas != d.Status.ReadyReplicas {
						log.Debugf("Release %s is not ready\n", d.Name)
					} else {
						log.Debugf("Release %s is Ready\n", d.Name)
						return
					}
				}
			}
			ss, err := clientset.AppsV1().StatefulSets(namespace).List(context.Background(), meta_v1.ListOptions{LabelSelector: label})
			if k8serrors.IsNotFound(err) {
				log.Debugf("StatefulSet %s in namespace %s not found\n", label, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				log.Debugf("Error getting statefulSet %s in namespace %s: %v\n",
					label, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				for _, s := range ss.Items {
					if *s.Spec.Replicas != s.Status.ReadyReplicas {
						log.Debugf("Release %s is not ready\n", s.Name)
					} else {
						log.Debugf("Release %s is Ready\n", s.Name)
						return
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (i *Install) CheckNamespace() error {
	_, err := i.cfg.KubeClient.GetNamespace(i.Namespace)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return i.cfg.KubeClient.CreateNamespace(i.Namespace)
		}
		return err
	}
	log.Infof("namespace %s already exists", i.Namespace)
	return nil
}

func (c *C7nConfiguration) CreateImagePullSecret(drs []resource.DockerRegistry) {
	if drs == nil {
		log.Debug("Skip up create image pull secret.")
		return
	}
	for _, ds := range drs {
		if _, err := c.KubeClient.CreateImagePullSecret(ds.Server, ds.Username, ds.Password, ds.SecretName); err != nil {
			log.Error("Create image pull secret %s failed: %s", ds.SecretName, err)
			continue
		}
		c.KubeClient.PatchServiceAccount(ds.ServiceAccount, ds.SecretName)
	}
}

func (c *C7nConfiguration) CheckResource(resources *v1.ResourceRequirements) error {
	request := resources.Requests

	reqMemory := request.Memory().Value()
	reqCpu := request.Cpu().Value()
	clusterMemory, clusterCpu := c.KubeClient.GetClusterResource()

	/*
		metrics.Memory = clusterMemory
		metrics.CPU = clusterCpu

		serverVersion, err := i.cfg.KubeClient.GetServerVersion()
		if err != nil {
			return std_errors.Wrap(err, "can't get your cluster version")
		}
		metrics.Version = serverVersion.String()
	*/
	if clusterMemory < reqMemory {
		return std_errors.New(fmt.Sprintf("cluster memory not enough, request %dGi", reqMemory/(1024*1024*1024)))
	}
	if clusterCpu < reqCpu {
		return std_errors.New(fmt.Sprintf("cluster cpu not enough, request %dc", reqCpu/1000))
	}
	return nil
}
