package client

var DefaultAnsibleVar = Vars{
	SkipVerifyNode:        false,
	KubeVersion:           "1.16.9",
	LbMode:                "openresty",
	LbKubeApiserverPort:   8443,
	KubePodSubnet:         "10.244.0.0/18",
	KubeServiceSubnet:     "10.244.64.0/18",
	KubeNetworkNodePrefix: 24,
	KubeMaxPods:           110,
	NetworkPlugin:         "calico",
	KubeletRootDir:        "/var/lib/kubelet",
	DockerStorageDir:      "/var/lib/docker",
	EtcdDataDir:           "/var/lib/etcd",
}

type Inventory struct {
	All Group `yaml:"all"`
}

type Group struct {
	Hosts    map[string]Host `yaml:"hosts,omitempty"`
	Vars     Vars            `yaml:"vars,omitempty"`
	Children Children        `yaml:"children,omitempty"`
}

type Host struct {
	AnsibleHost     string `yaml:"ansible_host,omitempty"`
	AnsiblePort     int    `yaml:"ansible_port,omitempty"`
	AnsibleUser     string `yaml:"ansible_user,omitempty"`
	AnsiblePassword string `yaml:"ansible_ssh_pass,omitempty"`
}

type Children struct {
	Lb         Hostname `yaml:"lb"`
	Etcd       Hostname `yaml:"etcd"`
	KubeMaster Hostname `yaml:"kube-master"`
	KubeWorker Hostname `yaml:"kube-worker"`
	NewMaster  Hostname `yaml:"new-master"`
	NewWorker  Hostname `yaml:"new-worker"`
	NewEtcd    Hostname `yaml:"new-etcd"`
}

type Hostname struct {
	Hosts map[string]interface{} `yaml:"hosts"`
}

type Vars struct {
	SkipVerifyNode        bool   `yaml:"skip_verify_node,omitempty"`
	KubeVersion           string `yaml:"kube_version,omitempty"`
	LbMode                string `yaml:"lb_mode,omitempty"`
	LbKubeApiserverPort   int    `yaml:"lb_kube_apiserver_port,omitempty"`
	KubePodSubnet         string `yaml:"kube_pod_subnet,omitempty"`
	KubeServiceSubnet     string `yaml:"kube_service_subnet,omitempty"`
	KubeNetworkNodePrefix int    `yaml:"kube_network_node_prefix,omitempty"`
	KubeMaxPods           int    `yaml:"kube_max_pods,omitempty"`
	NetworkPlugin         string `yaml:"network_plugin,omitempty"`
	KubeletRootDir        string `yaml:"kubelet_root_dir,omitempty"`
	DockerStorageDir      string `yaml:"docker_storage_dir,omitempty"`
	EtcdDataDir           string `yaml:"etcd_data_dir,omitempty"`
}
