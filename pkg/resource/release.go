package resource

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/consts"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/choerodon/c7nctl/pkg/utils"
	c7n_utils "github.com/choerodon/c7nctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/vinkdong/gox/http/downloader"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"time"
)

type Release struct {
	Name         string
	Chart        string
	Version      string
	Namespace    string
	RepoURL      string
	Values       []context.ChartValue
	Persistence  []*Persistence
	PreInstall   []ReleaseJob
	AfterInstall []ReleaseJob
	Requirements []string
	Resource     *config.Resource
	// TODO Remove
	Timeout     int
	Prefix      string
	SkipInput   bool
	PaaSVersion string
}

type ReleaseJob struct {
	Name     string
	InfraRef string `yaml:"infraRef"`
	Database string `yaml:"database"`
	Commands []string
	Mysql    []string
	Psql     []string `yaml:"psql"`
	Opens    []string
	Request  *Request
}

type Request struct {
	Header     []ChartValue
	Url        string
	Parameters []ChartValue
	Body       string
	Method     string
}

type ChartValue struct {
	Name  string
	Value string
	Input context.Input
	Case  string
	Check string
}

func (r *Request) parserParams() string {
	var params []string
	for _, p := range r.Parameters {
		params = append(params, fmt.Sprintf("%s=%s", p.Name, p.Value))
	}
	return strings.Join(params, "&")
}

func (r *Request) parserUrl() string {
	params := r.parserParams()
	url := r.Url
	if params != "" {
		url = fmt.Sprintf("%s?%s", url, params)
	}
	return url
}

func (pi *ReleaseJob) ExecuteSql(rls *Release, sqlType string) error {
	_, r := context.Ctx.GetJobInfo(pi.Name)
	if r != nil && r.Type == sqlType && r.Status == context.SucceedStatus {
		log.Infof("task %s had executed", pi.Name)
		return nil
	}
	log.Infof("executing %s , %s", rls.Name, pi.Name)

	r = &context.JobInfo{
		Name:     pi.Name,
		RefName:  rls.Name,
		Type:     context.TaskType,
		Status:   context.SucceedStatus,
		TaskType: context.SqlTask,
		Version:  rls.Version,
	}

	defer context.Ctx.AddJobInfo(r)

	sqlList := make([]string, 0)

	for _, v := range pi.Commands {
		sqlList = append(sqlList, v)
	}
	for _, v := range pi.Mysql {
		sqlList = append(sqlList, v)
	}
	for _, v := range pi.Psql {
		sqlList = append(sqlList, v)
	}
	res := utils.GetResource(pi.InfraRef)
	s := context.Ctx.Slaver
	if err := s.ExecuteRemoteSql(sqlList, res, pi.Database, sqlType); err != nil {
		r.Status = context.FailedStatus
		r.Reason = err.Error()
		return err
	}
	return nil
}

func (pi *ReleaseJob) ExecuteRequests(infra *Release) error {
	if pi.Request == nil {
		return nil
	}
	_, r := context.Ctx.GetJobInfo(pi.Name)
	if r != nil && r.Type == context.HttpGetTask {
		log.Infof("task %s had executed", pi.Name)
		return nil
	}

	r = &context.JobInfo{
		Name:     pi.Name,
		RefName:  infra.Name,
		Type:     context.TaskType,
		Status:   context.SucceedStatus,
		TaskType: context.HttpGetTask,
		Version:  infra.Version,
	}

	defer context.Ctx.AddJobInfo(r)

	req := pi.Request
	s := context.Ctx.Slaver
	header := make(map[string][]string)
	for _, h := range req.Header {
		header[h.Name] = []string{h.Value}
	}

	reqUrl := req.Url
	paramsString := req.parserParams()
	if paramsString != "" {
		reqUrl = reqUrl + "?" + paramsString
	}
	f := slaver.Forward{
		Url:    reqUrl,
		Body:   req.Body,
		Header: header,
		Method: req.Method,
	}

	_, err := s.ExecuteRemoteRequest(f)
	if err != nil {
		r.Status = context.FailedStatus
		r.Reason = err.Error()
	}
	return err
}

func (rls *Release) String() string {
	return rls.Name
}
func (rls *Release) ExecutePreCommands() error {
	err := rls.executeExternalFunc(rls.PreInstall)
	return err
}

