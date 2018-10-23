package slaver

import (
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/kube"
	pb "github.com/choerodon/c7n/pkg/protobuf"
	"github.com/vinkdong/gox/log"
	"k8s.io/api/core/v1"
	"testing"
	"time"
)

const C7nLabelKey = "c7n-usage"
const C7nLabelValue = "c7n-installer"

func TestInstall(t *testing.T) {

	port := v1.ContainerPort{
		Name:          "http",
		ContainerPort: 9000,
	}
	ports := []v1.ContainerPort{port}
	labels := make(map[string]string)
	labels[C7nLabelKey] = C7nLabelValue
	slaver := Slaver{
		Client:       kube.GetClient(),
		Version:      "0.1",
		Namespace:    "install",
		Name:         "c7n-slaver",
		CommonLabels: labels,
		Image:        "vinkdong/timing",
		Ports:        ports,
	}
	ds, err := slaver.Install()
	if err != nil {
		t.Error(err)
		t.Error("install daemonset failed")
	} else {
		log.Info(ds.Spec)
	}
}

func TestPortForward(t *testing.T) {
	labels := make(map[string]string)
	labels[C7nLabelKey] = C7nLabelValue
	slaver := Slaver{
		Client:       kube.GetClient(),
		Version:      "0.1",
		Namespace:    "install",
		Name:         "c7n-slaver",
		CommonLabels: labels,
		Image:        "vinkdong/timing",
		Ports: []v1.ContainerPort{v1.ContainerPort{
			Name:          "http",
			ContainerPort: 9800,
		}},
	}
	stopCh := make(chan struct{})
	port := slaver.ForwardPort("http", stopCh)
	log.Infof("success get listening port on %d", port)
	time.Sleep(time.Second * 60)
}

func TestCheckHealth(t *testing.T) {
	slaver := Slaver{
		GRpcAddress: "127.0.0.1:9001",
	}
	check := &pb.Check{
		Type:   "httpGet",
		Host:   "baidu.com",
		Port:   443,
		Schema: "https",
		Path:   "/",
	}
	log.Info(slaver.CheckHealth("test-service", check))
	check = &pb.Check{
		Type: "socket",
		Host: "baidu.com",
		Port: 445,
	}
	log.Info(slaver.CheckHealth("test-service", check))
}

func TestExecuteRemoteSql(t *testing.T) {
	slaver := Slaver{
		GRpcAddress: "127.0.0.1:9001",
	}
	sqlList := []string{
		"CREATE DATABASE abc",
		"DROP DATABASE abc",
	}
	r := &config.Resource{
		Host:     "192.168.99.100",
		Port:     3306,
		Username: "root",
		Password: "abc123",
	}
	log.Info(slaver.ExecuteRemoteSql(sqlList, r))
}

func TestExecuteRemoteCommand(t *testing.T) {
	slaver := Slaver{
		GRpcAddress: "127.0.0.1:9001",
	}
	cmdList := []string{
		"ls",
		"`mkdir -p abc/123`",
		"pwd",
		"`chown -R 1001:1001 abc`",
	}
	log.EnableDebug()
	log.Info(slaver.ExecuteRemoteCommand(cmdList))
}
