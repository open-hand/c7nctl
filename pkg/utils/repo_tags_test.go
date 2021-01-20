package utils

import (
	"github.com/choerodon/c7nctl/pkg/common/consts"
	"testing"
)

func TestGetReleaseTag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{consts.ChoerodonRegister, "0.24.1"},
		{consts.ChoerodonPlatform, "0.24.0"},
		{consts.ChoerodonAdmin, "0.24.0"},
		{consts.ChoerodonIam, "0.24.1"},
		{consts.ChoerodonAsgard, "0.24.0"},
		{consts.ChoerodonSwagger, "0.24.0"},
		{consts.ChoerodonGateWay, "0.24.0"},
		{consts.ChoerodonOauth, "0.24.2"},
		{consts.ChoerodonMonitor, "0.24.0"},
		{consts.ChoerodonFile, "0.24.0"},
		{consts.ChoerodonMessage, "0.24.0"},
		{consts.DevopsService, "0.24.2"},
		{consts.WorkflowService, "0.24.0"},
		{consts.GitlabService, "0.24.0"},
		{consts.AgileService, "0.24.2"},
		{consts.TestManagerService, "0.24.1"},
		{consts.KnowledgebaseService, "0.24.0"},
		{consts.ElasticsearchKb, "0.24.0"},
		{consts.ProdRepoService, "0.24.2"},
		{consts.CodeRepoService, "0.24.1"},
		{consts.ChoerodonFrontHzero, "0.24.0"},
		{consts.ChoerodonFront, "0.24.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "0.24")
		if VersionOrdinal(version) != VersionOrdinal(app.Version) {
			t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
			continue
		}
		t.Logf("%s version is %s", app.Name, version)
	}
}

func TestCheckMatch2(t *testing.T) {
	url, path := matchChartRepo(consts.DefaultRepoUrl)
	t.Logf("url: %s path: %s", url, path)
}
