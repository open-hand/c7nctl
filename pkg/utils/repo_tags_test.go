package utils

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	"testing"
)

func TestGetReleaseTag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{consts.ChoerodonRegister, "2.1.0"},
		{consts.ChoerodonPlatform, "2.1.0"},
		{consts.ChoerodonAdmin, "2.1.0"},
		{consts.ChoerodonIam, "2.1.0"},
		{consts.ChoerodonAsgard, "2.1.0"},
		{consts.ChoerodonSwagger, "2.1.0"},
		{consts.ChoerodonGateWay, "2.1.0"},
		{consts.ChoerodonOauth, "2.1.0"},
		{consts.ChoerodonMonitor, "2.1.0"},
		{consts.ChoerodonFile, "2.1.0"},
		{consts.ChoerodonMessage, "2.1.0"},
		{consts.DevopsService, "2.1.0"},
		{consts.WorkflowService, "2.1.0"},
		{consts.GitlabService, "2.1.0"},
		{consts.AgileService, "2.1.0"},
		{consts.TestManagerService, "2.0.0"},
		{consts.KnowledgebaseService, "2.0.0"},
		{consts.ElasticsearchKb, "2.0.0"},
		{consts.ProdRepoService, "2.0.0"},
		{consts.CodeRepoService, "2.0.0"},
		{consts.ChoerodonFrontHzero, "2.0.0"},
		{consts.ChoerodonFront, "2.0.0"},
		{consts.ChoerodonClusterAgent, "2.2.0"},
		{consts.ChoerodonIamServiceBusiness, "2.0.0"},
		{consts.DevopsServiceBusiness, "2.0.0"},
		{consts.AgileServiceBusiness, "2.0.0"},
		{consts.DocRepoService, "2.0.0"},
		//{consts.HrdsQA, "2.0.0"},
		//{consts.MarketService, "2.0.0"},
		{consts.ChoerodonFrontBusiness, "2.0.0"},
		{consts.TestManagerServiceBusiness, "2.0.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "2.2")
		//fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
		//if VersionOrdinal(version) != VersionOrdinal(app.Version) {
		//	t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
		//	continue
		//}
		fmt.Printf("helm pull vista-c7n/%s --version %s\n", app.Name, version)

		//t.Logf("%s: %s", app.Name, version)
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
