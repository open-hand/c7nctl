package slaver

import (
	"testing"
	"github.com/choerodon/c7n/pkg/kube"
	"k8s.io/api/core/v1"
	"github.com/vinkdong/gox/log"
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
		Ports:        []v1.ContainerPort{v1.ContainerPort{
			Name:"http",
			ContainerPort: 9800,
		}},
	}
	stopCh := make(chan struct{})
	port := slaver.ForwardPort(stopCh)
	log.Infof("success get listening port on %d",port)
	time.Sleep(time.Second * 60)
}
