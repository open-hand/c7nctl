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
		{consts.ChoerodonRegister, "0.23.1"},
		{consts.ChoerodonPlatform, "0.23.4"},
		{consts.ChoerodonAdmin, "0.23.1"},
		{consts.ChoerodonIam, "0.23.7"},
		{consts.ChoerodonAsgard, "0.23.5"},
		{consts.ChoerodonSwagger, "0.23.1"},
		{consts.ChoerodonGateWay, "0.23.2"},
		{consts.ChoerodonOauth, "0.23.1"},
		{consts.ChoerodonMonitor, "0.23.1"},
		{consts.ChoerodonFile, "0.23.1"},
		{consts.ChoerodonMessage, "0.23.7"},
		{consts.DevopsService, "0.23.7"},
		{consts.WorkflowService, "0.23.2"},
		{consts.GitlabService, "0.23.1"},
		{consts.AgileService, "0.23.6"},
		{consts.TestManagerService, "0.23.2"},
		{consts.KnowledgebaseService, "0.23.1"},
		{consts.ElasticsearchKb, "0.23.0"},
		{consts.ProdRepoService, "0.23.4"},
		{consts.CodeRepoService, "0.23.2"},
		{consts.ChoerodonFrontHzero, "0.23.1"},
		{consts.ChoerodonFront, "0.23.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(app.Name, "0.23")
		if VersionOrdinal(version) != VersionOrdinal(app.Version) {
			t.Errorf("%s  %s is not newest version", app.Name, app.Version)
			continue
		}
		t.Logf("%s version is %s", app.Name, version)
	}
}
