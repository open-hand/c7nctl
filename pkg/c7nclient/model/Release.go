package model

type Release struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec `yaml:"spec"`
}

type Metadata struct {
	Name string  `yaml:"name"`
}

type Spec struct {
	ChartName    string `yaml:"chartName"`
	ChartVersion string `yaml:"chartVersion"`
	RepoUrl      string `yaml:"repoUrl"`
	Values       string `yaml:"values"`
}
