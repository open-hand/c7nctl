package model

type Release struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	ChartName    string `json:"chartName"`
	ChartVersion string `json:"chartVersion"`
	RepoUrl      string `json:"repoUrl"`
	Values       string `json:"values"`
}
