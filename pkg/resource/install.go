package resource

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	c7ncfg "github.com/choerodon/c7nctl/pkg/config"
	c7nconsts "github.com/choerodon/c7nctl/pkg/consts"
	c7nctx "github.com/choerodon/c7nctl/pkg/context"
	c7nslaver "github.com/choerodon/c7nctl/pkg/slaver"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml_v2 "gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"text/template"
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
	t, err := c7nctx.GetReleaseTaskInfo(i.Namespace, r.Name)
	if err != nil {
		return err
	}
	if t.Status == c7nctx.UninitializedStatus {
		// 传入的参数是指针
		r.mergerResource(uc)
		t.Resource = *r.Resource

		if err = i.renderValues(r); err != nil {
			return err
		}
		t.Values = r.Values

		if err := i.render(r); err != nil {
			return err
		}

		// 保存渲染完成的 r
		t.Status = c7nctx.RenderedStatus
		if err = c7nctx.UpdateTaskToCM(i.Namespace, *t); err != nil {
			return err
		}
	}
	// 当 r 渲染完成但是没有完成安装——c7nctl install 会中断，二次执行
	if t.Status == c7nctx.RenderedStatus || t.Status == c7nctx.FailedStatus {
		r.Values = t.Values
		r.Resource = &t.Resource
		// 重新渲染 preCommand 等，避免在 TaskInfo 加入 PreCommand 导致循环依赖
		if err := i.render(r); err != nil {
			return err
		}
	}
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
	/*
		// TODO useless
				if r.Timeout > 0 {
					values = append(values, fmt.Sprintf("preJob.timeout=%d", r.Timeout))
				}
	*/
	var fileValsByte bytes.Buffer
	if uc == nil {
		fileVals, err := r.ValuesRaw(uc)
		if err != nil {
			return nil, err
		}
		fileValsByte, err = i.renderTpl(r.Name+"-file-values", fileVals)

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
		return errors.New(fmt.Sprintf("release %s values is empty", rls.Name))
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
func (i *InstallDefinition) getNamespace() string {
	return i.Namespace
}

func (i *InstallDefinition) withPrefix() string {
	if i.Prefix == "" {
		return ""
	}
	return i.Prefix + "-"
}

func (i *InstallDefinition) getReleaseName(rlsName string) string {
	return i.withPrefix() + rlsName
}

// TODO add storageClassName()
func (i *InstallDefinition) getStorageClass() string {
	return c7nctx.Ctx.UserConfig.GetStorageClassName()
}

func (i *InstallDefinition) getDatabaseUrl(rls string) string {
	return fmt.Sprintf(c7nconsts.DatabaseUrlTpl, i.getReleaseName("mysql"), i.getReleaseName(rls))
}

func (i *InstallDefinition) getResource(rls string) *c7ncfg.Resource {
	for _, r := range i.Spec.Release {
		if r.Name == rls {
			return r.Resource
		}
	}
	log.Fatal("Release cannot be empty")
	return nil
}

func (i *InstallDefinition) getReleaseValue(rls, value string) string {
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

func (i *InstallDefinition) encryptGitlabAccessToken() string {
	token := i.getReleaseValue("gitlab-service", "env.open.GITLAB_PRIVATETOKEN")
	dbKeyBase := i.getReleaseValue("gitlab", "core.env.GITLAB_SECRETS_DB_KEY_BASE")
	str := fmt.Sprintf("%s%s", token, dbKeyBase[:32])

	hash := sha256.New()
	hash.Write([]byte(str))

	// to lowercase hexits
	hex.EncodeToString(hash.Sum(nil))

	// to base64
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
