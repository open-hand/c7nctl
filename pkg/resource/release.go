package resource

import (
	"context"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/consts"
	c7nctx "github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/slaver"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vinkdong/gox/http/downloader"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"sync"
	"time"
)

type Release struct {
	Name         string
	Chart        string
	Version      string
	Namespace    string
	RepoURL      string
	Values       []c7nctx.ChartValue
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
	Input c7nctx.Input
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

func (pi *ReleaseJob) ExecuteSql(rls *Release, sqlType string, s *slaver.Slaver) error {
	task, err := c7nctx.GetTaskFromCM(rls.Namespace, pi.Name)
	if err != nil {
		return err
	}

	if task != nil && task.Status == c7nctx.SucceedStatus {
		log.Infof("task %s had executed", pi.Name)
		return nil
	}
	log.Infof("executing %s , %s", rls.Name, pi.Name)

	task = &c7nctx.TaskInfo{
		Name:     pi.Name,
		RefName:  rls.Name,
		Type:     c7nctx.TaskType,
		Status:   c7nctx.SucceedStatus,
		TaskType: c7nctx.SqlTask,
		Version:  rls.Version,
	}
	defer c7nctx.AddTaskToCM(rls.Namespace, *task)

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
	rlsRef, _ := c7nctx.GetReleaseTaskInfo(rls.Namespace, pi.InfraRef)
	res := rlsRef.Resource
	if err := s.ExecuteRemoteSql(sqlList, &res, pi.Database, sqlType); err != nil {
		task.Status = c7nctx.FailedStatus
		task.Reason = err.Error()
		return err
	}
	return nil
}

func (pi *ReleaseJob) ExecuteRequests(infra *Release, s *slaver.Slaver) error {
	if pi.Request == nil {
		return nil
	}
	_, r := c7nctx.Ctx.GetJobInfo(pi.Name)
	if r != nil && r.Type == c7nctx.HttpGetTask {
		log.Infof("task %s had executed", pi.Name)
		return nil
	}

	r = &c7nctx.TaskInfo{
		Name:     pi.Name,
		RefName:  infra.Name,
		Type:     c7nctx.TaskType,
		Status:   c7nctx.SucceedStatus,
		TaskType: c7nctx.HttpGetTask,
		Version:  infra.Version,
	}

	defer c7nctx.Ctx.AddJobInfo(r)

	req := pi.Request
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
		r.Status = c7nctx.FailedStatus
		r.Reason = err.Error()
	}
	return err
}

func (r *Release) String() string {
	return r.Name
}
func (r *Release) ExecutePreCommands(s *slaver.Slaver) error {

	err := r.executeExternalFunc(r.PreInstall, s)
	return err
}

func (r *Release) executeExternalFunc(c []ReleaseJob, s *slaver.Slaver) error {
	for _, pi := range c {
		if len(pi.Commands) > 0 {
			if err := pi.ExecuteSql(r, "mysql", s); err != nil {
				return err
			}
		}
		if len(pi.Mysql) > 0 {
			if err := pi.ExecuteSql(r, "mysql", s); err != nil {
				return err
			}
		}
		if len(pi.Psql) > 0 {
			if err := pi.ExecuteSql(r, "postgres", s); err != nil {
				return err
			}
		}
		if pi.Request != nil {
			if err := pi.ExecuteRequests(r, s); err != nil {
				return err
			}
		}
	}
	return nil
}

// 将 config.yml 中的值合并到 Release.Resource
func (r *Release) mergerResource(uc *config.C7nConfig) {
	cnf := uc
	if res := cnf.GetResource(r.Name); res == nil {
		log.Warnf("There is no resource in config.yaml of Release %s", r.Name)
	} else {
		// 直接使用外部配置
		if res.External {
			r.Resource = res
		} else {
			// TODO 有没有更加简便的方式
			if res.Domain != "" {
				r.Resource.Domain = res.Domain
			}
			if res.Schema != "" {
				r.Resource.Schema = res.Schema
			}
			if res.Url != "" {
				r.Resource.Url = res.Url
			}
			if res.Host != "" {
				r.Resource.Host = res.Host
			}
			if res.Port > 0 {
				r.Resource.Port = res.Port
			}
			if res.Username != "" {
				r.Resource.Username = res.Username
			}
			if res.Password != "" {
				r.Resource.Password = res.Password
			}
			if res.Persistence != nil {
				r.Resource.Persistence = res.Persistence
			}
		}
	}
}

func (r *Release) getUserValuesTpl(valuesPath string) ([]byte, error) {
	return c7nctx.Ctx.UserConfig.GetHelmValuesTpl(r.Name)
}

// convert yml format values template to yaml raw data
func (r *Release) ValuesRaw(uc *config.C7nConfig) (string, error) {
	dir := uc.Spec.HelmConfig.Values.Dir
	c7nversion := uc.Version
	if dir == "" {
		dir = "values"
	}
	// values.yaml 与 rls 名一致
	valuesFilepath := fmt.Sprintf("%s/%s.yaml", dir, r.Name)

	var data []byte
	_, err := os.Stat(valuesFilepath)
	// 当本地文件不存在时拉取远程文件
	if err != nil {
		if os.IsNotExist(err) {
			url := fmt.Sprintf(consts.RemoteInstallResourceRootUrl, c7nversion, r.Name)
			nData, statusCode, err := downloader.GetFileContent(url)
			if statusCode == 200 && err == nil {
				data = nData
			}
		} else {
			return "", std_errors.WithMessage(err, fmt.Sprintf("get user values for %s failed", r.Name))
		}
	}
	if err == nil {
		data, _ = ioutil.ReadFile(valuesFilepath)
	}

	if len(data) > 0 {
		return string(data[:]), nil
	}
	return "", nil
}