func (rls *Release) executeExternalFunc(c []ReleaseJob) error {
	for _, pi := range c {
		if len(pi.Commands) > 0 {
			if err := pi.ExecuteSql(rls, "mysql"); err != nil {
				return err
			}
		}
		if len(pi.Mysql) > 0 {
			if err := pi.ExecuteSql(rls, "mysql"); err != nil {
				return err
			}
		}
		if len(pi.Psql) > 0 {
			if err := pi.ExecuteSql(rls, "postgres"); err != nil {
				return err
			}
		}
		if pi.Request != nil {
			if err := pi.ExecuteRequests(rls); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rls *Release) getUserValuesTpl() ([]byte, error) {
	return context.Ctx.UserConfig.GetHelmValuesTpl(rls.Name)
}

// convert yml format values template to yaml raw data
func (rls *Release) ValuesRaw() string {
	data, err := rls.getUserValuesTpl()

	if err != nil {
		log.Error(err)
		log.Errorf("get user values for %s failed", rls.Name)
		os.Exit(127)
	}
	if len(data) == 0 {
		url := fmt.Sprintf(consts.RemoteInstallResourceRootUrl, rls.PaaSVersion, "values/"+rls.Name+".yaml")
		nData, statusCode, err := downloader.GetFileContent(url)
		if statusCode == 200 && err == nil {
			data = nData
		}
	}
	if len(data) > 0 {
		//return rls.renderValue(string(data[:]))
		return utils.RenderReleaseValue(rls.Name, string(data[:]))
	}
	return ""
}

// convert yml values to values list as xxx=yyy
func (rls *Release) HelmValues() []string {
	values := make([]string, len(rls.Values))
	// store values for feature use
	for k, v := range rls.Values {
		// 解决特殊字符
		values[k] = fmt.Sprintf("%s=%s", v.Name, v.Value)
	}
	return values
}

// resource infra
func (rls *Release) Install() error {
	_, ji := context.Ctx.GetJobInfo(rls.Name)
	// 如果已经安装直接返回
	if ji.Status == context.SucceedStatus {
		log.Infof("already installed %s", rls.Name)
		return nil
	}

	for _, r := range rls.Requirements {
		checkReleasePodRunning(r)
	}

	err := rls.ExecutePreCommands()
	c7n_utils.CheckErr(err)

	values := rls.HelmValues()
	releaseName := rls.Name
	if context.Ctx.Prefix != "" {
		releaseName = fmt.Sprintf("%s-%s", context.Ctx.Prefix, rls.Name)
	}
	chartArgs := client.ChartArgs{
		ReleaseName: releaseName,
		Namespace:   context.Ctx.Namespace,
		RepoUrl:     context.Ctx.RepoUrl,
		Verify:      false,
		Version:     rls.Version,
		ChartName:   rls.Chart,
	}

	log.Infof("installing %s", rls.Name)
	for _, k := range values {
		log.Debugf(k)
	}
	if rls.Timeout > 0 {
		values = append(values, fmt.Sprintf("preJob.timeout=%d", rls.Timeout))
	}
	raw := rls.ValuesRaw()
	err = context.Ctx.HelmClient.InstallRelease(values, raw, chartArgs)

	if err != nil {
		ji.Status = context.FailedStatus
		context.Ctx.UpdateJobInfo(ji)
		return err
	}

	if len(rls.AfterInstall) > 0 {
		err = rls.ExecuteAfterTasks()
		c7n_utils.CheckErr(err)
	}
	ji.Status = context.SucceedStatus
	// 更新 jobinfo 状态
	context.Ctx.UpdateJobInfo(ji)
	return nil
}

func (rls *Release) InstallComponent() error {

	values := rls.HelmValues()
	releaseName := rls.Name
	if rls.Prefix != "" {
		releaseName = fmt.Sprintf("%s-%s", context.Ctx.Prefix, rls.Name)
	}
	chartArgs := client.ChartArgs{
		ReleaseName: releaseName,
		Namespace:   context.Ctx.Namespace,
		RepoUrl:     context.Ctx.RepoUrl,
		Verify:      false,
		Version:     rls.Version,
		ChartName:   rls.Chart,
	}

	log.Infof("installing %s", rls.Name)
	for _, k := range values {
		log.Debug(k)
	}
	if rls.Timeout > 0 {
		values = append(values, fmt.Sprintf("preJob.timeout=%d", rls.Timeout))
	}
	// raw := rls.ValuesRaw()
	err := context.Ctx.HelmClient.InstallRelease(values, "", chartArgs)
	return err
}

//
func (rls *Release) ExecuteAfterTasks() error {
	checkReleasePodRunning(rls.Name)

	log.Infof("%s: started, will execute required commands and requests", rls.Name)
	err := rls.executeExternalFunc(rls.AfterInstall)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (rls *Release) CatchInitJobs() error {
	client := *context.Ctx.KubeClient
	jobInterface := client.BatchV1().Jobs(context.Ctx.UserConfig.Metadata.Namespace)
	jobList, err := jobInterface.List(meta_v1.ListOptions{
		LabelSelector: fmt.Sprintf("choerodon.io/release=%s", rls.Name),
	})
	if err != nil {
		return err
	}
	for _, job := range jobList.Items {
		if job.Status.Active > 0 {
			log.Infof("job %s haven't finished yet. please wait patiently", job.Name)
			jobLabelSelector := fmt.Sprintf("job-name=%s", job.Name)
			rls.catchPodLogs(jobLabelSelector, client)
		}
	}
	log.Debugf("still resource %s", rls.Name)
	return nil
}

func (rls *Release) catchPodLogs(labelSelector string, client kubernetes.Interface) error {
	podList, err := client.CoreV1().Pods(rls.Namespace).List(meta_v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return err
	}

	for _, po := range podList.Items {
		if po.Status.Phase == core_v1.PodRunning {
			log.Debugf("you can watch logs by execute follow command:\nkubectl logs -f %s -n %s",
				po.Name, po.Namespace)
		}
	}
	return nil
}

// 基础组件——比如 gitlab-ha ——有 app 标签，c7n 有 choerodon.io/release 标签
// TODO 所有组件设置统一的label
func checkReleasePodRunning(rls string) {

	clientset := *context.Ctx.KubeClient
	namespace := context.Ctx.Namespace
	time.Sleep(time.Second * 3)

	for {
		pods1, err := clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{LabelSelector: fmt.Sprintf("choerodon.io/release=%s", rls)})
		if err != nil {
			log.Error(err)
		}
		pods2, err := clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", rls)})
		if err != nil {
			log.Error(err)
		}

		pods := append(pods1.Items, pods2.Items...)
		log.Infof("%s has %d pods in the cluster\n", rls, len(pods))
		isReady := true
		for _, p := range pods {
			if p.Status.Phase != core_v1.PodRunning {
				log.Infof("Waiting for pods %s of release %s", p.Name, rls)
				isReady = false
			}
		}
		if isReady {
			log.Infof("%s's pods is running", rls)
			break
		} else {
			time.Sleep(time.Second * 4)
		}
	}
}
