package config

type ChoerodonVersion struct {
	Name string `yaml:"name"`

	// 配置文件的版本
	Version string  `yaml:"version"`
	Spec    PkgSpec `yaml:"spec"`
}

type PkgSpec struct {
	// 需要生成离线安装包的猪齿鱼版本
	VersionRegexp string `yaml:"version-regexp"`
	Offline       bool   `yaml:"offline"`
	Chart         Chart  `yaml:"chart"`
	Image         Image  `yaml:"image"`
}

type Chart struct {
	DefaultSource ChartRepository  `yaml:"default-source"`
	DefaultTarget ChartRepository  `yaml:"default-target"`
	Component     []ChartComponent `yaml:"component"`
}

type ChartComponent struct {
	Name     string          `yaml:"name"`
	Version  string          `yaml:"version"`
	Source   ChartRepository `yaml:"source"`
	Target   ChartRepository `yaml:"target"`
	Category string          `yaml:"category"`
}

type ChartRepository struct {
	Url      string `yaml:"url"`
	Repo     string `yaml:"repo"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Image struct {
	Registry []ImageRegistry `yaml:"registry"`
	Images   []string        `yaml:"images"`
}

type ImageRegistry struct {
	Domain     string `yaml:"domain"`
	Repository string `yaml:"repository"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Insecure   string `yaml:"insecure"`
}
