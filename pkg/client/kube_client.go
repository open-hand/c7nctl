package client

import (
	stderrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
)

type kubeInterface struct {
	Namespace string
	ClientSet *kubernetes.Clientset
}

var kubeClient *kubeInterface
var once sync.Once

// GetKubeInterface 单例模式初始化客户端
func GetKubeInterface() *kubeInterface {
	once.Do(func() {
		var err error
		kubeClient = new(kubeInterface)
		kubeClient.ClientSet, err = getKubeClient("")
		if err != nil {
			log.Error(err)
		}
	})

	return kubeClient
}

func getKubeClient(kubeconfig string) (kubeClient *kubernetes.Clientset, err error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		return nil, stderrors.WithMessage(err, "Get kubeconfig failed")
	}

	kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		err = stderrors.WithMessage(err, "init kubernetes client error")
	}

	return kubeClient, err
}

func (k *kubeInterface) CreateImagePullSecret(server, username, password, secretName string) (*v1.Secret, error) {
	return nil, nil
}

func (k *kubeInterface) PatchServiceAccount(sa, ips string) error {
	return nil
}

func (k *kubeInterface) CreateNamespace(namespace string) (*v1.Namespace, error) {
	return nil, nil
}

func (k *kubeInterface) CheckNamespace(namespace string) bool {
	return false
}

func (k *kubeInterface) CreatePersistentVolume(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	return nil, nil
}

func (k *kubeInterface) CreatePersistentClaim(pv *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return nil, nil
}

func (k *kubeInterface) GetClusterInfo() {

}

func (k *kubeInterface) ExecCommand(podName, cmds []string) error {
	return nil
}

func (k *kubeInterface) ExecSQL(podName, sqls []string) error {
	return nil
}

func (k *kubeInterface) ExecRequest(podName, reqs []string) error {
	return nil
}

func (k *kubeInterface) CreateDefaultIndex() error {
	return nil
}

func (k *kubeInterface) UpgradeIngress() error {
	return nil
}

func (k *kubeInterface) InstallIngress() error {
	return nil
}

func (k *kubeInterface) CheckClusterDomain(domain string) error {
	return nil
}

func createDeployment() {

}

func createService(svc *v1.Service) {

}
