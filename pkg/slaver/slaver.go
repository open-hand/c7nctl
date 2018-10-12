package slaver

import (
	"k8s.io/client-go/kubernetes"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"github.com/vinkdong/gox/log"
	"k8s.io/apimachinery/pkg/labels"
	"fmt"
	"k8s.io/client-go/transport/spdy"
	"github.com/choerodon/c7n/pkg/kube"
	"net/http"
	"os"
	"k8s.io/client-go/tools/portforward"
	"net"
	"strconv"
	"time"
	"k8s.io/api/extensions/v1beta1"
)

type Slaver struct {
	Client       kubernetes.Interface
	Version      string
	Namespace    string
	Name         string
	CommonLabels map[string]string
	Image        string
	Ports        []core_v1.ContainerPort
	Env          []core_v1.EnvVar
	volumeMounts []core_v1.VolumeMount
	PodList      *core_v1.PodList
}

type Dir struct {
	Mode string
	Path string
}

func (s *Slaver) CheckInstall() (*v1beta1.DaemonSet, error) {
	ds, err := s.Client.ExtensionsV1beta1().DaemonSets(s.Namespace).Get(s.Name, meta_v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Infof("deploying daemonSet %s", s.Name)
			return s.Install()
		}
		return nil, err
	}
	return ds, err
}

func (s *Slaver) Install() (*v1beta1.DaemonSet, error) {

	dsContainer := core_v1.Container{
		Name:         s.Name,
		Image:        s.Image,
		Ports:        s.Ports,
		Env:          s.Env,
		VolumeMounts: s.volumeMounts,
	}

	tmp := core_v1.PodTemplateSpec{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels: s.CommonLabels,
		},
		Spec: core_v1.PodSpec{
			Containers: []core_v1.Container{dsContainer},
		},
	}

	selector := &meta_v1.LabelSelector{
		MatchLabels: s.CommonLabels,
	}
	ds := &v1beta1.DaemonSet{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "v1beta2",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   s.Name,
			Labels: s.CommonLabels,
		},
		Spec: v1beta1.DaemonSetSpec{
			Template: tmp,
			Selector: selector,
		},
	}
	daemonSet, err := s.Client.ExtensionsV1beta1().DaemonSets(s.Namespace).Create(ds)

	if err != nil {
		return nil, err
	}
	return daemonSet, err
}

func (s *Slaver) GetPods() (*core_v1.PodList, error) {
	set := labels.Set(s.CommonLabels)
	fmt.Println(set.AsSelector().String())
	opts := meta_v1.ListOptions{
		LabelSelector: set.AsSelector().String(),
	}
	return s.Client.CoreV1().Pods(s.Namespace).List(opts)
}

func (s *Slaver) CheckRunning() bool {
	poList, err := s.GetPods()
	if err != nil || poList.Size() < 1 {
		log.Error(err)
		return false
	}
	for _, po := range poList.Items {
		if po.Status.Phase != core_v1.PodRunning {
			return false
		}
	}
	s.PodList = poList
	return true
}

func (s *Slaver) getForwardPorts(localPort int) string {
	for _, port := range s.Ports {
		if port.Name == "http" {
			return fmt.Sprintf("%d:%d", localPort, port.ContainerPort)
		}
	}
	log.Error("no slave http port found")
	os.Exit(129)
	return ""
}

func (s *Slaver) ForwardPort(stopCh <-chan struct{}) int {

	rest := s.Client.CoreV1().RESTClient()

	var pod core_v1.Pod

	if !s.CheckRunning() {
		return 0
	}
	pod = s.PodList.Items[0]

	req := rest.Post().Resource("pods").
		Namespace(pod.Namespace).
		Name(pod.Name).
		SubResource("portforward")

	config, err := kube.GetConfig()
	if err != nil {
		log.Error(err)
	}

	transport, upgrader, err := spdy.RoundTripperFor(config)
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())

	readyCh := make(chan struct{})

	port := 8000
getFreePort:
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("", strconv.Itoa(port)), time.Second)
	if conn != nil {
		port += 1
		goto getFreePort
		conn.Close()
	}
	log.Info(port)

	fw, err := portforward.New(dialer, []string{s.getForwardPorts(port)}, stopCh, readyCh, os.Stdout, os.Stderr)

	if err != nil {
		log.Error(err)
	}
	go fw.ForwardPorts()
	<-readyCh
	return port
}

func (s *Slaver) MakeDir(dir Dir) error {
	log.Infof("create dir %s with mode %s", dir.Path, dir.Mode)
	return nil
}

func (s *Slaver) ExecuteSql(sql string) error {
	log.Infof("executed sql %s", sql)
	return nil
}

func (s *Slaver) Uninstall() error {
	return s.Client.AppsV1().DaemonSets(s.Namespace).Delete(s.Name, &meta_v1.DeleteOptions{})
}
