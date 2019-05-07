package c7nclient

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/c7nclient/model"
	"github.com/pkg/errors"
	"io"
	"strconv"
	"time"
)

const baseFormat = "2006-01-01 00:00:00"

func (c *C7NClient) ListClusters(out io.Writer, organizationId int) {
	if organizationId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["page"] = "0"
	paras["size"] = "10"
	req, err := c.newRequest("POST", fmt.Sprintf("devops/v1/organizations/%d/clusters/page_cluster", organizationId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var clusters = model.Clusters{}
	_, err = c.do(req, &clusters)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	clusterInfos := []model.ClusterInfo{}
	for _, cluster := range clusters.Content {
		var status string
		if cluster.Connect {
			status = "已连接"
		} else {
			status = "未连接"
		}
		clusterInfo := model.ClusterInfo{
			Name:        cluster.Name,
			Code:        cluster.Code,
			Description: cluster.Description,
			Status:      status,
		}
		clusterInfos = append(clusterInfos, clusterInfo)
	}
	model.PrintClusterInfo(clusterInfos, out)

}

func (c *C7NClient) GetCluster(out io.Writer, organizationId int, clusterCode string) (error error, result model.Cluster) {
	if organizationId == 0 {
		return errors.New("the organization is not found"), model.Cluster{}
	}
	paras := make(map[string]interface{})
	paras["code"] = clusterCode
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/organizations/%d/clusters/query_by_code", organizationId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")
	}
	var cluster = model.Cluster{}
	_, err = c.do(req, &cluster)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return err, cluster
	}
	return nil, cluster
}

func (c *C7NClient) ListClusterNode(out io.Writer, organizationId int, clusterId int) {
	if organizationId == 0 {
		return
	}
	paras := make(map[string]interface{})
	paras["cluster_id"] = strconv.Itoa(clusterId)
	paras["page"] = "0"
	paras["size"] = "10"
	req, err := c.newRequest("GET", fmt.Sprintf("devops/v1/organizations/%d/clusters/page_nodes", organizationId), paras, nil)
	if err != nil {
		fmt.Printf("build request error")

	}
	var nodes = model.Nodes{}
	_, err = c.do(req, &nodes)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	nodeInfos := []model.NodeInfo{}
	for _, node := range nodes.Content {
		now := time.Now()
		creationTime, _ := time.Parse(baseFormat, node.CreateTime)
		nodeInfo := model.NodeInfo{
			Status:        node.Status,
			NodeName:      node.NodeName,
			NodeType:      node.Type,
			CpuRequest:    node.CPURequest,
			CpuLimit:      node.CPULimit,
			MemoryRequest: node.MemoryRequest,
			MemoryLimit:   node.MemoryLimit,
			CreationDate:  c.getTime(now.Sub(creationTime).Seconds()),
		}
		nodeInfos = append(nodeInfos, nodeInfo)
	}
	model.PrintNodeInfo(nodeInfos, out)

}

func (c *C7NClient) CreateCluster(out io.Writer, organizationId int, clusterPostInfo *model.ClusterPostInfo) {
	req, err := c.newRequest("POST", fmt.Sprintf("/devops/v1/organizations/%d/clusters", organizationId), nil, clusterPostInfo)
	if err != nil {
		fmt.Printf("build request error")
	}
	var clusterInfo string
	_, err = c.doHandleString(req, &clusterInfo)
	if err != nil {
		fmt.Printf("request err:%v", err)
		return
	}
	fmt.Println(clusterInfo)

}
