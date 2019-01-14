package c7nclient

import "time"

// 設定型

type C7NPlatformContext struct {
	Token string `yaml:"token"`
	Timeout time.Duration `yaml:"timeout"`
	Server string    `yaml:"server"`
	Organization string `yaml:"organization"`
	Project  string  `yaml:"project"`
	ProjectId int `yaml:"projectId"`
	OrganizationId int `yaml:"organizationId"`
	OrganizationCode string `yaml:"organizationCode"`
	ProjectCode  string  `yaml:"projectCode"`
	Name  string
}
type C7NPlatformConfig struct {
	Context C7NPlatformContext `yaml:"context"`
	Name  string `yaml:"name"`
}

type C7NConfig struct {
	Contexts  []C7NPlatformConfig `yaml:"contexts"`
	CurrentContext  string `yaml:"current-context"`
}


