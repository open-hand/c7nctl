package client

import (
	"context"
	"fmt"
	stderrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	kubeInterface *kubernetes.Interface
}

func NewK8sClient(kclient *kubernetes.Interface) *K8sClient {
	return &K8sClient{
		kubeInterface: kclient,
	}
}

// 创建 kubernetes 的客户端
func GetKubeClient(kubeconfig string) (kubeClient *kubernetes.Clientset, err error) {
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homeDir(), ".kube", "config")
	}
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

func (k *K8sClient) CreateNamespace(namespace string) error {
	client := *k.kubeInterface
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	log.Infof("creating namespace %s", namespace)

	ns, err := client.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil {
		return stderrors.Wrap(err, fmt.Sprintf("Create Namespace %+v failed", namespace))
	}
	return nil
}

func (k *K8sClient) GetNamespace(namespace string) (ns *v1.Namespace, err error) {
	client := *k.kubeInterface
	return client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
}

func (k *K8sClient) SaveToCM(namespace string, cmName string, data map[string]string) (*v1.ConfigMap, error) {
	client := *k.kubeInterface

	cm, err := k.GetOrCreateCM(namespace, cmName)
	if err != nil {
		return nil, err
	}
	cm.Data = data
	return client.CoreV1().ConfigMaps(namespace).Update(context.Background(), cm, metav1.UpdateOptions{})
}

func (k *K8sClient) GetOrCreateCM(namespace, cmName string) (*v1.ConfigMap, error) {
	client := *k.kubeInterface

	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), cmName, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		log.Infof("Config map %s isn't existing, now Create it", cmName)
		return CreateCM(namespace, cmName)
	} else if err != nil {
		return nil, stderrors.WithMessage(err, fmt.Sprintf("Failed To get config maps %s from kubernetes", cmName))
	}
	return cm, nil
}

func (k *K8sClient) GetCM(namespace, cmName string) (*v1.ConfigMap, error) {
	client := *k.kubeInterface

	cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), cmName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func (k *K8sClient) reateCM(namespace string, cmName string) (*v1.ConfigMap, error) {
	client := *k.kubeInterface

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: namespace,
		},
	}

	cm, err := client.CoreV1().ConfigMaps(namespace).Create(context.Background(), cm, metav1.CreateOptions{})
	if err != nil {
		err = stderrors.WithMessage(err, fmt.Sprintf("Failed to create namesapce %s", cmName))
		return nil, err
	}
	log.Infof("Successfully created namespace %s in namespace %s", cmName, namespace)
	return cm, err
}

// Get exist pv
func (k *K8sClient) GetPv(pvName string) (pv *v1.PersistentVolume, err error) {
	client := *k.kubeInterface

	pv, err = client.CoreV1().PersistentVolumes().Get(context.Background(), pvName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = stderrors.WithMessage(err, fmt.Sprintf("PV %s isn't existing", pvName))
		}
		return nil, err
	}
	return pv, nil
}

// Get exist pvc
func GetPvc(namespace, pvcName string) (pvc *v1.PersistentVolumeClaim, err error) {
	client, err := GetKubeClient()
	if err != nil {
		return nil, stderrors.WithMessage(err, "Failed To get kubernetes client")
	}
	pvc, err = client.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), pvcName, metav1.GetOptions{})

	if err != nil {
		if k8serrors.IsNotFound(err) {
			err = stderrors.WithMessage(err, fmt.Sprintf("PVC %s isn't existing", pvcName))
		}
		return nil, err
	}
	return pvc, nil
}

func CreatePv(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	client, err := GetKubeClient()
	if err != nil {
		return nil, stderrors.WithMessage(err, "Failed To get kubernetes client")
	}
	return client.CoreV1().PersistentVolumes().Create(context.Background(), pv, metav1.CreateOptions{})
}

func CreatePvc(namespace string, pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	client, err := GetKubeClient()
	if err != nil {
		return nil, stderrors.WithMessage(err, "Failed To get kubernetes client")
	}
	return client.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(), pvc, metav1.CreateOptions{})
}

func GetClusterResource() (int64, int64) {
	client, err := GetKubeClient()
	if err != nil {
		// return nil, stderrors.WithMessage(err, "Failed To get kubernetes client")
		log.Error(err)
	}
	var sumMemory int64
	var sumCpu int64
	list, _ := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	for _, v := range list.Items {
		sumMemory += v.Status.Capacity.Memory().Value()
		sumCpu += v.Status.Capacity.Cpu().Value()
	}
	return sumMemory, sumCpu
}

func GetServerVersion() (*version.Info, error) {
	client, err := GetKubeClient()
	if err != nil {
		return nil, stderrors.WithMessage(err, "Failed To get kubernetes client")
	}
	return client.Discovery().ServerVersion()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
