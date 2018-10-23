package slaver

import (
	"bytes"
	"context"
	"encoding/json"
	sys_errors "errors"
	"fmt"
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/kube"
	pb "github.com/choerodon/c7n/pkg/protobuf"
	"github.com/vinkdong/gox/log"
	"google.golang.org/grpc"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
	VolumeMounts []core_v1.VolumeMount
	PodList      *core_v1.PodList
	Address      string
	GRpcAddress  string
	PvcName      string
	DataPath     string
}

const IngressCheckPath = "/c7n/acme-challenge"

type Dir struct {
	Mode string
	Own  string
	Path string
}

/**
Type: httpGet or socket
*/

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
		Name:            s.Name,
		Image:           s.Image,
		Ports:           s.Ports,
		Env:             s.Env,
		VolumeMounts:    s.VolumeMounts,
		ImagePullPolicy: "Always",
	}

	volume := core_v1.Volume{
		Name: "data",
		VolumeSource: core_v1.VolumeSource{
			PersistentVolumeClaim: &core_v1.PersistentVolumeClaimVolumeSource{
				ClaimName: s.PvcName,
			},
		},
	}

	tmp := core_v1.PodTemplateSpec{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels: s.CommonLabels,
		},
		Spec: core_v1.PodSpec{
			Containers: []core_v1.Container{dsContainer},
			Volumes:    []core_v1.Volume{volume},
		},
	}

	selector := &meta_v1.LabelSelector{
		MatchLabels: s.CommonLabels,
	}
	s.CommonLabels["app"] = s.Name
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
	opts := meta_v1.ListOptions{
		LabelSelector: set.AsSelector().String(),
	}
	return s.Client.CoreV1().Pods(s.Namespace).List(opts)
}

func (s *Slaver) CheckRunning() bool {
	log.Info("waiting slaver running...")
	poList, err := s.GetPods()
	if err != nil || len(poList.Items) < 1 {
		return false
	}
	for _, po := range poList.Items {
		if po.Status.Phase != core_v1.PodRunning {
			time.Sleep(time.Second * 5)
			return false
		}
	}
	s.PodList = poList
	return true
}

func (s *Slaver) getForwardPorts(portName string, localPort int) string {
	for _, port := range s.Ports {
		if port.Name == portName {
			return fmt.Sprintf("%d:%d", localPort, port.ContainerPort)
		}
	}
	log.Errorf("no slave %s port found", portName)
	os.Exit(129)
	return ""
}

func (s *Slaver) ForwardPort(portName string, stopCh <-chan struct{}) int {

	rest := s.Client.CoreV1().RESTClient()

	var pod core_v1.Pod

loop:
	for {
		select {
		case <-time.Tick(time.Second):
			if s.CheckRunning() {
				break loop
			}
		}
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

	fw, err := portforward.New(dialer, []string{s.getForwardPorts(portName, port)}, stopCh, readyCh, os.Stdout, os.Stderr)

	if err != nil {
		log.Error(err)
	}
	go fw.ForwardPorts()
	<-readyCh

	return port
}

func (s *Slaver) MakeDir(dir Dir) error {
	log.Infof("create dir %s with mode %s own %s ", dir.Path, dir.Mode, dir.Own)

	if len(s.VolumeMounts) < 1 {
		err := sys_errors.New("slaver have not mount any volumes")
		return err
	}
	rootPath := s.VolumeMounts[0].MountPath

	cmdList := []string{
		fmt.Sprintf("`mkdir -p %s/%s -m %s`", rootPath, dir.Path, dir.Mode),
	}
	if dir.Own != "" {
		cmdList = append(cmdList, fmt.Sprintf("`chown -R %s %s/%s`", dir.Own, rootPath, dir.Path))
	}

	if created := s.ExecuteRemoteCommand(cmdList); created != false {
		return nil
	}

	return sys_errors.New(fmt.Sprintf("can't create dir %s with mode %s", dir.Path, dir.Mode))
}

func (s *Slaver) connectGRpc() (*grpc.ClientConn, error) {

	return grpc.Dial(s.GRpcAddress, grpc.WithInsecure())
}

func (s *Slaver) CheckHealth(name string, check *pb.Check) bool {
	conn, err := s.connectGRpc()
	if err != nil {
		log.Errorf("connect %s grpc path  failed", s.GRpcAddress)
		return false
	}
	c := pb.NewRouteCallClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1)
remoteCheck:
	r, err := c.CheckHealth(ctx, check)
	if err != nil {
		log.Debugf("check %s health failed with msg: '%s' retry ..", name, err.Error())
		time.Sleep(time.Second * 10)
		goto remoteCheck
	}
	defer cancel()

	if r.Success == false {
		log.Debugf("check health failed with msg: %s retry..", r.Message)
		time.Sleep(time.Second * 10)
		goto remoteCheck
	}
	return true
}

type Request struct {
	Url    string
	Method string
	Body   string
}

type Forward struct {
	Url    string              `json:"url"`
	Body   string              `json:"body"`
	Method string              `json:"method"`
	Header map[string][]string `json:"header"`
}

func (s *Slaver) ExecuteRemoteRequest(f Forward) error {
	url := fmt.Sprint(s.Address, "/forward")

	data, err := json.Marshal(f)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header = f.Header
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		log.Errorf("request %s \nheader: %v \n with body \n %s\nfailed", f.Url, f.Header, f.Body)
		return sys_errors.New(fmt.Sprintf("resp code %d not is 2xx or 3xx", resp.StatusCode))
	}
	return nil
}

