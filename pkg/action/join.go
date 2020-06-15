package action

import (
	"github.com/choerodon/c7nctl/pkg/client"
	c7n_utils "github.com/choerodon/c7nctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	addMasterPlaybook = "91-add-master.yml"
	addWorkerPlaybook = "92-add-worker.yml"
	addEtcdPlaybook   = "93-add-etcd.yml"
)

func (i *InstallK8s) RunJoinNode() {
	inventory := client.Inventory{}
	in, err := ioutil.ReadFile(hostPath)
	c7n_utils.CheckErrAndExit(err, 1)
	err = yaml.Unmarshal(in, &inventory)
	c7n_utils.CheckErr(err)
	i.checkJoinIPs(&inventory)

	if len(i.MasterIPs) > 0 {
		joinMaster(i.MasterIPs, &inventory)
		joinNode(i.MasterIPs, &inventory)
	}
	if len(i.NodeIPs) > 0 {
		joinNode(i.NodeIPs, &inventory)
	}
}

func joinNode(ps []string, inventory *client.Inventory) {
	log.WithField("IPs", ps).Info("Start join worker node into k8s")
	inventory.AddHosts(ps)
	inventory.AddNewWorkers(ps)
	writeHosts(inventory)
	execAnsiblePlaybook(addWorkerPlaybook)
	inventory.MoveToKubeWorkers()
	writeHosts(inventory)
}

func joinMaster(ps []string, inventory *client.Inventory) {
	log.WithField("IPs", ps).Info("Start join master node into k8s")
	inventory.AddHosts(ps)
	inventory.AddNewMasters(ps)
	writeHosts(inventory)

	execAnsiblePlaybook(addMasterPlaybook)
	inventory.MoveToKubeMasters()
	writeHosts(inventory)

	for _, ip := range ps {
		log.WithField("IP", ip).Info("Start join etcd node into k8s")
		inventory.AddNewEtcd(ip)
		writeHosts(inventory)

		execAnsiblePlaybook(addEtcdPlaybook)
		inventory.MoveToEtcd()
		writeHosts(inventory)
	}
}

func writeHosts(inventory *client.Inventory) {
	i, err := yaml.Marshal(inventory)
	c7n_utils.CheckErr(err)
	ioutil.WriteFile(hostPath, i, os.ModePerm)
}

func (i *InstallK8s) checkJoinIPs(inventory *client.Inventory) {
	oldIPs := inventory.GetHosts()
	var masters []string
	masters = append(masters, i.MasterIPs...)
	masters = append(masters, inventory.GetKubeMaster()...)
	checkMasterIP(masters)

	var hosts []string
	hosts = append(hosts, oldIPs...)
	hosts = append(hosts, i.MasterIPs...)
	hosts = append(hosts, i.NodeIPs...)
	checkIP(hosts)
}
