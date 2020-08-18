package upgrade

import (
	"context"
	"fmt"
	"github.com/buger/jsonparser"

	"github.com/choerodon/c7nctl/pkg/resource"
	"github.com/choerodon/c7nctl/pkg/utils"
	"github.com/vinkdong/gox/log"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/batch/v1"
	"strings"
	"time"
)

type Upgrader struct {
	// HelmClient *c7n_helm.HelmClient
	Version  string
	Metadata Metadata
	Spec     Spec
}

type Upgrade struct {
	Name              string
	Chart             string
	Version           string
	InstalledVersion  string
	Namespace         string
	ConstraintVersion string
	Values            []byte
	SetKey            []*SetKey
	ChangeKey         []*ChangeKey
	DeleteKey         []string
}

type Metadata struct {
	Name string
}

type Spec struct {
	Basic     Basic
	Install   []*resource.Release
	Uninstall []*Uninstall
	Upgrade   []*Upgrade
}

type Uninstall struct {
	Kind string
	Name string
}

type Basic struct {
	RepoURL string
}

type SetKey struct {
	Name  string
	Value string
	Input utils.Input
}

type ChangeKey struct {
	Old string
	New string
}

func (u *Upgrader) Init() {
	/*	helmClient := &c7n_helm.Client{
			TillerTunnel: kube.GetTunnel(),
			KubeClient:   kube.GetClient(),
		}
		helmClient.InitClient()
		u.HelmInstall = helmClient*/
}

func (u *Upgrader) GetReleaseValues(upgrade *Upgrade) error {
	/*	ls, err := u.HelmInstall.HelmInstall.ReleaseContent(upgrade.Name)
		if err != nil {
			return err
		}
		config := ls.GetRelease().GetConfig()
		upgrade.InstalledVersion = ls.GetRelease().GetChart().GetMetadata().GetVersion()
		upgrade.Namespace = ls.GetRelease().GetNamespace()
		log.Debugf("Get raw values:\n%s", config.GetRaw())
		bytes, err := yaml.YAMLToJSON([]byte(config.GetRaw()))
		if err != nil {
			return err
		}
		upgrade.Values = bytes*/
	return nil
}

func upgradeRelease(u *Upgrader, upgrade *Upgrade) error {
	if len(upgrade.Values) != 0 {
		/*
			// 解析变量
			e := upgradeValues(upgrade)
			if e != nil {
				log.Error(e)
				return e
			}
			raw, err := yaml.JSONToYAML(upgrade.Values)
			log.Debugf("After rendering values:\n%s", string(raw))
			if err != nil {
				return err
			}
			chartArgs := c7n_helm.ChartArgs{
				ReleaseName: upgrade.Name,
				RepoUrl:     u.Spec.Basic.RepoURL,
				Verify:      false,
				Version:     upgrade.Version,
				ChartName:   upgrade.Chart,
			}
			log.Infof("Upgrade %s to %s version,please waiting.", upgrade.Name, upgrade.Version)
			return u.HelmInstall.UpgradeRelease(
				raw,
				chartArgs,
			)*/
	}
	return nil
}

func (u *Upgrader) UpgradeReleases() error {
	for _, upgrade := range u.Spec.Upgrade {
		if err := upgradeRelease(u, upgrade); err != nil {
			return err
		}
	}
	return nil
}

func upgradeValues(upgrade *Upgrade) error {
	for _, v := range upgrade.SetKey {
		if v.Input.Enabled {
			var err error
			value := ""
			if v.Input.Password {
				v.Input.Twice = true
				value, err = utils.AcceptUserPassword(v.Input)
			} else {
				value, err = utils.AcceptUserInput(v.Input)
			}
			if err != nil {
				return err
			}
			v.Value = value
		}
		bytes, e := setValueByKey(upgrade.Values, v.Name, v.Value)
		if e != nil {
			return e
		}
		upgrade.Values = bytes
	}
	for _, v := range upgrade.ChangeKey {
		value, e := getValueByKey(upgrade.Values, v.Old)
		if e != nil {
			log.Errorf("Key: %s not found", v.Old)
			value, e = utils.AcceptUserInput(utils.Input{
				Tip:   fmt.Sprintf("Please value for Key: %s\n", v.New),
				Regex: ".+",
			})
			if e != nil {
				return e
			}
		}
		bytes, e := setValueByKey(upgrade.Values, v.New, value)
		if e != nil {
			return e
		}
		upgrade.Values = bytes
		upgrade.Values = deleteByKey(upgrade.Values, v.Old)
	}
	for _, v := range upgrade.DeleteKey {
		upgrade.Values = deleteByKey(upgrade.Values, v)
	}
	return nil
}

func getValueByKey(data []byte, key string) (string, error) {
	value, err := jsonparser.GetUnsafeString(data, strings.Split(key, ".")...)
	return string(value), err
}

func setValueByKey(data []byte, key, value string) ([]byte, error) {
	if value == "true" || value == "false" {
		return jsonparser.Set(data, []byte(value), strings.Split(key, ".")...)
	}
	if !(strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) {
		value = fmt.Sprintf(`"%s"`, value)
	}
	return jsonparser.Set(data, []byte(value), strings.Split(key, ".")...)
}

func deleteByKey(data []byte, key string) []byte {
	return jsonparser.Delete(data, strings.Split(key, ".")...)
}

func (u *Upgrader) Install() error {
	return nil
}
func (u *Upgrader) Uninstall() error {
	return nil
}

func checkJobDeleted(jobInterface v1.JobInterface) {
	jobList, err := jobInterface.List(context.Background(), meta_v1.ListOptions{})
	if len(jobList.Items) > 0 || err != nil {
		log.Infof("Deleting job,Please wait.")
		time.Sleep(5 * time.Second)
		checkJobDeleted(jobInterface)
	}
}

func (u *Upgrader) preUpgrade() error {
	namespace := ""
	for _, v := range u.Spec.Upgrade {
		if err := u.GetReleaseValues(v); err == nil {
			log.Debugf("Got %s,version %s", v.Name, v.InstalledVersion)
			constraintVersion := fmt.Sprintf("%s,<=%s", v.ConstraintVersion, v.Version)
			b, e := utils.CheckVersion(v.InstalledVersion, constraintVersion)
			if !b || e != nil {
				return fmt.Errorf("Can't auto upgrade of %s installed version. Want version %s,but got version %s.",
					v.Name, constraintVersion, v.InstalledVersion)
			}
			if v.Namespace != namespace {
				/*jobInterface := u.HelmInstall.KubeClient.BatchV1().Jobs(v.Namespace)
				jobList, err := jobInterface.List(meta_v1.ListOptions{})
				if err != nil {
					return err
				}
				log.Info("clean history jobs...")
				delOpts := &meta_v1.DeleteOptions{}
				for _, job := range jobList.Items {
					if job.Status.Active > 0 {
						log.Infof("job %s still active ignored..", job.Name)
					} else {
						if err := jobInterface.Delete(job.Name, delOpts); err != nil {
							return err
						}
						log.Successf("deleted job %s", job.Name)
					}
				}
				checkJobDeleted(jobInterface)
				namespace = v.Namespace*/
			}
		} else {
			log.Infof("Get Release %s error,Skip it. %s", v.Name, err)
		}
	}
	return nil
}

func (u *Upgrader) Run(args ...string) error {
	u.Init()
	if err := u.preUpgrade(); err != nil {
		return err
	}
	if err := u.Install(); err != nil {
		return err
	}
	if err := u.UpgradeReleases(); err != nil {
		return err
	}
	if err := u.Uninstall(); err != nil {
		return err
	}
	return nil
}
