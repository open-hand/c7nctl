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
		{consts.ChoerodonRegister, "1.0.1"},
		{consts.ChoerodonPlatform, "1.0.0"},
		{consts.ChoerodonAdmin, "1.0.0"},
		{consts.ChoerodonIam, "1.0.4"},
		{consts.ChoerodonAsgard, "1.0.0"},
		{consts.ChoerodonSwagger, "1.0.0"},
		{consts.ChoerodonGateWay, "1.0.1"},
		{consts.ChoerodonOauth, "1.0.1"},
		{consts.ChoerodonMonitor, "1.0.0"},
		{consts.ChoerodonFile, "1.0.0"},
		{consts.ChoerodonMessage, "1.0.1"},
		{consts.DevopsService, "1.0.8"},
		{consts.WorkflowService, "1.0.0"},
		{consts.GitlabService, "1.0.1"},
		{consts.AgileService, "1.0.3"},
		{consts.TestManagerService, "1.0.2"},
		{consts.KnowledgebaseService, "1.0.1"},
		{consts.ElasticsearchKb, "1.0.0"},
		{consts.ProdRepoService, "1.0.0"},
		{consts.CodeRepoService, "1.0.7"},
		{consts.ChoerodonFrontHzero, "1.0.0"},
		{consts.ChoerodonFront, "1.0.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "1.0")
		// fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
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

func Test19to21ReleaseTag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{"go-register-server", "0.21.0"},
		{"manager-service", "0.21.0"},
		{"asgard-service", "0.21.1"},
		{"notify-service", "0.21.2"},
		{"base-service", "0.21.7"},
		{"api-gateway", "0.21.0"},
		{"oauth-server", "0.21.1"},
		{"file-service", "0.21.1"},
		{consts.DevopsService, "0.21.13"},
		{consts.WorkflowService, "0.21.1"},
		{consts.GitlabService, "0.21.0"},
		{consts.AgileService, "0.21.2"},
		{consts.TestManagerService, "0.21.1"},
		{consts.KnowledgebaseService, "0.21.0"},
		{consts.ElasticsearchKb, "0.21.0"},
		{consts.ChoerodonFront, "0.21.5"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "0.21")
		// fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
		if VersionOrdinal(version) != VersionOrdinal(app.Version) {
			t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
			continue
		}
		t.Logf("%s version is %s", app.Name, version)
	}
}

func TestHzeroReleaseTag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{"hzero-register", "0.22.2"},
		{"hzero-admin", "0.22.3"},
		{"hzero-iam", "0.22.5"},
		{"hzero-asgard", "0.22.5"},
		{"hzero-swagger", "0.22.1"},
		{"hzero-gateway", "0.22.4"},
		{"hzero-oauth", "0.22.2"},
		{"hzero-platform", "0.22.2"},
		{"hzero-monitor", "0.22.4"},
		{"hzero-file", "0.22.4"},
		{"hzero-message", "0.22.10"},
		{"hzero-front", "0.22.1"},
		{consts.DevopsService, "0.22.12"},
		{consts.WorkflowService, "0.22.2"},
		{consts.GitlabService, "0.22.1"},
		{consts.AgileService, "0.22.2"},
		{consts.TestManagerService, "0.22.2"},
		{consts.KnowledgebaseService, "0.22.1"},
		{consts.ElasticsearchKb, "0.22.1"},
		{consts.ChoerodonFront, "0.22.1"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "0.22")
		//fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)

		if VersionOrdinal(version) != VersionOrdinal(app.Version) {
			t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
			continue
		}
		t.Logf("%s version is %s", app.Name, version)
	}
}
