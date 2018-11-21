package install

import (
	"fmt"
	"github.com/choerodon/c7n/pkg/kube"
	"github.com/vinkdong/gox/log"
	"testing"
)

func TestGetNewsData(t *testing.T) {
	ctx := Context{
		Client:    kube.GetClient(),
		Namespace: "test",
	}
	log.Info(ctx.GetOrCreateConfigMapData(staticLogName, staticLogKey))
}

func TestSaveNewsData(t *testing.T) {
	ctx := Context{
		Client:    kube.GetClient(),
		Namespace: "test",
	}

	news := &News{
		Name:      "testnews2",
		Namespace: "test",
		Type:      PvcType,
		Status:    FailedStatus,
		Reason:    "reason1 ",
	}
	ctx.SaveNews(news)
}

func TestRandomString(t *testing.T) {
	fmt.Println(RandomString(), RandomString())
	if RandomString() == RandomString() {
		t.Error("Func RandowString not work")
	}
}

func TestRandomToken(t *testing.T) {
	fmt.Println(RandomToken(17), RandomToken(12))
}

func TestCheckMatch(t *testing.T)  {
	input := Input{
		Regex: "^\\S{8}$",
		Include: []KV{KV{Name:"必须包含数字",Value:"\\d+"},KV{Name:"必须包含大写",Value:"[A-Z]+"}},
	}

	if CheckMatch("abcdef",input) {
		t.Error("regex check error")
	}

	if !CheckMatch("ce123Qwe",input) {
		t.Error("regex check error")
	}

	if CheckMatch("ce123456",input) {
		t.Error("regex check error")
	}

	if CheckMatch("cdsfaQes",input) {
		t.Error("regex check error")
	}
}
