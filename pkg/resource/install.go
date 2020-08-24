package resource

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	c7nconsts "github.com/choerodon/c7nctl/pkg/common/consts"
	c7ncfg "github.com/choerodon/c7nctl/pkg/config"
	c7nslaver "github.com/choerodon/c7nctl/pkg/slaver"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml_v2 "gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"text/template"
	"time"
)

const (
	eurekaClientServerUrlTpl = "%s://%s:8000/eureka/"
	resourceDomainUrlTpl     = "%s://%s"
)

type InstallDefinition struct {
	// api 版本
	Version string
	// Choerodon 平台版本
	PaaSVersion string
	Metadata    Metadata
	Spec        Spec
	// TODO REMOVE
	CommonLabels       map[string]string
	DefaultAccessModes []v1.PersistentVolumeAccessMode `yaml:"accessModes"`
	StorageClass       string
	SkipInput          bool
	Timeout            int
	Prefix             string
	Namespace          string
	Mail               string
}

type Metadata struct {
	Name      string
	Namespace string
}

type Spec struct {
	Basic     Basic
	Resources v1.ResourceRequirements
	Release   []*Release
	Runner    *Release `json:"runner"`
	Component []*Release
}

type Basic struct {
	RepoURL string
	Slaver  c7nslaver.Slaver
}

// TODO 需要初始化 namespace prefix 等参数
func (i *InstallDefinition) RenderRelease(r *Release, uc *c7ncfg.C7nConfig) error {
	t, err := r.Client.GetTaskInfoFromCM(i.Namespace, r.Name)

	if err != nil {
		// 创建一些基本错误类型
		if err.Error() == "Task info is not found" {
			t = c7nclient.TaskInfo{
				Name:      r.Name,
				Namespace: r.Namespace,
				RefName:   "",
				Status:    c7nconsts.UninitializedStatus,
				Type:      c7nconsts.StaticReleaseKey,
				Date:      time.Now(),
				Version:   r.Version,
				Prefix:    r.Prefix,
			}
		} else {
			return err
		}
	}
	if t.Status == c7nconsts.UninitializedStatus {
		// 传入的参数是指针
		r.mergerResource(uc)
		if err = i.renderValues(r); err != nil {
			return err
		}
		if err := i.render(r); err != nil {
			return err
		}

		// 保存渲染完成的 r
		t.Resource = *r.Resource
		t.Values = r.Values
		t.Status = c7nconsts.RenderedStatus
		if err = r.Client.SaveTaskInfoToCM(i.Namespace, t); err != nil {
			return err
		}
	} else {
		// 当 r 渲染完成但是没有完成安装——c7nctl install 会中断，二次执行
		r.Values = t.Values
		r.Resource = &t.Resource
		// 重新渲染 preCommand 等，避免在 TaskInfo 加入 PreCommand 导致循环依赖
		if err := i.render(r); err != nil {
			return err
		}
	}
	log.Infof("Successfully render Release %s", r.Name)
	return nil
}

func (i *InstallDefinition) RenderComponent(rls *Release) error {
	rlsByte, _ := yaml_v2.Marshal(rls)
	renderedRls, err := i.renderTpl(rls.Name, string(rlsByte))
	if err != nil {
		return err
	}
	return yaml_v2.Unmarshal(renderedRls.Bytes(), rls)
}

//
func (i *InstallDefinition) RenderHelmValues(r *Release, uc *c7ncfg.C7nConfig) (map[string]interface{}, error) {
	rlsVals := r.HelmValues()
	var fileValsByte bytes.Buffer
	if uc != nil {
		fileVals, err := r.ValuesRaw(uc)
		if err != nil {
			return nil, err
		}
		fileValsByte, err = i.renderTpl(r.Name+"-file-values", fileVals)
		if err != nil {
			return nil, err
		}
	}

	return c7nutils.Vals(rlsVals, fileValsByte.String())
}

// 渲染 release
func (i *InstallDefinition) render(r *Release) error {
	rlsByte, _ := yaml_v2.Marshal(r)
	renderedRls, err := i.renderTpl(r.Name, string(rlsByte))
	if err != nil {
		return err
	}
	if err := yaml_v2.Unmarshal(renderedRls.Bytes(), r); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Unmarshal Release %s failed", r))
	}
	return nil
}

