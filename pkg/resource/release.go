package resource

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/consts"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/choerodon/c7nctl/pkg/utils"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vinkdong/gox/http/downloader"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func (rls *Release) ValuesRaw() (string, error) {
	data, err := rls.getUserValuesTpl()

	if err != nil {
		return "", std_errors.WithMessage(err, fmt.Sprintf("get user values for %s failed", rls.Name))
	}
	if len(data) == 0 {
		url := fmt.Sprintf(consts.RemoteInstallResourceRootUrl, context.Ctx.Version, "values/"+rls.Name+".yaml")
		nData, statusCode, err := downloader.GetFileContent(url)
		if statusCode == 200 && err == nil {
			data = nData
		}
	}
	if len(data) > 0 {
		//return rls.renderValue(string(data[:]))
		return utils.RenderReleaseValue(rls.Name, string(data[:])), nil
	}
	return "", nil
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
	if ji.Status == context.SucceedStatus {
		log.Infof("Release %s is already installed", rls.Name)
		return nil
	}
	if ji.Status == context.RenderedStatus || ji.Status == context.FailedStatus {
		// 等待依赖项安装完成
		for _, r := range rls.Requirements {
			checkReleasePodRunning(r)
		}

		if err := rls.ExecutePreCommands(); err != nil {
			return std_errors.WithMessage(err, fmt.Sprintf("Release %s execute pre commands failed", rls.Name))
		}

		values := rls.HelmValues()
		var releaseName string
		if context.Ctx.Prefix != "" {
			releaseName = fmt.Sprintf("%s-%s", context.Ctx.Prefix, rls.Name)
		} else {
			releaseName = rls.Name

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
			log.WithField("release", rls.Name).Debug(k)
		}
		// TODO useless
		if rls.Timeout > 0 {
			values = append(values, fmt.Sprintf("preJob.timeout=%d", rls.Timeout))
		}
		raw, err := rls.ValuesRaw()
		if err != nil {
			return std_errors.WithMessage(err, "Release %s get value failed")
		}
		err = context.Ctx.HelmClient.InstallRelease(values, raw, chartArgs)

		if err != nil {
			ji.Status = context.FailedStatus
			context.Ctx.UpdateJobInfo(ji)
			return err
		} else {
			ji.Status = context.InstalledStatus
			context.Ctx.UpdateJobInfo(ji)
		}
	}
	if ji.Status == context.InstalledStatus {
		// TODO 解决 after task失败二次执行
		if len(rls.AfterInstall) > 0 {
			go rls.ExecuteAfterTasks()
			// return std_errors.WithMessage(err, "Execute after task failed")
		}
		ji.Status = context.SucceedStatus
		// 更新 jobinfo 状态
		context.Ctx.UpdateJobInfo(ji)
	}

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

// 基础组件——比如 gitlab-ha ——有 app 标签，c7n 有 choerodon.io/release 标签
// TODO 所有组件设置统一的label
func checkReleasePodRunning(rls string) {

	clientset := *context.Ctx.KubeClient
	namespace := context.Ctx.Namespace
	time.Sleep(time.Second * 3)

	labels := []string{
		fmt.Sprintf("choerodon.io/release=%s", rls),
		fmt.Sprintf("app=%s", rls),
	}
	for {
		for _, label := range labels {
			deploy, err := clientset.AppsV1().Deployments(namespace).List(meta_v1.ListOptions{LabelSelector: label})
			if errors.IsNotFound(err) {
				log.Infof("Deployment %s in namespace %s not found\n", label, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				log.Infof("Error getting deployment %s in namespace %s: %v\n",
					label, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				for _, d := range deploy.Items {
					if *d.Spec.Replicas != d.Status.ReadyReplicas {
						log.Infof("Deployment %s is not ready\n", d.Name)
					} else {
						log.Infof("Deployment %s is Ready\n", d.Name)
						return
					}
				}
			}
			ss, err := clientset.AppsV1().StatefulSets(namespace).List(meta_v1.ListOptions{LabelSelector: label})
			if errors.IsNotFound(err) {
				log.Infof("StatefulSet %s in namespace %s not found\n", label, namespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				log.Infof("Error getting statefulSet %s in namespace %s: %v\n",
					label, namespace, statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				for _, s := range ss.Items {
					if *s.Spec.Replicas != s.Status.ReadyReplicas {
						log.Infof("StatefulSet %s is not ready\n", s.Name)
					} else {
						log.Infof("statefulSet %s is Ready\n", s.Name)
						return
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}
