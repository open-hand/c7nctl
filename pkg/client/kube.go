package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	stderrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/credentialprovider"
	"os"
)

type K8sClient struct {
	Namespace     string
	kubeInterface *kubernetes.Clientset
}

func NewK8sClient(kclient *kubernetes.Clientset, ns string) *K8sClient {
	return &K8sClient{
		kubeInterface: kclient,
		Namespace:     ns,
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

func (k *K8sClient) CreateImagePullSecret(server, username, password, secretName string) (*v1.Secret, error) {
	client := *k.kubeInterface

	imagePullSecret, err := client.CoreV1().Secrets(k.Namespace).Get(context.Background(), secretName, metav1.GetOptions{})
	if !k8serrors.IsNotFound(err) {
		log.Infof("Image pull secret %s already exists", secretName)
		return imagePullSecret, nil
	}
	auth := credentialprovider.DockerConfigJson{
		Auths: map[string]credentialprovider.DockerConfigEntry{
			server: credentialprovider.DockerConfigEntry{
				Username: username,
				Password: password,
			},
		},
		HttpHeaders: nil,
	}
	authByte, err := json.Marshal(auth)
	if err != nil {
		log.Errorf("Create image pull secret failed: %s", err)
	}
	imagePullSecret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: k.Namespace,
		},
		Immutable: nil,
		Data: map[string][]byte{
			".dockerconfigjson": authByte,
		},
		Type: v1.SecretTypeDockerConfigJson,
	}
	return client.CoreV1().Secrets(k.Namespace).Create(context.Background(), imagePullSecret, metav1.CreateOptions{})
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

func (k *K8sClient) DeletePvc(namespace, pvc string) error {
	client := *k.kubeInterface

	if err := client.CoreV1().PersistentVolumeClaims(namespace).Delete(context.Background(), pvc, metav1.DeleteOptions{}); err != nil {
		return stderrors.WithMessage(err, fmt.Sprintf("Failed to delete pvc %s in namesapce %s", pvc, namespace))
	}
	log.Infof("Successfully created pvc %s in namespace %s", pvc, namespace)
	return nil
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

func (k *K8sClient) DeleteDaemonSet(namespace, daemonSet string) error {
	client := *k.kubeInterface

	if err := client.AppsV1().DaemonSets(namespace).Delete(context.Background(), daemonSet, metav1.DeleteOptions{}); err != nil {
		return stderrors.WithMessage(err, fmt.Sprintf("Failed to delete daemonSet %s in namesapce %s", daemonSet, namespace))
	}
	log.Infof("Successfully created daemonSet %s in namespace %s", daemonSet, namespace)
	return nil
}

func (k *K8sClient) ExecCommand(podName, command string) error {

	cmd := []string{
		"sh",
		"-c",
		command,
	}

	req := k.kubeInterface.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(k.Namespace).
		SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	req.VersionedParams(option, scheme.ParameterCodec)

	config, err := GetConfig()
	if err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		log.Error(stderr.String())
		return err
	}
	log.Infof("Exec command result: %s", stdout.String())
	return nil
}

func (k *K8sClient) PatchServiceAccount(sa, ips string) {
	defaultSA, err := k.kubeInterface.CoreV1().ServiceAccounts(k.Namespace).Get(context.Background(), sa, metav1.GetOptions{})
	if err != nil {
		log.Error(err)
	}
	if defaultSA.ImagePullSecrets == nil {
		defaultSA.ImagePullSecrets = []v1.LocalObjectReference{}
	}
	// 当不重复时继续时添加
	for _, i := range defaultSA.ImagePullSecrets {
		if i.Name == ips {
			return
		}
	}
	defaultSA.ImagePullSecrets = append(defaultSA.ImagePullSecrets, v1.LocalObjectReference{
		Name: ips,
	})
	_, err = k.kubeInterface.CoreV1().ServiceAccounts(k.Namespace).Update(context.Background(), defaultSA, metav1.UpdateOptions{})
	if err != nil {
		log.Error(err)
	}
}
