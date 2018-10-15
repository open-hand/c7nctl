package install

import (
	"bytes"
	"fmt"
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/helm"
	"github.com/choerodon/c7n/pkg/slaver"
	"github.com/vinkdong/gox/log"
	"k8s.io/client-go/kubernetes"
	"os"
	"text/template"
)

func (infra *InfraResource) executePreCommands() error {
	s := Ctx.Slaver
	for _, pi := range infra.PreInstall {
		for _, c := range pi.Commands {
			if err := s.ExecuteSql(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (infra *InfraResource) preparePersistence(client kubernetes.Interface, config *config.Config) error {
	getPvs := config.Spec.Persistence.GetPersistentVolumeSource
	namespace := config.Metadata.Namespace
	for _, persistence := range infra.Persistence {
		persistence.Client = client

		// check or create dir
		dir := slaver.Dir{
			Mode: persistence.Mode,
			Path: persistence.Path,
		}
		if err := Ctx.Slaver.MakeDir(dir); err != nil {
			return err
		}

		if err := persistence.CheckOrCreatePv(getPvs(persistence.Path)); err != nil {
			return err
		}
		if persistence.PvcEnabled {
			persistence.Namespace = namespace
			err := persistence.CheckOrCreatePvc()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (infra *InfraResource) executePreValues() error {
	return infra.PreValues.prepareValues()
}

func (infra *InfraResource) renderValue(tplString string) string {
	tpl, err := template.New(infra.Name).Parse(tplString)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	var data bytes.Buffer
	err = tpl.Execute(&data, infra)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	return data.String()
}

func (infra *InfraResource) GetPreValue(key string) string {
	return infra.PreValues.getValues(key)
}

func (infra *InfraResource) GetRequirement(app, key string) string {
	return "ok"
}

// convert yml values to values list as xxx=yyy
func (infra *InfraResource) HelmValues() ([]string, []ChartValue) {
	values := make([]string, len(infra.Values))
	// store values for feature use
	cvList := make([]ChartValue, len(infra.Values))
	for k, v := range infra.Values {
		value := ""
		if v.Input.Enabled {
			password, err := AcceptUserPassword(v.Input)
			if err != nil {
				log.Error(err)
				os.Exit(128)
			}
			value = password
		} else {
			value = infra.renderValue(v.Value)
		}
		values[k] = fmt.Sprintf("%s=%s", v.Name, value)
		v.Value = value
		cvList[k] = v
	}
	return values, cvList
}

// install infra
func (infra *InfraResource) Install() error {
	values, cvList := infra.HelmValues()
	chartArgs := helm.ChartArgs{
		ReleaseName: infra.Name,
		Namespace:   infra.Namespace,
		RepoUrl:     infra.RepoURL,
		Verify:      false,
		Version:     infra.Version,
		ChartName:   infra.Chart,
	}
	log.Infof("installing %s", infra.Name)
	err := infra.Client.InstallRelease(values, chartArgs)

	news := &News{
		Name:      infra.Name,
		Namespace: infra.Namespace,
		RefName:   infra.Name,
		Status:    FailedStatues,
		Type:      ReleaseTYPE,
		Resource:  infra.Resource,
		Values:    cvList,
	}
	defer Ctx.SaveNews(news)

	if err != nil {
		news.Reason = err.Error()
		return err
	}
	news.Status = SucceedStatus
	return nil
}

func (infra *InfraResource) CheckInstall() error {
	news := Ctx.GetSucceed(infra.Name, ReleaseTYPE)
	// 初始化value
	if err := infra.executePreValues(); err != nil {
		return err
	}
	if news != nil {
		log.Infof("using exist release %s", news.RefName)
		return nil
	}
	// 执行安装前命令
	if err := infra.executePreCommands(); err != nil {
		return err
	}
	return infra.Install()
}
