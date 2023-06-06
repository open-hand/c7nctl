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
		{consts.ChoerodonRegister, "2.4.0"},
		{consts.ZknowPlatform, "2.4.0"},
		{consts.ZknowAdmin, "2.4.0"},
		{consts.ZknowIam, "2.4.0"},
		{consts.ZknowAsgard, "2.4.0"},
		{consts.ChoerodonSwagger, "2.4.0"},
		{consts.ZknowGateway, "2.4.0"},
		{consts.ZknowOauth, "2.4.0"},
		{consts.ChoerodonMonitor, "2.4.0"},
		{consts.ChoerodonBase, "2.4.0"},
		{consts.ZknowFile, "2.4.0"},
		{consts.ZknowMessage, "2.4.0"},
		{consts.DevopsService, "2.4.0"},
		{consts.WorkflowService, "2.4.0"},
		{consts.GitlabService, "2.4.0"},
		{consts.AgileService, "2.4.0"},
		{consts.TestManagerService, "2.4.0"},
		{consts.KnowledgebaseService, "2.4.0"},
		{consts.ElasticsearchKb, "2.4.0"},
		{consts.ProdRepoService, "2.4.0"},
		{consts.CodeRepoService, "2.4.0"},
		{consts.ChoerodonFrontHzero, "2.4.0"},
		{consts.ChoerodonFront, "2.4.0"},
		{consts.ChoerodonClusterAgent, "2.4.0"},
		{consts.ChoerodonIamServiceBusiness, "2.4.0"},
		{consts.DevopsServiceBusiness, "2.4.0"},
		{consts.AgileServiceBusiness, "2.4.0"},
		{consts.DocRepoService, "2.4.0"},
		//{consts.HrdsQA, "2.0.0"},
		//{consts.MarketService, "2.0.0"},
		{consts.ChoerodonFrontBusiness, "2.4.0"},
		{consts.TestManagerServiceBusiness, "2.4.0"},
		{consts.ChoerodonFrontBase, "2.4.0"},
		{consts.ChoerodonFrontCodeRepo, "2.4.0"},
		{consts.ChoerodonFrontDevops, "2.4.0"},
		{consts.ChoerodonFrontDocRepo, "2.4.0"},
		{consts.ChoerodonFrontProdRepo, "2.4.0"},
		{consts.ChoerodonFrontAgilePro, "2.4.0"},
		{consts.ChoerodonFrontKnowledgebase, "2.4.0"},
		{consts.ChoerodonFrontMobile, "2.4.0"},
		{consts.ChoerodonFrontBaseBusiness, "2.4.0"},
		{consts.ChoerodonFrontAsgard, "2.4.0"},
		{consts.ChoerodonFrontManager, "2.4.0"},
		{consts.ChoerodonFrontNotify, "2.4.0"},
		{consts.ChoerodonFrontTestPro, "2.4.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "2.4")
		//fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
		//if VersionOrdinal(version) != VersionOrdinal(app.Version) {
		//	t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
		//	continue
		//}
		fmt.Printf("%s : %s\n", app.Name, version)

		//t.Logf("%s: %s", app.Name, version)
	}
}

func TestGetRelease23Tag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{consts.ChoerodonRegister, "2.2.0"},
		{consts.ZknowPlatform, "2.2.0"},
		{consts.ZknowAdmin, "2.2.0"},
		{consts.ZknowIam, "2.2.0"},
		{consts.ZknowAsgard, "2.2.0"},
		{consts.ChoerodonSwagger, "2.2.0"},
		{consts.ZknowGateway, "2.2.0"},
		{consts.ZknowOauth, "2.2.0"},
		{consts.ChoerodonMonitor, "2.2.0"},
		{consts.ChoerodonBase, "2.2.0"},
		{consts.ZknowFile, "2.2.0"},
		{consts.ZknowMessage, "2.2.0"},
		{consts.DevopsService, "2.2.0"},
		{consts.WorkflowService, "2.2.0"},
		{consts.GitlabService, "2.2.0"},
		{consts.AgileService, "2.2.0"},
		{consts.TestManagerService, "2.2.0"},
		{consts.KnowledgebaseService, "2.2.0"},
		{consts.ElasticsearchKb, "2.2.0"},
		{consts.ProdRepoService, "2.2.0"},
		{consts.CodeRepoService, "2.2.0"},
		{consts.ChoerodonFrontHzero, "2.2.0"},
		{consts.ChoerodonFront, "2.2.0"},
		{consts.ChoerodonClusterAgent, "2.2.0"},
		{consts.ChoerodonIamServiceBusiness, "2.2.0"},
		{consts.DevopsServiceBusiness, "2.2.0"},
		{consts.AgileServiceBusiness, "2.2.0"},
		{consts.DocRepoService, "2.2.0"},
		//{consts.HrdsQA, "2.0.0"},
		//{consts.MarketService, "2.0.0"},
		{consts.ChoerodonFrontBusiness, "2.2.0"},
		{consts.TestManagerServiceBusiness, "2.2.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "2.3.0-alpha")
		//fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
		//if VersionOrdinal(version) != VersionOrdinal(app.Version) {
		//	t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
		//	continue
		//}
		fmt.Printf("%s : %s\n", app.Name, version)

		//t.Logf("%s: %s", app.Name, version)
	}
}

