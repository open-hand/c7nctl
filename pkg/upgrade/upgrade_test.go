package upgrade

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/config"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/choerodon/c7nctl/pkg/common"
	"github.com/ghodss/yaml"
	"github.com/vinkdong/gox/log"
)

var data = `
# 这里是注释
env:
  config:
    GITLAB_DEFAULT_CAN_CREATE_GROUP: true
    GITLAB_EMAIL_FROM: git.sys@example.com
    GITLAB_EXTERNAL_URL: http://git.staging.saas.hand-china.com/
    GITLAB_TIMEZONE: Asia/Shanghai
    MYSQL_DATABASE: gitlabhq_production
    MYSQL_HOST: gitlab-mysql
    MYSQL_PASSWORD: password
    MYSQL_PORT: 3306
    MYSQL_USERNAME: root
    NODE_EXPORTER_ENABLE: false
    PROMETHEUS_ENABLE: false
    REDIS_HOST: gitlab-redis
    REDIS_PORT: 6379
    SMTP_ENABLE: false
ingress:
  enabled: true
  hosts:
  - git.staging.saas.hand-china.com
persistence:
  enabled: true
  existingClaim: gitlab-pvc`

func initData() []byte {
	r, _ := yaml.YAMLToJSON([]byte(data))
	return r
}

func TestGetRelease(t *testing.T) {
	u := Upgrader{}
	upgrade := Upgrade{
		Name: "mysql-test",
	}
	e := u.GetReleaseValues(&upgrade)
	if e != nil {
		t.Error(e)
	} else {
		s, e := getValueByKey(upgrade.Values, "env.MYSQL_ROOT_PASSWORD")
		if e != nil {
			t.Error(e)
		}
		t.Log(s)
	}
}

func TestGetValueByKey(t *testing.T) {
	type args struct {
		data []byte
		key  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "string",
			args: args{
				data: initData(),
				key:  "env.config.MYSQL_PASSWORD",
			},
			want:    "password",
			wantErr: false,
		},
		{name: "int",
			args: args{
				data: initData(),
				key:  "env.config.MYSQL_PORT",
			},
			want:    "3306",
			wantErr: false,
		},
		{name: "bool",
			args: args{
				data: initData(),
				key:  "ingress.enabled",
			},
			want:    "true",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getValueByKey(tt.args.data, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValueByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetValueByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetValueByKey(t *testing.T) {
	var setData = `
env:
  config:
    GITLAB_DEFAULT_CAN_CREATE_GROUP: true
    GITLAB_EMAIL_FROM: git.sys@example.com
    GITLAB_EXTERNAL_URL: http://git.staging.saas.hand-china.com/
    GITLAB_TIMEZONE: Asia/Shanghai
    MYSQL_DATABASE: gitlabhq_production
    MYSQL_HOST: gitlab-mysql
    MYSQL_PASSWORD: password
    MYSQL_PORT: 3306
    MYSQL_USERNAME: root
    NODE_EXPORTER_ENABLE: false
    PROMETHEUS_ENABLE: false
    REDIS_HOST: gitlab-redis
    REDIS_PORT: 6379
    SMTP_ENABLE: false
ingress:
  enabled: false
  hosts:
  - git.staging.saas.hand-china.com
persistence:
  enabled: true
  existingClaim: gitlab-pvc`

	r, _ := yaml.YAMLToJSON([]byte(setData))
	type args struct {
		data  []byte
		value string
		key   string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "string",
			args: args{
				data:  initData(),
				key:   "ingress.enabled",
				value: "false",
			},
			want:    r,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setValueByKey(tt.args.data, tt.args.value, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetValueByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.Compare(string(got), string(tt.want)) != 0 {
				t.Errorf("SetValueByKey() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestDeleteByKey(t *testing.T) {
	var delete = `
env:
  config:
    GITLAB_DEFAULT_CAN_CREATE_GROUP: true
    GITLAB_EMAIL_FROM: git.sys@example.com
    GITLAB_EXTERNAL_URL: http://git.staging.saas.hand-china.com/
    GITLAB_TIMEZONE: Asia/Shanghai
    MYSQL_DATABASE: gitlabhq_production
    MYSQL_HOST: gitlab-mysql
    MYSQL_PASSWORD: password
    MYSQL_PORT: 3306
    MYSQL_USERNAME: root
    NODE_EXPORTER_ENABLE: false
    PROMETHEUS_ENABLE: false
    REDIS_HOST: gitlab-redis
    REDIS_PORT: 6379
    SMTP_ENABLE: false
ingress:
  hosts:
  - git.staging.saas.hand-china.com
persistence:
  enabled: true
  existingClaim: gitlab-pvc`

	r, _ := yaml.YAMLToJSON([]byte(delete))
	type args struct {
		data []byte
		key  string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{name: "string",
			args: args{
				data: initData(),
				key:  "ingress.enabled",
			},
			want: r,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteByKey(tt.args.data, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpgrader_UpgradeRelease(t *testing.T) {
	log.EnableDebug()
	u := &Upgrader{
		Spec: Spec{
			Basic: Basic{
				RepoURL: "https://openchart.choerodon.com.cn/choerodon/c7n/",
			},
		},
	}
	u.Init()
	up := &Upgrade{
		Name:    "mysql-test",
		Chart:   "mysql",
		Version: "0.1.0",
		SetKey: []*SetKey{
			&SetKey{
				Name:  "env.aa",
				Value: "123,123,env.aa",
			},
			&SetKey{
				Name:  "env.bb",
				Value: "bb",
			},
			&SetKey{
				Name:  "env.cc",
				Value: "true",
			},
			&SetKey{
				Name:  "env.dd",
				Value: "123",
			},
		},
		ChangeKey: []*ChangeKey{
			&ChangeKey{
				Old: "env.cc",
				New: "env.ccc",
			},
		},
		DeleteKey: []string{
			"env.dd",
			"env.sss",
		},
	}
	release := upgradeRelease(u, up)
	fmt.Print(release)
}

func Test_CheckVersion(t *testing.T) {
	b, e := common.CheckVersion("0.11.1", ">=0.11.0")
	fmt.Print(b, e)
}

func Test(t *testing.T) {
	var (
		data []byte
		err  error
	)
	resp, err := http.Get(fmt.Sprintf("%s%s", "https://file.choerodon.com.cn/choerodon-install", "/upgrade.yml"))
	if err != nil {
		log.Error(err)
		os.Exit(127)
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Get resource %s failed", "/upgrade.yml")
		log.Error(err)
		os.Exit(127)
	}
	fmt.Printf("%s", string(data))
	var u Upgrader
	yaml.Unmarshal(data, &u)
	fmt.Printf("%s", string(data))
}

func TestUpgrader_Run(t *testing.T) {
	log.EnableDebug()
	r := config.ResourceDefinition{}
	r.LocalFile = "../../upgrade.yml"
	data, _ := r.GetUpgradeResourceDate("0.12")
	u := Upgrader{}
	yaml.Unmarshal(data, &u)
	e := u.Run()
	fmt.Println(e)
}

func Test_a(t *testing.T) {
	for i := 0; i != 10; i = i + 1 {
		fmt.Fprintf(os.Stdout, "result is %d\r", i)
		time.Sleep(time.Second * 1)
	}
	fmt.Println("over")
}
