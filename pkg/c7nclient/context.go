package c7nclient

// 設定型

type C7NContext struct {
	Name   string  `yaml:"name"`
	Server string  `yaml:"server"`
	User   C7NUser `yaml:"user"`
}

type C7NUser struct {
	UserName         string `yaml:"userName"`
	Token            string `yaml:"token"`
	ProjectId        int    `yaml:"projectId"`
	OrganizationId   int    `yaml:"organizationId"`
	OrganizationCode string `yaml:"organizationCode"`
	ProjectCode      string `yaml:"projectCode"`
}

type C7NConfig struct {
	Contexts       []C7NContext `yaml:"contexts"`
	CurrentContext string       `yaml:"currentContext"`
}
