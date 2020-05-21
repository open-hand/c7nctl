package resource

import (
	"bytes"
	"fmt"
	"github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/consts"
	"github.com/choerodon/c7nctl/pkg/context"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/choerodon/c7nctl/pkg/utils"
	c7n_utils "github.com/choerodon/c7nctl/pkg/utils"
	"github.com/vinkdong/gox/http/downloader"
	"github.com/vinkdong/gox/log"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"text/template"
	"time"
)

type Release struct {
	Name      string
	Chart     string
	Version   string
	Namespace string
	RepoURL   string
	// TODO change value set method
	Values []context.ChartValue
	// TODO remove Persistence
	Persistence  []*Persistence
	PreInstall   []ReleaseJob
	AfterInstall []ReleaseJob
	// preValues unused
	PreValues    context.PreValueList
	Requirements []string
	Health       context.Health

	Resource *config.Resource
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

type PreValueList []*PreValue

type ChartValue struct {
	Name  string
	Value string
	Input context.Input
	Case  string
	Check string
}

type PreValue struct {
	Name  string
	Value string
	Check string
	Input context.Input
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

func (pi *ReleaseJob) ExecuteSql(infra *Release, sqlType string) error {

	news := context.Ctx.GetSucceedTask(pi.Name, infra.Name, context.SqlTask)
	if news != nil {
		log.Successf("task %s had executed", pi.Name)
		return nil
	}
	log.Infof("executing %s , %s", infra.Name, pi.Name)

	news = &context.JobInfo{
		Name:     pi.Name,
		RefName:  infra.Name,
		Type:     context.TaskType,
		Status:   context.SucceedStatus,
		TaskType: context.SqlTask,
		Version:  infra.Version,
	}

	defer context.Ctx.SaveNews(news)

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
	r := utils.GetResource(pi.InfraRef)
	s := context.Ctx.Slaver
	if err := s.ExecuteRemoteSql(sqlList, r, pi.Database, sqlType); err != nil {
		news.Status = context.FailedStatus
		news.Reason = err.Error()
		return err
	}
	return nil
}

func (pl *PreValueList) prepareValues() error {

	for _, v := range *pl {
		if err := v.renderValue(); err != nil {
			return err
		}
	}
	return nil
}

func (pl *PreValueList) getValues(key string) string {
	for _, v := range *pl {
		if v.Name == key {
			return v.Value
		}
	}
	return ""
}

func (p *PreValue) RandomToken(length int) string {
	return c7n_utils.RandomToken(length)
}

func (p *PreValue) RandomLowCaseToken(length int) string {
	return c7n_utils.GenerateRunnerToken(length)
}

func (p *PreValue) renderValue() error {

	var value string
	if p.Input.Enabled && !context.Ctx.SkipInput {
		log.Lock()
		var err error
		if p.Input.Password {
			p.Input.Twice = true
			value, err = utils.AcceptUserPassword(p.Input)
		} else {
			value, err = utils.AcceptUserInput(p.Input)
		}
		log.Unlock()
		if err != nil {
			log.Error(err)
			os.Exit(128)
		}
	} else {
		tpl, err := template.New(p.Name).Parse(p.Value)
		if err != nil {
			return err
		}
		var data bytes.Buffer
		err = tpl.Execute(&data, p)
		if err != nil {
			return err
		}
		value = data.String()
	}

	switch p.Check {
	case "clusterdomain":
		//todo: add check domain
		log.Debugf("PreValue %s: %s, checking: %s", p.Name, p.Value, p.Check)
		if err := context.Ctx.Slaver.CheckClusterDomain(p.Value); err != nil {
			log.Errorf("请检查您的域名: %s 已正确解析到集群", value)
			return err
		}
	}

	p.Value = value
	return nil
}

// 获取基础组件信息
func (p *PreValue) GetResource(key string) *config.Resource {
	news := context.Ctx.GetSucceed(key, context.ReleaseTYPE)
	// get info from succeed
	if news != nil {
		return &news.Resource
	} else {
		// 从用户配置文件中读取
		if r, ok := context.Ctx.UserConfig.Spec.Resources[key]; ok {
			return r
		}
	}
	log.Errorf("can't get required resource [%s]", key)
	context.Ctx.CheckExist(188)
	return nil
}

func (pi *ReleaseJob) ExecuteRequests(infra *Release) error {
	if pi.Request == nil {
		return nil
	}
	news := context.Ctx.GetSucceedTask(pi.Name, infra.Name, context.HttpGetTask)
	if news != nil {
		log.Successf("task %s had executed", pi.Name)
		return nil
	}

	news = &context.JobInfo{
		Name:     pi.Name,
		RefName:  infra.Name,
		Type:     context.TaskType,
		Status:   context.SucceedStatus,
		TaskType: context.HttpGetTask,
		Version:  infra.Version,
	}

	defer context.Ctx.SaveNews(news)

	// pi.Request.Render(infra)
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
		news.Status = context.FailedStatus
		news.Reason = err.Error()
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

func (rls *Release) GetUserStorageClassName() string {
	return context.Ctx.UserConfig.GetStorageClassName()
}

func (rls *Release) IgnorePv() bool {
	return context.Ctx.UserConfig.IgnorePv()
}

func (rls *Release) getUserConfig() *config.Resource {
	return context.Ctx.UserConfig.GetResource(rls.Name)
}

func (rls *Release) getUserValuesTpl() ([]byte, error) {
	return context.Ctx.UserConfig.GetHelmValuesTpl(rls.Name)
}

func (rls *Release) ApplyUserResource() error {
	r := rls.getUserConfig()
	if r == nil {
		//log.Infof("no user config resource for %s", rls.Name)
		return nil
	}
	if r.External {
		rls.Resource = r
		return nil
	}
	// just override domain,host and schema
	if r.Domain != "" {
		rls.Resource.Domain = r.Domain
	}

	if r.Schema != "" && rls.Resource.Schema == "" {
		rls.Resource.Schema = r.Schema
	}

	return nil
}

func (rls *Release) ExecutePreValues() error {
	/*return rls.PreValues.prepareValues()*/
	return nil
}

func (rls *Release) renderValue(tplString string) string {
	tpl, err := template.New(rls.Name).Parse(tplString)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	var data bytes.Buffer
	err = tpl.Execute(&data, rls)
	if err != nil {
		log.Error(err)
		os.Exit(255)
	}
	return data.String()
}

func (rls *Release) GetPods() (*core_v1.PodList, error) {
	selectLabel := make(map[string]string)
	selectLabel["choerodon.io/release"] = c7n_utils.WithPrefix() + rls.Name
	set := labels.Set(selectLabel)
	opts := meta_v1.ListOptions{
		LabelSelector: set.AsSelector().String(),
	}
	return (*context.Ctx.KubeClient).CoreV1().Pods(rls.Namespace).List(opts)
}

func (rls *Release) GetPodIp() string {
reget:
	poList, err := rls.GetPods()
	if err != nil || len(poList.Items) < 1 {
		log.Errorf("can't get a pod from %s, retry...", rls.Name)
		goto reget
	}
	for _, po := range poList.Items {
		if po.Status.Phase == core_v1.PodRunning {
			return po.Status.PodIP
		}
	}
	log.Debugf("can't get a running pod from %s, retry...", rls.Name)
	time.Sleep(time.Second * 3)
	goto reget
}

func (rls *Release) GetPreValue(key string) string {
	/*return rls.PreValues.getValues(key)*/
	return ""
}

func (rls *Release) GetRequire(app string) *Release {
	news := context.Ctx.GetSucceed(app, context.ReleaseTYPE)
	i := &Release{
		Name:      app,
		Namespace: rls.Namespace,
		//PreValues: news.PreValue,
		Values: news.Values,
	}
	return i
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

func (rls *Release) GetValue(key string) string {
	for _, v := range rls.Values {
		if v.Name == key {
			return v.Value
		}
	}
	log.Infof("can't get value '%s' of %s", key, rls.Name)
	return ""
}

// only used for save log
func (rls *Release) RenderResource() config.Resource {
	//todo: just render password now, add more
	r := rls.Resource
	tpl, err := template.New(fmt.Sprintf("r-%s-%s", rls.Name, "password")).Parse(r.Password)
	if err != nil {
		log.Info(err)
		os.Exit(125)
	}
	var data bytes.Buffer
	if err := tpl.Execute(&data, rls); err != nil {
		log.Error(err)
		os.Exit(125)
	}
	r.Password = data.String()
	if r.Password != "" {
		log.Debugf("%s: resource password is %s", rls.Name, r.Password)
	}
	r.Url = rls.renderValue(r.Url)
	r.Host = rls.renderValue(r.Host)
	if r.Url != "" {
		log.Debugf("%s: resource url is %s", rls.Name, r.Url)
	}
	return *r
}

func (rls *Release) Run() error {
	log.Infof("start resource %s", rls.Name)

	if r := context.Ctx.UserConfig.GetResource(rls.Name); r != nil && r.External {
		log.Infof("using external %s", rls.Name)
	}

	// check requirement started
	// TODO make sure rls is running after install
	/*for _, r := range rls.Requirements {
		i := i.GetRelease(r)
		if i.Prefix == "" {
			i.Prefix = infra.Prefix
		}
		if err := i.CheckRunning(); err != nil {
			return err
		}
	}*/

	news := context.Ctx.GetSucceed(rls.Name, context.ReleaseTYPE)

	if news != nil {
		log.Successf("using exist release %s", news.RefName)
		if news.Status != context.SucceedStatus {
			//rls.PreValues = news.PreValue
			rls.ExecuteAfterTasks()
		}
		return nil
	}

	// 渲染 Release
	/*if err := renderRls(rls); err != nil {
		return err
	}*/

	// 执行安装前命令
	if err := rls.ExecutePreCommands(); err != nil {
		return err
	}

	statusCh := make(chan error)

	go func() {
		err := rls.Install()
		statusCh <- err
	}()

	for {
		select {
		case <-time.Tick(time.Second * 10):
			_ = rls.CatchInitJobs()
		case err := <-statusCh:
			return err
		}
	}
}

// resource infra
func (rls *Release) Install() error {
	ji := context.Ctx.GetJobInfo(rls.Name)
	// 如果已经安装直接返回
	if ji.Status == context.SucceedStatus {
		log.Infof("already installed %s", rls.Name)
		return nil
	}

	for _, r := range rls.Requirements {
		checkReleasePodRunning(r)
	}

	if err := rls.ExecutePreCommands(); err != nil {
		log.Error(err)
	}

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
		log.Debugf(k)
	}
	if rls.Timeout > 0 {
		values = append(values, fmt.Sprintf("preJob.timeout=%d", rls.Timeout))
	}
	raw := rls.ValuesRaw()
	err := context.Ctx.HelmClient.InstallRelease(values, raw, chartArgs)

	if err != nil {
		ji.Status = context.FailedStatus
		context.Ctx.UpdateJobInfo(ji)
		return err
	}

	// 更新 jobinfo 状态
	if len(rls.AfterInstall) > 0 {
		_ = rls.ExecuteAfterTasks()
	}
	ji.Status = context.SucceedStatus
	context.Ctx.UpdateJobInfo(ji)
	return nil
}

//
func (rls *Release) ExecuteAfterTasks() error {
	checkReleasePodRunning(rls.Name)

	log.Successf("%s: started, will execute required commands and requests", rls.Name)
	err := rls.executeExternalFunc(rls.AfterInstall)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (rls *Release) executeAfterTasks(task *context.BackendTask) error {
	// TODO
	checkReleasePodRunning(rls.Name)

	log.Successf("%s: started, will execute required commands and requests", rls.Name)
	err := rls.executeExternalFunc(rls.AfterInstall)
	if err != nil {
		log.Error(err)
		return err
	}

	if err := context.Ctx.UpdateCreated(rls.Name, rls.Namespace); err == nil {
		task.Success = true
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
	time.Sleep(time.Second * 5)

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
