package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/chr4/pwgen"
	"github.com/vinkdong/gox/log"
	"github.com/vinkdong/gox/random"
	"math/rand"
	"os"
	"text/template"
	"time"
)

const DatabaseUrlTpl = "jdbc:mysql://%s:3306/%s?useUnicode=true&characterEncoding=utf-8&useSSL=false&useInformationSchema=true&remarks=true&allowMultiQueries=true&serverTimezone=Asia/Shanghai"

/*
func RenderInstallDef(tplStr string) []byte {
	// context = uc
	tpl, err := template.New("c7n-install-def").Funcs(c7nFunc).Parse(tplStr)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	var intDef bytes.Buffer

	err = tpl.Execute(&intDef, context.Ctx.UserConfig)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	return intDef.Bytes()
}*/

func RenderRelease(rlsName, tplStr string) []byte {
	rls := renderTpl(rlsName, tplStr)
	return rls.Bytes()
}
func RenderReleaseValue(values, tplStr string) string {
	rlsValues := renderTpl(values+"-values", tplStr)
	return rlsValues.String()
}

func renderTpl(name, tplStr string) bytes.Buffer {
	tpl, err := template.New(name).Funcs(c7nFunc).Parse(tplStr)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	var result bytes.Buffer
	err = tpl.Execute(&result, context.Ctx.UserConfig)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	return result
}

var c7nFunc = template.FuncMap{
	"getResource":              GetResource,
	"withPrefix":               WithPrefix,
	"getReleaseName":           getReleaseName,
	"getStorageClass":          getStorageClass,
	"randomToken":              randomToken,
	"getImageRepo":             getImageRepo,
	"randomLowCaseToken":       randomLowCaseToken,
	"getDatabaseUrl":           getDatabaseUrl,
	"getNamespace":             getNameSpace,
	"getReleaseValue":          getReleaseValue,
	"generateAlphaNum":         generateAlphaNum,
	"encryptGitlabAccessToken": encryptGitlabAccessToken,
}

func getImageRepo(rls string) string {
	return "registry.cn-shanghai.aliyuncs.com/c7n/" + rls
}

// release 信息会被保存在 cm 中，所以将 resource 合并保存在 release 中
// 应该警惕获取的 rls 未安装导致的 value 不完整
func GetResource(rlsName string) *config.Resource {
	//news := context.Ctx.GetSucceed(rlsName, context.ReleaseTYPE)
	_, ji := context.Ctx.GetJobInfo(rlsName)
	// get info from succeed
	if ji.Name != "" {
		return &ji.Resource
	} else if r, ok := context.Ctx.UserConfig.Spec.Resources[rlsName]; ok {
		return r
	}
	errMsg := fmt.Sprintf("can't get required resource [%s]", rlsName)
	log.Error(errMsg)
	context.Ctx.Metrics.ErrorMsg = append(context.Ctx.Metrics.ErrorMsg, errMsg)
	context.Ctx.CheckExist(188)

	return nil
}

func getReleaseValue(rls, valueName string) string {
	_, ji := context.Ctx.GetJobInfo(rls)
	var err error
	if ji.Name != "" {
		for _, key := range ji.Values {
			if key.Name == valueName && key.Value != "" {
				return key.Value
			}
			err = errors.New(fmt.Sprintf("can't get required value %s in release %s", valueName, rls))
		}
	}
	err = errors.New(fmt.Sprintf("can't found release %s which installed", rls))
	log.Error(err)
	context.Ctx.Metrics.ErrorMsg = append(context.Ctx.Metrics.ErrorMsg, err.Error())
	context.Ctx.CheckExist(188)

	return ""
}

func WithPrefix() string {
	if context.Ctx.Prefix == "" {
		return ""
	}
	return context.Ctx.Prefix + "-"
}

func getReleaseName(rlsName string) string {
	return WithPrefix() + rlsName
}

func getStorageClass() string {
	return context.Ctx.UserConfig.GetStorageClassName()
}

func randomToken(length int) string {
	return GenerateRunnerToken(length)
}

func randomLowCaseToken(length int) string {
	return GenerateRunnerToken(length)
}

func getDatabaseUrl(rls string) string {
	return fmt.Sprintf(DatabaseUrlTpl, getReleaseName("mysql"), getReleaseName(rls))
}

func getNameSpace() string {
	return context.Ctx.Namespace
}

func generateAlphaNum(length int) string {
	return pwgen.AlphaNum(length)
}

func encryptGitlabAccessToken() string {
	token := getReleaseValue("gitlab-service", "env.open.GITLAB_PRIVATETOKEN")
	dbKeyBase := getReleaseValue("gitlab", "core.env.GITLAB_SECRETS_DB_KEY_BASE")
	str := fmt.Sprintf("%s%s", token, dbKeyBase[:32])

	hash := sha256.New()
	hash.Write([]byte(str))

	// to lowercase hexits
	hex.EncodeToString(hash.Sum(nil))

	// to base64
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func RandomToken(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		random.Seed(time.Now().UnixNano())
		op := random.RangeIntInclude(random.Slice{Start: 48, End: 57},
			random.Slice{Start: 65, End: 90}, random.Slice{Start: 97, End: 122})
		b[i] = byte(op) //A=65 and Z = 65+25
	}
	return string(b)
}

func GenerateRunnerToken(length int) string {
	b := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		random.Seed(time.Now().UnixNano())
		op := random.RangeIntInclude(random.Slice{Start: 48, End: 57},
			random.Slice{Start: 97, End: 122})
		b[i] = byte(op) //A=65 and Z = 65+25
	}
	return string(b)
}