// 传指针的方式好呢，还是返回值的方式好？
//
// 在渲染 release 前将 values 渲染完成
// 获取用户输入或者根据 value 的模版值渲染
func (i *InstallDefinition) renderValues(rls *Release) error {
	if rls.Values == nil {
		log.Debugf("release %s values is empty", rls.Name)
		return nil
	}
	for idx, v := range rls.Values {
		// 输入 value
		if v.Input.Enabled && !i.SkipInput {
			var err error
			var value string
			if v.Input.Password {
				v.Input.Twice = true
				value, err = c7nutils.AcceptUserPassword(v.Input)
			} else {
				value, err = c7nutils.AcceptUserInput(v.Input)
			}
			if err != nil {
				return err
			}
			// v.Values 是复制
			rls.Values[idx].Value = value
		} else {
			v, err := i.renderTpl(v.Name+"-values", v.Value)
			if err != nil {
				return err
			}
			rls.Values[idx].Value = v.String()
		}
	}
	return nil
}

// 根据模版和 InstallDefinition 渲染
func (i *InstallDefinition) renderTpl(name, tplStr string) (bytes.Buffer, error) {
	tpl, err := template.New(name).Funcs(c7nutils.C7nFunc).Parse(tplStr)
	if err != nil {
		return bytes.Buffer{}, err
	}
	var result bytes.Buffer
	err = tpl.Execute(&result, i)
	if err != nil {
		return bytes.Buffer{}, errors.WithMessage(err, fmt.Sprintf("Failed to render release %s", name))
	}
	return result, nil
}

/*
  template 内嵌函数
*/
func (i *InstallDefinition) GetNamespace() string {
	return i.Namespace
}

func (i *InstallDefinition) WithPrefix() string {
	if i.Prefix == "" {
		return ""
	}
	return i.Prefix + "-"
}

func (i *InstallDefinition) GetReleaseName(rlsName string) string {
	return i.WithPrefix() + rlsName
}

// TODO add storageClassName()
func (i *InstallDefinition) GetStorageClass() string {
	//return c7nctx.Ctx.UserConfig.GetStorageClassName()
	return i.StorageClass
}

func (i *InstallDefinition) GetDatabaseUrl(rls string) string {
	return fmt.Sprintf(c7nconsts.DatabaseUrlTpl, i.GetReleaseName("mysql"), i.GetReleaseName(rls))
}

func (i *InstallDefinition) GetResource(rls string) *c7ncfg.Resource {
	for _, r := range i.Spec.Release {
		if r.Name == rls {
			return r.Resource
		}
	}
	log.Fatal("Release cannot be empty")
	return nil
}

func (i *InstallDefinition) GetReleaseValue(rls, value string) string {
	for _, r := range i.Spec.Release {
		if r.Name == rls {
			for _, v := range r.Values {
				if v.Name == value {
					return v.Value
				}
			}
			log.WithField("Release values", value).Fatal("Release value cannot be empty")
		}
	}
	log.WithField("Release", rls).Fatal("Release cannot be empty")
	return ""
}

func (i *InstallDefinition) EncryptGitlabAccessToken() string {
	token := i.GetReleaseValue("gitlab-service", "env.open.GITLAB_PRIVATETOKEN")
	dbKeyBase := i.GetReleaseValue("gitlab", "core.env.GITLAB_SECRETS_DB_KEY_BASE")
	str := fmt.Sprintf("%s%s", token, dbKeyBase[:32])

	hash := sha256.New()
	hash.Write([]byte(str))

	// to lowercase hexits
	hex.EncodeToString(hash.Sum(nil))

	// to base64
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (i *InstallDefinition) GetPersistence(rls string, index int) *Persistence {
	for _, r := range i.Spec.Release {
		if r.Name == rls && len(r.Persistence) > index {
			return r.Persistence[index]
		}
	}
	log.WithField("Release", rls).Fatal("Release cannot be empty")
	return nil
}

func (i *InstallDefinition) GetRunnerPersistence(index int) *Persistence {
	if len(i.Spec.Runner.Persistence) > index {
		return i.Spec.Runner.Persistence[index]
	}
	log.WithField("Release", "gitlab-runner").Fatal("Release cannot be empty")
	return nil
}

func (i *InstallDefinition) GetRunnerValues(values string) string {
	for _, v := range i.Spec.Runner.Values {
		if v.Name == values {
			return v.Value
		}
	}
	return ""
}

func (i *InstallDefinition) GetEurekaUrl() string {
	for _, r := range i.Spec.Release {
		if r.Name == c7nconsts.HzeroRegister {
			return fmt.Sprintf(eurekaClientServerUrlTpl, r.Resource.Schema, r.Resource.Host)
		}
	}
	return ""
}

func (i *InstallDefinition) GetResourceDomainUrl(rls string) string {
	for _, r := range i.Spec.Release {
		if r.Name == rls {
			return fmt.Sprintf(resourceDomainUrlTpl, r.Resource.Schema, r.Resource.Domain)
		}
	}
	return ""
}
