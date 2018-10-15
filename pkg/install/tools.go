package install

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/config"
	"github.com/choerodon/c7n/pkg/slaver"
	"github.com/vinkdong/gox/log"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"os"
	"regexp"
	"syscall"
	"time"
)

var Ctx Context

const (
	PvType             = "pv"
	PvcType            = "pvc"
	CRDType            = "crd"
	ReleaseTYPE        = "helm"
	SucceedStatus      = "succeed"
	FailedStatues      = "failed"
	staticLogName      = "c7n-logs"
	staticLogKey       = "logs"
	staticInstalledKey = "installed"
	randomLength       = 4
)

type Context struct {
	Client        kubernetes.Interface
	Namespace     string
	CommonLabels  map[string]string
	SlaverAddress string
	Slaver        *slaver.Slaver
	UserConfig    *config.Config
}

// i want use log but it make ...
type News struct {
	Name      string
	Namespace string
	RefName   string
	Type      string
	Status    string
	Reason    string
	Date      time.Time
	Resource  config.Resource
	Values    []ChartValue
	PreValue  []PreValue
}

type NewsResourceList struct {
	News []News `yaml:"logs"`
}

func (ctx *Context) SaveNews(news *News) error {
	data := ctx.GetOrCreateConfigMapData(staticLogName, staticLogKey)
	nr := &NewsResourceList{}
	yaml.Unmarshal([]byte(data), nr)
	news.Date = time.Now()
	if news.RefName == "" {
		news.RefName = news.Name
	}
	nr.News = append(nr.News, *news)
	newData, err := yaml.Marshal(nr)
	if err != nil {
		log.Error(err)
		return err
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, staticLogKey)

	if news.Status == SucceedStatus {
		ctx.SaveSucceed(news)
	}
	return nil
}

func (ctx *Context) SaveSucceed(news *News) error {

	news.Date = time.Now()
	nr := ctx.getSucceedData()
	nr.News = append(nr.News, *news)
	newData, err := yaml.Marshal(nr)
	if err != nil {
		log.Error(err)
		return err
	}
	ctx.saveConfigMapData(string(newData[:]), staticLogName, staticInstalledKey)
	return nil
}

func (ctx *Context) GetSucceed(name string, resourceType string) *News {
	nr := ctx.getSucceedData()
	for _, v := range nr.News {
		if v.Name == name && v.Type == resourceType {
			// todo: make sure gc effort
			p := v
			return &p
		}
	}
	return nil
}

func (ctx *Context) getSucceedData() *NewsResourceList {
	data := ctx.GetOrCreateConfigMapData(staticLogName, staticInstalledKey)
	nr := &NewsResourceList{}
	yaml.Unmarshal([]byte(data), nr)
	return nr
}

func (ctx *Context) GetOrCreateConfigMapData(cmName, cmKey string) string {
	cm, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Get(cmName, meta_v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("creating logs to cluster")
			cm = ctx.createNewsData()
		}
	}
	return cm.Data[cmKey]
}

func (ctx *Context) createNewsData() *v1.ConfigMap {

	data := make(map[string]string)
	data[staticLogKey] = ""
	cm := &v1.ConfigMap{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   staticLogName,
			Labels: ctx.CommonLabels,
		},
		Data: data,
	}
	configMap, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Create(cm)
	if err != nil {
		log.Error(err)
		os.Exit(122)
	}
	return configMap
}

func (ctx *Context) saveConfigMapData(data, cmName, cmKey string) *v1.ConfigMap {

	cm, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Get(cmName, meta_v1.GetOptions{})
	cm.Data[cmKey] = data
	configMap, err := ctx.Client.CoreV1().ConfigMaps(ctx.Namespace).Update(cm)
	if err != nil {
		log.Error(err)
		os.Exit(122)
	}
	return configMap
}

func IsNotFound(err error) bool {
	errorStatus, ok := err.(*errors.StatusError)
	if ok && errorStatus.Status().Code == 404 {
		return true
	}
	return false
}

func RandomString() string {
	len := randomLength
	bytes := make([]byte, len)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		bytes[i] = byte(97 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func AcceptUserPassword(input Input) (string, error) {
start:
	fmt.Print(input.Tip)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(input.Regex)
	if !r.MatchString(string(bytePassword[:])) {
		log.Error("password format not correct,try again")
		goto start
	}

	fmt.Print("enter again:")
	bytePassword2, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	if len(bytePassword2) != len(bytePassword) {
		log.Error("password length not match, please try again")
		goto start
	}
	for k, v := range bytePassword {
		if bytePassword2[k] != v {
			log.Error("password not match, please try again")
			goto start
		}
	}

	log.Info("waiting...")

	return string(bytePassword[:]), nil
}
