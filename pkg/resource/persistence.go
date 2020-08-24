package resource

import (
	"fmt"
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	c7nconsts "github.com/choerodon/c7nctl/pkg/common/consts"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Persistence struct {
	Client       *c7nclient.K8sClient
	CommonLabels map[string]string
	AccessModes  []v1.PersistentVolumeAccessMode
	Capacity     v1.ResourceList
	Name         string
	PvcEnabled   bool
	Path         string
	RootPath     string
	Size         string
	Namespace    string
	RefPvName    string
	RefPvcName   string
	Mode         string
	Own          string
	MountOptions []string
	StorageClass string
}

// check and create pv with defined pv schema
func (p *Persistence) CheckOrCreatePv(pvs v1.PersistentVolumeSource) error {
	if p.RefPvName == "" {
		p.RefPvName = p.Name
	}
	ti, err := p.Client.GetTaskInfoFromCM(p.Namespace, p.Name)
	if err != nil {
		if err.Error() == "Task info is not found" {
			ti = c7nclient.TaskInfo{
				Name:    p.Name,
				RefName: p.Name,
				Type:    c7nconsts.StaticPersistentKey,
				Status:  c7nconsts.UninitializedStatus,
			}
		} else {
			return err
		}
	}
	if ti.Status == c7nconsts.SucceedStatus && ti.Type == c7nconsts.PvType {
		log.Infof("using exist pv [%s]", ti.RefName)
		p.RefPvName = ti.RefName
		return nil
	}

	// 获得一个不重复的 pv name
	for {
		if got, _ := p.getPv(); got {
			p.RefPvName = fmt.Sprintf("%s-%s", p.Name, c7nutils.RandomString())
		} else {
			break
		}
	}
	return p.createPv(pvs)
}

func (p *Persistence) CheckOrCreatePvc(sc string) error {
	if p.RefPvcName == "" {
		p.RefPvcName = p.Name
	}
	ti, err := p.Client.GetTaskInfoFromCM(p.Namespace, p.Name)
	if err != nil {
		if err.Error() == "Task info is not found" {
			ti = c7nclient.TaskInfo{
				Name:    p.Name,
				RefName: p.Name,
				Type:    c7nconsts.StaticPersistentKey,
				Status:  c7nconsts.UninitializedStatus,
			}
		} else {
			return err
		}
	}
	if ti.Name != "" && ti.Status == c7nconsts.SucceedStatus {
		p.RefPvcName = ti.RefName
		log.Infof("using existing pvc %s", ti.RefName)
		return nil
	}
	// 获得一个不重复的 pvc name
	for {
		if got, _ := p.getPvc(); got {
			p.RefPvcName = fmt.Sprintf("%s-%s", p.Name, c7nutils.RandomString())
		} else {
			break
		}
	}
	return p.createPvc(sc)
}

func (p *Persistence) createPv(pvs v1.PersistentVolumeSource) error {
	log.Infof("creating pv %s", p.RefPvName)
	if len(p.AccessModes) == 0 {
		p.AccessModes = []v1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	}

	if p.Capacity == nil {
		p.Capacity = make(map[v1.ResourceName]resource.Quantity)
		q := resource.MustParse(p.Size)
		p.Capacity["storage"] = q
	}

	pv := &v1.PersistentVolume{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   p.RefPvName,
			Labels: p.CommonLabels,
		},
		Spec: v1.PersistentVolumeSpec{
			AccessModes:            p.AccessModes,
			Capacity:               p.Capacity,
			PersistentVolumeSource: pvs,
			MountOptions:           p.MountOptions,
			StorageClassName:       p.StorageClass,
		},
	}

	news := p.prepareTaskInfo()
	defer p.Client.SaveTaskInfoToCM(p.Namespace, *news)

	_, err := p.Client.CreatePv(pv)
	if err != nil {
		news.Status = c7nconsts.FailedStatus
		news.Reason = err.Error()
		return err
	}
	log.Info("created pv [%s]", p.RefPvName)
	return nil
}

func (p *Persistence) createPvc(sc string) error {
	q := resource.MustParse(p.Size)

	pvc := &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   p.RefPvcName,
			Labels: p.CommonLabels,
			// 基于 NFS storageClass 的 PVC 自动创建
			Annotations: map[string]string{
				"volume.beta.kubernetes.io/storage-class": sc,
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: p.AccessModes,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": q,
				},
			},
			VolumeName:       p.RefPvName,
			StorageClassName: &sc,
		},
	}

	ti := p.prepareTaskInfo()
	ti.RefName = p.RefPvcName
	defer p.Client.SaveTaskInfoToCM(p.Namespace, *ti)

	_, err := p.Client.CreatePvc(p.Namespace, pvc)
	if err != nil {
		log.Error(err)
		ti.Status = c7nconsts.FailedStatus
		ti.Reason = err.Error()
		return err
	}
	log.Infof("created pvc [%s]", p.RefPvcName)
	return nil
}

func (p *Persistence) prepareTaskInfo() *c7nclient.TaskInfo {
	ti := &c7nclient.TaskInfo{
		Name:      p.Name,
		Namespace: p.Namespace,
		Type:      c7nconsts.StaticPersistentKey,
		Status:    c7nconsts.SucceedStatus,
		RefName:   p.RefPvName,
	}
	return ti
}

// Get exist pv
func (p *Persistence) getPv() (hasFound bool, pv *v1.PersistentVolume) {
	pv, err := p.Client.GetPv(p.RefPvName)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, pv
		}
	}
	return true, pv
}

// Get exist pvc
func (p *Persistence) getPvc() (hasFound bool, pvc *v1.PersistentVolumeClaim) {
	pvc, err := p.Client.GetPvc(p.Namespace, p.RefPvcName)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, pvc
		}
	}
	return true, pvc
}