// convert yml values to values list as xxx=yyy
func (r *Release) HelmValues() []string {
	values := make([]string, len(r.Values))
	// store values for feature use
	for k, v := range r.Values {
		// 解决特殊字符
		values[k] = fmt.Sprintf("%s=%s", v.Name, v.Value)
	}
	return values
}

// resource infra
func (r *Release) Install(s *slaver.Slaver, vals map[string]interface{}) error {
	ti, err := c7nctx.GetReleaseTaskInfo(r.Namespace, r.Name)
	if err != nil {
		return err
	}
	if ti.Status == c7nctx.SucceedStatus {
		log.Infof("Release %s is already installed", r.Name)
		return nil
	}
	if ti.Status == c7nctx.RenderedStatus || ti.Status == c7nctx.FailedStatus {
		// 等待依赖项安装完成
		for _, r := range r.Requirements {
			CheckReleasePodRunning(r)
		}

		if err := r.ExecutePreCommands(s); err != nil {
			return std_errors.WithMessage(err, fmt.Sprintf("Release %s execute pre commands failed", r.Name))
		}

		values := r.HelmValues()
		var releaseName string
		if c7nctx.Ctx.Prefix != "" {
			releaseName = fmt.Sprintf("%s-%s", r.Prefix, r.Name)
		} else {
			releaseName = r.Name
		}

		chartArgs := client.ChartArgs{
			ReleaseName: releaseName,
			Namespace:   c7nctx.Ctx.Namespace,
			RepoUrl:     c7nctx.Ctx.RepoUrl,
			Verify:      false,
			Version:     r.Version,
			ChartName:   r.Chart,
		}

		log.Infof("installing %s", r.Name)
		for _, k := range values {
			log.WithField("release", r.Name).Debug(k)
		}

		if err != nil {
			ti.Status = c7nctx.FailedStatus
			c7nctx.Ctx.UpdateJobInfo(ti)
			return err
		} else {
			ti.Status = c7nctx.InstalledStatus
			c7nctx.Ctx.UpdateJobInfo(ti)
		}
	}
	if ti.Status == c7nctx.InstalledStatus {
		// TODO 解决 after task失败二次执行
		if len(r.AfterInstall) > 0 {
			go r.ExecuteAfterTasks()
			// return std_errors.WithMessage(err, "Execute after task failed")
		}
		ti.Status = c7nctx.SucceedStatus
		// 更新 jobinfo 状态
		c7nctx.Ctx.UpdateJobInfo(ti)
	}

	return nil
}

func (r *Release) InstallComponent() error {

	values := r.HelmValues()
	releaseName := r.Name
	if r.Prefix != "" {
		releaseName = fmt.Sprintf("%s-%s", c7nctx.Ctx.Prefix, r.Name)
	}
	chartArgs := client.ChartArgs{
		ReleaseName: releaseName,
		Namespace:   c7nctx.Ctx.Namespace,
		RepoUrl:     c7nctx.Ctx.RepoUrl,
		Verify:      false,
		Version:     r.Version,
		ChartName:   r.Chart,
	}

	log.Infof("installing %s", r.Name)
	for _, k := range values {
		log.Debug(k)
	}
	if r.Timeout > 0 {
		values = append(values, fmt.Sprintf("preJob.timeout=%d", r.Timeout))
	}
	// raw := r.ValuesRaw()
	err := c7nctx.Ctx.HelmClient.InstallRelease(values, "", chartArgs)
	return err
}

//
func (r *Release) ExecuteAfterTasks(s *slaver.Slaver, wg *sync.WaitGroup) error {
	ti, err := c7nctx.GetTaskFromCM(r.Namespace, r.Name)
	if err != nil {
		return err
	}
	defer c7nctx.UpdateTaskToCM(r.Namespace, *ti)

	CheckReleasePodRunning(r.Name)

	log.Infof("%s: started, will execute required commands and requests", r.Name)
	err = r.executeExternalFunc(r.AfterInstall, s)
	if err != nil {
		log.Error(err)
		ti.Status = c7nctx.FailedStatus
		return err
	}
	ti.Status = c7nctx.SucceedStatus
	wg.Done()
	return nil
}

// 基础组件——比如 gitlab-ha ——有 app 标签，c7n 有 choerodon.io/release 标签
// TODO 所有组件设置统一的label
func CheckReleasePodRunning(rls string) {

	clientset := *c7nctx.Ctx.KubeClient
	namespace := c7nctx.Ctx.Namespace
	time.Sleep(time.Second * 3)

	labels := []string{
		fmt.Sprintf("choerodon.io/release=%s", rls),
		fmt.Sprintf("app=%s", rls),
	}
	for {
		for _, label := range labels {
			deploy, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), meta_v1.ListOptions{LabelSelector: label})
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
			ss, err := clientset.AppsV1().StatefulSets(namespace).List(context.Background(), meta_v1.ListOptions{LabelSelector: label})
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
