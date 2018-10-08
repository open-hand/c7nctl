package install

import (
	"testing"
	"github.com/choerodon/c7n/pkg/kube"
	"github.com/vinkdong/gox/log"
	"fmt"
)

func TestGetNewsData(t *testing.T)  {
	ctx := Context{
		Client: kube.GetClient(),
		Namespace: "test",
	}
	log.Info(ctx.GetOrCreateConfigMapData(staticLogName,staticLogKey))
}

func TestSaveNewsData(t *testing.T)  {
	ctx := Context{
		Client: kube.GetClient(),
		Namespace: "test",
	}

	news := &News{
		Name: "testnews2",
		Namespace: "test",
		Type: PvcType,
		Status:FailedStatues,
		Reason:"reason1 ",
	}
	ctx.SaveNews(news)
}

func TestRandomString(t *testing.T) {
	fmt.Println(RandomString(),RandomString())
	if RandomString() == RandomString(){
		t.Error("Func RandowString not work")
	}
}

