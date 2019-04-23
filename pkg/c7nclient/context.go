package c7nclient

// 設定型

type C7NPlatformContext struct {
	Token            string `yaml:"token"`
	Server           string `yaml:"server"`
	ProjectId        int    `yaml:"projectId"`
	OrganizationId   int    `yaml:"organizationId"`
	OrganizationCode string `yaml:"organizationCode"`
	ProjectCode      string `yaml:"projectCode"`
}
