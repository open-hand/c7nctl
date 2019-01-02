package install

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/kube"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestCreatePv(t *testing.T) {
	client := kube.GetClient()

	pvLabels := make(map[string]string)
	pvLabels["usage"] = "test"

	pvCapacity := make(map[v1.ResourceName]resource.Quantity)
	q := resource.MustParse("10Gi")

	pvCapacity["storage"] = q

	pvs := v1.PersistentVolumeSource{
		NFS: &v1.NFSVolumeSource{
			Server:   "192.168.12.175",
			Path:     fmt.Sprintf("%s/%s", "/u01", "abc"),
			ReadOnly: false,
		},
	}
	p := Persistence{
		Client:       client,
		CommonLabels: pvLabels,
		AccessModes:  []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
		Capacity:     pvCapacity,
		Name:         "test-pv",
	}
	p.CreatePv(pvs)
}

func TestCreatePvc(t *testing.T) {
	client := kube.GetClient()

	pvLabels := make(map[string]string)
	pvLabels["usage"] = "test"

	p := Persistence{
		Client:       client,
		CommonLabels: pvLabels,
		AccessModes:  []v1.PersistentVolumeAccessMode{"ReadWriteMany"},
		Name:         "test-pv998",
		Namespace:    "default",
		Size:         "10Gi",
	}
	p.CreatePvc()
}
