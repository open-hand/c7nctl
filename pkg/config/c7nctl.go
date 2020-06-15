package config

var Cfg Config

// c7nctl 默认的配置项，包括安装的基本信息和连接c7n的信息
type Config struct {
	Version string
	// 安装 c7n 时，用户需要输入的邮箱信息
	Terms   Terms
	OpsMail string `yaml:"opsMail"`

	// 暂时不知道有什么用处
	Clusters       []*NamedCluster
	CurrentCluster string `yaml:"current-cluster"`
	Users          []*NamedUser

	// 连接 c7n 的上下文信息
	Contexts       []C7NContext `yaml:"contexts"`
	CurrentContext string       `yaml:"currentContext"`
}

type Terms struct {
	Accepted bool
}

type NamedUser struct {
	Name string `yaml:"name"`
	User *User  `yaml:"user"`
}

type NamedCluster struct {
	Name    string   `yaml:"name"`
	Cluster *Cluster `yaml:"cluster"`
}

type User struct {
	Mail       string
	Token      string
	Name       string
	ValidUntil int64
}

type Cluster struct {
	Server       string
	User         *User  `yaml:"-"`
	SelectedUser string `yaml:"user"`
}

// C7n 的连接上下文
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