func (s *Slaver) ExecuteRemoteSql(sqlList []string, resource *config.Resource) bool {
	conn, err := s.connectGRpc()
	if err != nil {
		log.Errorf("connect %s grpc path  failed", s.GRpcAddress)
		return false
	}
	c := pb.NewRouteCallClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1)
	defer cancel()
	stream, err := c.ExecuteSql(ctx)
	if err != nil {
		log.Error(err)
		return false
	}
	m := &pb.Mysql{
		Host:     resource.Host,
		Port:     resource.Port,
		Username: resource.Username,
		Password: resource.Password,
	}
	r := &pb.RouteSql{
		Mysql: m,
	}
	err = stream.Send(r)
	if err != nil {
		log.Error(err)
	}

	for _, sql := range sqlList {
		r := &pb.RouteSql{
			Sql: sql,
		}
		log.Debugf("executing: %s", sql)
		stream.Send(r)
		resp, err := stream.Recv()
		if err != nil {
			log.Error(err)
			return false
		}
		if !resp.Success {
			log.Error(resp.Message)
			return false
		}
	}

	return true
}

func (s *Slaver) ExecuteRemoteCommand(commands []string) bool {
	conn, err := s.connectGRpc()
	if err != nil {
		log.Errorf("connect %s grpc path  failed", s.GRpcAddress)
		return false
	}
	c := pb.NewRouteCallClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*1)
	defer cancel()
	stream, err := c.ExecuteCommand(ctx)

	for _, c := range commands {
		routeCommand := &pb.RouteCommand{
			Name: "sh",
			Args: []string{"-c", c},
		}
		log.Debugf("executed %s %s", routeCommand.Name, strings.Join(routeCommand.Args, " "))
		if err := stream.Send(routeCommand); err != nil {
			log.Error(err)
			return false
		}
		result, err := stream.Recv()
		if err != nil {
			log.Error(err)
			return false
		}
		if !result.Success {
			log.Error(result.Message)
			return false
		}
		log.Debugf(result.Message)
	}
	return true
}

func (s *Slaver) InstallService() (*core_v1.Service, error) {
	svcInterface := s.Client.CoreV1().Services(s.Namespace)

	svc, err := svcInterface.Get(s.Name, meta_v1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if svc != nil && err == nil {
		return svc, err
	}
	port := intstr.IntOrString{
		Type:   1,
		StrVal: "http",
	}
	servicePort := core_v1.ServicePort{
		Port:       80,
		Protocol:   "TCP",
		TargetPort: port,
	}

	service := &core_v1.Service{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   s.Name,
			Labels: s.CommonLabels,
		},
		Spec: core_v1.ServiceSpec{
			Ports:    []core_v1.ServicePort{servicePort},
			Selector: s.CommonLabels,
		},
	}

	return svcInterface.Create(service)
}

func (s *Slaver) UpdateIngress(ingress *v1beta1.Ingress, domain string) error {
	for _, r := range ingress.Spec.Rules {
		if r.Host == domain {
			return nil
		}
	}
	ruleList := ingress.Spec.Rules
	ingressRule, err := s.getIngressRule(domain)
	if err != nil {
		return err
	}
	ingress.Spec.Rules = append(ruleList, ingressRule)

	ingressInterface := s.Client.ExtensionsV1beta1().Ingresses(s.Namespace)

	_, err = ingressInterface.Update(ingress)
	return err
}

func (s *Slaver) getIngressRule(domain string) (v1beta1.IngressRule, error) {
	port := intstr.IntOrString{
		Type:   1,
		StrVal: "http",
	}
	svc, err := s.InstallService()

	if err != nil {
		return v1beta1.IngressRule{}, err
	}

	backend := v1beta1.IngressBackend{
		ServiceName: svc.Name,
		ServicePort: port,
	}

	ingressPath := v1beta1.HTTPIngressPath{
		Path:    IngressCheckPath,
		Backend: backend,
	}
	ingressRuleValue := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{ingressPath},
		},
	}
	ingressRule := v1beta1.IngressRule{
		Host:             domain,
		IngressRuleValue: ingressRuleValue,
	}

	return ingressRule, nil
}

func (s *Slaver) InstallIngress(domain string) error {

	ingressInterface := s.Client.ExtensionsV1beta1().Ingresses(s.Namespace)

	ing, err := ingressInterface.Get(s.Name+"checker", meta_v1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	if err == nil {
		return s.UpdateIngress(ing, domain)
	}

	ingressRule, err := s.getIngressRule(domain)

	if err != nil {
		return err
	}
	ingress := &v1beta1.Ingress{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   s.Name + "checker",
			Labels: s.CommonLabels,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{ingressRule},
		},
	}
	_, err = ingressInterface.Create(ingress)
	return err
}

func (s *Slaver) CheckClusterDomain(domain string) error {
	err := s.InstallIngress(domain)
	if err != nil {
		return err
	}
	log.Infof("send msg to check domain %s", domain)
	return nil
}

func (s *Slaver) Uninstall() error {
	return s.Client.AppsV1().DaemonSets(s.Namespace).Delete(s.Name, &meta_v1.DeleteOptions{})
}
