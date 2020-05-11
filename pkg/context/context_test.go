package context

import (
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/vinkdong/gox/log"
	"testing"
)

func TestGetNewsData(t *testing.T) {

	ctx := Context{
		KubeClient: client.GetKubeClient(),
		Namespace:  "test",
	}
	log.Info(ctx.GetOrCreateConfigMapData(staticLogName, staticLogKey))
}

func TestContext_GetJobInfo(t *testing.T) {
	mysql := Ctx.GetJobInfo("mysql")
	t.Log(mysql)
}

func TestSaveNewsData(t *testing.T) {
	ctx := Context{
		KubeClient: client.GetKubeClient(),
		Namespace:  "test",
	}

	news := &JobInfo{
		Name:      "testnews2",
		Namespace: "test",
		Type:      PvcType,
		Status:    FailedStatus,
		Reason:    "reason1 ",
	}
	ctx.SaveNews(news)
}

func TestRandomToken(t *testing.T) {
}
