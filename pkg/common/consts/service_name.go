package consts

// 服务列表
const (
	ChartMuseum          = "chartmuseum"
	Redis                = "c7n-redis"
	Mysql                = "c7n-mysql"
	Gitlab               = "gitlab"
	Harbor               = "harbor"
	Sonarqube            = "sonarqube"
	ChoerodonRegister    = "choerodon-register"
	ChoerodonPlatform    = "choerodon-platform"
	ChoerodonAdmin       = "choerodon-admin"
	ChoerodonIam         = "choerodon-iam"
	ChoerodonOauth       = "choerodon-oauth"
	ChoerodonGateWay     = "choerodon-gateway"
	ChoerodonAsgard      = "choerodon-asgard"
	ChoerodonSwagger     = "choerodon-swagger"
	ChoerodonMessage     = "choerodon-message"
	ChoerodonMonitor     = "choerodon-monitor"
	ChoerodonFile        = "choerodon-file"
	DevopsService        = "devops-service"
	GitlabService        = "gitlab-service"
	WorkflowService      = "workflow-service"
	AgileService         = "agile-service"
	TestManagerService   = "test-manager-service"
	KnowledgebaseService = "knowledgebase-service"
	ElasticsearchKb      = "elasticsearch-kb"
	ProdRepoService      = "prod-repo-service"
	CodeRepoService      = "code-repo-service"
	ChoerodonFrontHzero  = "choerodon-front-hzero"
	ChoerodonFront       = "choerodon-front"

	ChoerodonClusterAgent       = "choerodon-cluster-agent"
	ChoerodonIamServiceBusiness = "choerodon-iam-service-business"
	DevopsServiceBusiness       = "devops-service-business"
	AgileServiceBusiness        = "agile-service-business"
	DocRepoService              = "doc-repo-service"
	// HrdsQA                      = "hrds-qa"
	// MarketService               = "market-service"
	TestManagerServiceBusiness = "test-manager-service-business"
	ChoerodonFrontBusiness     = "choerodon-front-business"
)

var ServerListBiz = []string{
	ChoerodonRegister,
	ChoerodonPlatform,
	ChoerodonAdmin,
	ChoerodonIamServiceBusiness,
	ChoerodonMessage,
	ChoerodonOauth,
	ChoerodonGateWay,
	ChoerodonAsgard,
	ChoerodonSwagger,
	ChoerodonMonitor,
	ChoerodonFile,
	DevopsServiceBusiness,
	GitlabService,
	WorkflowService,
	AgileServiceBusiness,
	TestManagerService,
	ElasticsearchKb,
	KnowledgebaseService,
	ProdRepoService,
	CodeRepoService,
	DocRepoService,
	ChoerodonFrontHzero,
	ChoerodonFrontBusiness,
}
