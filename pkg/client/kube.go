package client

import (
	"context"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/consts"
	stderrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	kubeInterface *kubernetes.Clientset
}

func NewK8sClient(kclient *kubernetes.Clientset) *K8sClient {
	return &K8sClient{
		kubeInterface: kclient,
	}
}

// 创建 kubernetes 的客户端
func GetKubeClient(kubeconfig string) (kubeClient *kubernetes.Clientset, err error) {
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

func GetConfig() (*rest.Config, error) {
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

func (k *K8sClient) GetClientSet() *kubernetes.Clientset {
	return k.kubeInterface
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
		return k.CreateCM(namespace, cmName)
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

func (k *K8sClient) CreateCM(namespace string, cmName string) (*v1.ConfigMap, error) {
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
		return nil, err
	}
	return pv, nil
}

// Get exist pvc
func (k *K8sClient) GetPvc(namespace, pvcName string) (pvc *v1.PersistentVolumeClaim, err error) {
	client := *k.kubeInterface

	pvc, err = client.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), pvcName, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}
	return pvc, nil
}

func (k *K8sClient) CreatePv(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	client := *k.kubeInterface

	return client.CoreV1().PersistentVolumes().Create(context.Background(), pv, metav1.CreateOptions{})
}

func (k *K8sClient) CreatePvc(namespace string, pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	client := *k.kubeInterface

	return client.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(), pvc, metav1.CreateOptions{})
}

func (k *K8sClient) GetClusterResource() (int64, int64) {
	client := *k.kubeInterface

	var sumMemory int64
	var sumCpu int64
	list, _ := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	for _, v := range list.Items {
		sumMemory += v.Status.Capacity.Memory().Value()
		sumCpu += v.Status.Capacity.Cpu().Value()
	}
	return sumMemory, sumCpu
}

func (k *K8sClient) GetServerVersion() (*version.Info, error) {
	client := *k.kubeInterface

	return client.Discovery().ServerVersion()
}

func (k *K8sClient) GetTaskInfoFromCM(namespace, taskName string) (TaskInfo, error) {
	logs, err := k.GetOrCreateCM(namespace, consts.StaticLogsCM)
	if err != nil {
		return TaskInfo{}, err
	}
	keys := []string{consts.StaticReleaseKey, consts.StaticTaskKey, consts.StaticPersistentKey}
	for _, key := range keys {
		var tasks []TaskInfo
		if err := yaml.Unmarshal([]byte(logs.Data[key]), &tasks); err != nil {
			return TaskInfo{}, err
		}
		for _, t := range tasks {
			if t.Name == taskName {
				return t, nil
			}
		}
	}
	return TaskInfo{}, stderrors.New("Task info is not found")
}

func (k *K8sClient) SaveTaskInfoToCM(namespace string, task TaskInfo) error {
	c7nLogs, err := k.GetOrCreateCM(namespace, consts.StaticLogsCM)
	if err != nil {
		return err
	}

	var tasks []TaskInfo
	if err := yaml.Unmarshal([]byte(c7nLogs.Data[task.Type]), &tasks); err != nil {
		return err
	}
	// 如果存在就替换，不存在则添加
	var isExisting bool
	for idx, t := range tasks {
		if t.Name == task.Name {
			tasks[idx] = task
			isExisting = true
			break
		}
	}
	if !isExisting {
		tasks = append(tasks, task)
	}
	tasksStr, err := yaml.Marshal(tasks)
	if err != nil {
		return err
	}
	if c7nLogs.Data == nil {
		c7nLogs.Data = map[string]string{}
	}
	c7nLogs.Data[task.Type] = string(tasksStr)

	_, err = k.SaveToCM(namespace, consts.StaticLogsCM, c7nLogs.Data)
	if err != nil {
		log.Error(err)
	}
	return nil
}
