package client

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

const testHost = `all:
  hosts:
    local-node1:
      ansible_host: 192.168.56.201
      ansible_port: 22
      ansible_user: root
      ansible_ssh_pass: yishuida
  children:
    lb:
      hosts:
    etcd:
      hosts:
        local-node1:
    kube-master:
      hosts:
        local-node1:
    kube-worker:
      hosts:
        local-node1:
    new-master:
      hosts:
    new-worker:
      hosts:
    new-etcd:
      hosts:
  vars:
    skip_verify_node: true
    lb_mode: openresty
    lb_kube_apiserver_port: 8443
    kube_pod_subnet: 10.244.0.0/18
    kube_service_subnet: 10.244.64.0/18
    kube_network_node_prefix: 24
    kube_max_pods: 110
    network_plugin: calico
    kubelet_root_dir: /var/lib/kubelet
    docker_storage_dir: /var/lib/docker
    etcd_data_dir: /var/lib/etcd
`

func TestAnsible_inventory(t *testing.T) {
	inventory := Inventory{
		All: Group{
			Hosts: map[string]Host{
				"local-node1": {
					AnsibleHost:     "172.24.33.193",
					AnsiblePort:     22,
					AnsibleUser:     "root",
					AnsiblePassword: "yishuidaA!",
				},
				"local-node2": {
					AnsibleHost:     "172.24.33.194",
					AnsiblePort:     22,
					AnsibleUser:     "root",
					AnsiblePassword: "yishuidaA!",
				},
				"local-node3": {
					AnsibleHost:     "172.24.33.195",
					AnsiblePort:     22,
					AnsibleUser:     "root",
					AnsiblePassword: "yishuidaA!",
				},
			},
			Vars: Vars{
				SkipVerifyNode:        true,
				KubeVersion:           "1.15.7",
				LbMode:                "openresty",
				NetworkPlugin:         "calico",
				LbKubeApiserverPort:   8443,
				KubePodSubnet:         "10.244.0.0/18",
				KubeServiceSubnet:     "10.244.64.0/18",
				KubeNetworkNodePrefix: 24,
				KubeMaxPods:           110,
				KubeletRootDir:        "/var/lib/kubelet",
				DockerStorageDir:      "/var/lib/docker",
				EtcdDataDir:           "/var/lib/etcd",
			},
			Children: Children{
				Etcd: Hostname{
					Hosts: map[string]interface{}{
						"local-node1": struct{}{},
						"local-node2": struct{}{},
						"local-node3": struct{}{},
					},
				},
				KubeMaster: Hostname{
					Hosts: map[string]interface{}{
						"local-node1": struct{}{},
						"local-node2": struct{}{},
						"local-node3": struct{}{},
					},
				},
				KubeWorker: Hostname{
					Hosts: map[string]interface{}{
						"local-node1": struct{}{},
						"local-node2": struct{}{},
						"local-node3": struct{}{},
					},
				},
			},
		},
	}
	i, _ := yaml.Marshal(inventory)
	_ = ioutil.WriteFile("/tmp/host.yaml", i, 0766)
	t.Log(inventory)

	_ = yaml.Unmarshal([]byte(testHost), &inventory)
	t.Log(inventory)
}
