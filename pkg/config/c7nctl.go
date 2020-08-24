package config

// c7nctl 默认的配置项，包括安装的基本信息和连接c7n的信息
type Config struct {
	Version string `yaml:"version"`
	// 安装 c7n 时，用户需要输入的邮箱信息
	Terms   Terms  `yaml:"terms"`
	OpsMail string `yaml:"opsMail"`

	// 暂时不知道有什么用处
	Clusters       []*NamedCluster `yaml:"clusters"`
	CurrentCluster string          `yaml:"current-cluster"`
	Users          []*NamedUser    `yaml:"users"`

	// 连接 c7n 的上下文信息
	Contexts       []C7NContext `yaml:"contexts"`
	CurrentContext string       `yaml:"currentContext"`
}

type Terms struct {
	Accepted bool `yaml:"accepted"`
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
	Mail       string `yaml:"mail"`
	Token      string `yaml:"token"`
	Name       string `yaml:"name"`
	ValidUntil int64  `yaml:"valid-until"`
}

type Cluster struct {
	Server       string `yaml:"server"`
	User         *User  `yaml:"user"`
	SelectedUser string `yaml:"selected-user"`
}

// C7n 的连接上下文
type C7NContext struct {
	Name   string  `yaml:"name"`
	Server string  `yaml:"server"`
	User   C7NUser `yaml:"user"`
}

type C7NUser struct {
	UserName         string `yaml:"user-name"`
	Token            string `yaml:"token"`
	ProjectId        int    `yaml:"project-id"`
	OrganizationId   int    `yaml:"organization-id"`
	OrganizationCode string `yaml:"organization-code"`
	ProjectCode      string `yaml:"project-code"`
}