func TestGetRelease22Tag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{consts.ChoerodonRegister, "2.2.0"},
		{consts.ChoerodonPlatform, "2.2.0"},
		{consts.ChoerodonAdmin, "2.2.0"},
		{consts.ChoerodonIam, "2.2.0"},
		{consts.ChoerodonAsgard, "2.2.0"},
		{consts.ChoerodonSwagger, "2.2.0"},
		{consts.ChoerodonGateWay, "2.2.0"},
		{consts.ChoerodonOauth, "2.2.0"},
		{consts.ChoerodonMonitor, "2.2.0"},
		{consts.ChoerodonFile, "2.2.0"},
		{consts.ChoerodonMessage, "2.2.0"},
		{consts.DevopsService, "2.2.0"},
		{consts.WorkflowService, "2.2.0"},
		{consts.GitlabService, "2.2.0"},
		{consts.AgileService, "2.2.0"},
		{consts.TestManagerService, "2.2.0"},
		{consts.KnowledgebaseService, "2.2.0"},
		{consts.ElasticsearchKb, "2.2.0"},
		{consts.ProdRepoService, "2.2.0"},
		{consts.CodeRepoService, "2.2.0"},
		{consts.ChoerodonFrontHzero, "2.2.0"},
		{consts.ChoerodonFront, "2.2.0"},
		{consts.ChoerodonClusterAgent, "2.2.0"},
		{consts.ChoerodonIamServiceBusiness, "2.2.0"},
		{consts.DevopsServiceBusiness, "2.2.0"},
		{consts.AgileServiceBusiness, "2.2.0"},
		{consts.DocRepoService, "2.2.0"},
		//{consts.HrdsQA, "2.0.0"},
		//{consts.MarketService, "2.0.0"},
		{consts.ChoerodonFrontBusiness, "2.2.0"},
		{consts.TestManagerServiceBusiness, "2.2.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "2.2")
		//fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
		//if VersionOrdinal(version) != VersionOrdinal(app.Version) {
		//	t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
		//	continue
		//}
		fmt.Printf("%s : %s\n", app.Name, version)

		//t.Logf("%s: %s", app.Name, version)
	}
}

func TestGetRelease20Tag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{consts.ChoerodonRegister, "2.0.0"},
		{consts.ChoerodonPlatform, "2.0.0"},
		{consts.ChoerodonAdmin, "2.0.0"},
		{consts.ChoerodonIam, "2.0.0"},
		{consts.ChoerodonAsgard, "2.0.0"},
		{consts.ChoerodonSwagger, "2.0.0"},
		{consts.ChoerodonGateWay, "2.0.0"},
		{consts.ChoerodonOauth, "2.0.0"},
		{consts.ChoerodonMonitor, "2.0.0"},
		{consts.ChoerodonFile, "2.0.0"},
		{consts.ChoerodonMessage, "2.0.0"},
		{consts.DevopsService, "2.0.0"},
		{consts.WorkflowService, "2.0.0"},
		{consts.GitlabService, "2.0.0"},
		{consts.AgileService, "2.0.0"},
		{consts.TestManagerService, "2.0.0"},
		{consts.KnowledgebaseService, "2.0.0"},
		{consts.ElasticsearchKb, "2.0.0"},
		{consts.ProdRepoService, "2.0.0"},
		{consts.CodeRepoService, "2.0.0"},
		{consts.ChoerodonFrontHzero, "2.0.0"},
		{consts.ChoerodonFront, "2.0.0"},
		{consts.ChoerodonClusterAgent, "2.0.0"},
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
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "2.0")
		//fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
		if VersionOrdinal(version) != VersionOrdinal(app.Version) {
			t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
			continue
		}
		//fmt.Printf("helm pull vista-c7n/%s --version %s\n", app.Name, version)

		t.Logf("%s: %s", app.Name, version)
	}
}

func TestGetRelease11Tag(t *testing.T) {
	apps := []struct {
		Name    string
		Version string
	}{
		{consts.ChoerodonRegister, "1.1.0"},
		{consts.ChoerodonPlatform, "1.1.0"},
		{consts.ChoerodonAdmin, "1.1.0"},
		{consts.ChoerodonIam, "1.1.0"},
		{consts.ChoerodonAsgard, "1.1.0"},
		{consts.ChoerodonSwagger, "1.1.0"},
		{consts.ChoerodonGateWay, "1.1.0"},
		{consts.ChoerodonOauth, "1.1.0"},
		{consts.ChoerodonMonitor, "1.1.0"},
		{consts.ChoerodonFile, "1.1.0"},
		{consts.ChoerodonMessage, "1.1.0"},
		{consts.DevopsService, "1.1.0"},
		{consts.WorkflowService, "1.1.0"},
		{consts.GitlabService, "1.1.0"},
		{consts.AgileService, "1.1.0"},
		{consts.TestManagerService, "1.0.0"},
		{consts.KnowledgebaseService, "1.0.0"},
		{consts.ElasticsearchKb, "1.1.0"},
		{consts.ProdRepoService, "1.1.0"},
		{consts.CodeRepoService, "1.1.0"},
		{consts.ChoerodonFrontHzero, "1.1.0"},
		{consts.ChoerodonFront, "1.1.0"},
		{consts.ChoerodonClusterAgent, "1.1.0"},
		//{consts.DevopsServiceBusiness, "1.1.0"},
		//{consts.HrdsQA, "2.0.0"},
		//{consts.MarketService, "2.0.0"},
		//{consts.ChoerodonFrontBusiness, "1.1.0"},
		//{consts.TestManagerServiceBusiness, "1.1.0"},
	}
	for _, app := range apps {
		version, _ := GetReleaseTag(consts.DefaultRepoUrl, app.Name, "1.1")
		//fmt.Printf("    [\"registry.cn-shanghai.aliyuncs.com/c7n/%s\"]=\"%s\"\n", app.Name, version)
		//if VersionOrdinal(version) != VersionOrdinal(app.Version) {
		//	t.Errorf("%s  %s is not newest version, newer version is %s", app.Name, app.Version, version)
		//	continue
		//}
		fmt.Printf("%s: %s\n", app.Name, version)

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
