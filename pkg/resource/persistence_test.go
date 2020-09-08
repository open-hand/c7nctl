package resource

import (
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"testing"
)

func TestCreatePv(t *testing.T) {
	//client, _ := c7nclient.GetKubeClient("")

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
		CommonLabels: pvLabels,
		AccessModes:  []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
		Capacity:     pvCapacity,
		Name:         "test-pv",
	}
	p.createPv(pvs)
}

func TestCreatePvc(t *testing.T) {
	/*
		client, _ := c7nclient.GetKubeClient("")

		pvLabels := make(map[string]string)
		pvLabels["usage"] = "test"

		_ := Persistence{
			CommonLabels: pvLabels,
			AccessModes:  []v1.PersistentVolumeAccessMode{"ReadWriteMany"},
			GetName:         "test-pv998",
			Namespace:    "default",
			Size:         "10Gi",
		}
		p.createPvc()

	*/
}
