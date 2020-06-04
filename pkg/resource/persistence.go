package resource

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/choerodon/c7nctl/pkg/utils"
	"github.com/vinkdong/gox/log"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Persistence struct {
	Client       kubernetes.Interface
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
}

func (p *Persistence) PrepareJobInfo() *context.JobInfo {
	ji := &context.JobInfo{
		Name:      p.Name,
		Namespace: p.Namespace,
		Type:      context.PvType,
		Status:    context.SucceedStatus,
		RefName:   p.RefPvName,
	}
	return ji
}

// Get exist pv
func (p *Persistence) getPv() (hasFound bool, pv *v1.PersistentVolume) {
	client := p.Client
	pv, err := client.CoreV1().PersistentVolumes().Get(p.RefPvName, meta_v1.GetOptions{})
	if err != nil {
		if context.IsNotFound(err) {
			return false, pv
		}
	}
	return true, pv
}

// Get exist pvc
func (p *Persistence) getPvc() (hasFound bool, pvc *v1.PersistentVolumeClaim) {
	client := p.Client
	pvc, err := client.CoreV1().PersistentVolumeClaims(p.Namespace).Get(p.RefPvcName, meta_v1.GetOptions{})
	if err != nil {
		if context.IsNotFound(err) {
			return false, pvc
		}
	}
	return true, pvc
}

// check and create pv with defined pv schema
func (p *Persistence) CheckOrCreatePv(pvs v1.PersistentVolumeSource) error {
	if p.RefPvName == "" {
		p.RefPvName = p.Name
	}
	if _, ji := context.Ctx.GetJobInfo(p.Name); ji != nil && ji.Type == context.PvType {
		log.Infof("using exist pv [%s]", ji.RefName)
		p.RefPvName = ji.RefName
		return nil
	}

	if context.Ctx.UserConfig.IgnorePv() {
		p.RefPvName = ""
		log.Debug("ignore create pv because specify storage class and no other persistence config")
		return nil
	}

	// create dir
	dir := slaver.Dir{
		Mode: p.Mode,
		Path: p.Path,
		Own:  p.Own,
	}
	if context.Ctx.Slaver == nil {
		goto checkpv
	}
	if err := context.Ctx.Slaver.MakeDir(dir); dir.Path != "" && err != nil {
		return err
	}

checkpv:
	if got, _ := p.getPv(); got {
		p.RefPvName = fmt.Sprintf("%s-%s", p.Name, utils.RandomString())
		goto checkpv
	}
	return p.CreatePv(pvs)
}

func (p *Persistence) CheckOrCreatePvc() error {
	if p.RefPvcName == "" {
		p.RefPvcName = p.Name
	}
	if _, ji := context.Ctx.GetJobInfo(p.Name); ji != nil && ji.Type == context.PvType {
		p.RefPvcName = ji.RefName
		return nil
	}
checkpvc:
	if got, _ := p.getPvc(); got {
		p.RefPvcName = fmt.Sprintf("%s-%s", p.Name, utils.RandomString())
		goto checkpvc
	}
	return p.CreatePvc()
}

func (p *Persistence) CreatePv(pvs v1.PersistentVolumeSource) error {
	log.Infof("creating pv %s", p.RefPvName)
	client := p.Client
	if len(p.AccessModes) == 0 {
		p.AccessModes = []v1.PersistentVolumeAccessMode{"ReadWriteOnce"}
	}

	if p.Capacity == nil {
		p.Capacity = make(map[v1.ResourceName]resource.Quantity)
		q := resource.MustParse(p.Size)
		p.Capacity["storage"] = q
	}

	mountOptions := p.MountOptions

	storageClassName := context.Ctx.UserConfig.GetStorageClassName()

	pv := &v1.PersistentVolume{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   p.RefPvName,
			Labels: p.CommonLabels,
		},
		Spec: v1.PersistentVolumeSpec{
			AccessModes:            p.AccessModes,
			Capacity:               p.Capacity,
			PersistentVolumeSource: pvs,
			MountOptions:           mountOptions,
			StorageClassName:       storageClassName,
		},
	}

	news := p.PrepareJobInfo()
	defer context.Ctx.AddJobInfo(news)

	_, err := client.CoreV1().PersistentVolumes().Create(pv)
	if err != nil {
		news.Status = context.FailedStatus
		news.Reason = err.Error()
		return err
	}
	log.Successf("created pv [%s]", p.RefPvName)
	return nil
}

func (p *Persistence) CreatePvc() error {
	client := p.Client

	q := resource.MustParse(p.Size)

	resList := v1.ResourceList{
		"storage": q,
	}
	res := v1.ResourceRequirements{
		Requests: resList,
	}

	storageClassName := context.Ctx.UserConfig.GetStorageClassName()

	pvc := &v1.PersistentVolumeClaim{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   p.RefPvcName,
			Labels: p.CommonLabels,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes:      p.AccessModes,
			Resources:        res,
			VolumeName:       p.RefPvName,
			StorageClassName: &storageClassName,
		},
	}

	ji := p.PrepareJobInfo()
	ji.Type = context.PvcType
	ji.RefName = p.RefPvcName

	defer context.Ctx.AddJobInfo(ji)

	_, err := client.CoreV1().PersistentVolumeClaims(p.Namespace).Create(pvc)
	if err != nil {
		log.Error(err)
		ji.Status = context.FailedStatus
		ji.Reason = err.Error()
		return err
	}
	log.Successf("created pvc [%s]", p.RefPvcName)
	return nil
}
