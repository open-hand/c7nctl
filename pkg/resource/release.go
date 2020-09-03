package resource

import (
	"encoding/json"
	"fmt"
	c7nclient "github.com/choerodon/c7nctl/pkg/client"
	"github.com/choerodon/c7nctl/pkg/common/consts"
	c7nerrors "github.com/choerodon/c7nctl/pkg/common/errors"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/choerodon/c7nctl/pkg/slaver"
	"github.com/choerodon/c7nctl/pkg/utils"
	std_errors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

type Release struct {
	Name         string
	Chart        string
	Version      string
	Namespace    string
	RepoURL      string
	Values       []c7nclient.ChartValue
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
	Header     []c7nclient.ChartValue
	Url        string
	Parameters []c7nclient.ChartValue
	Body       string
	Method     string
}

// TODO 移动到 action 包
func (r *Release) InstallComponent() error {
	values := r.HelmValues()
	releaseName := r.Name
	if r.Prefix != "" {
		releaseName = fmt.Sprintf("%s-%s", r.Prefix, r.Name)
	}
	chartArgs := c7nclient.ChartArgs{
		ReleaseName: releaseName,
		Namespace:   r.Namespace,
		// RepoUrl:     r.RepoUrl,
		Verify:    false,
		Version:   r.Version,
		ChartName: r.Chart,
	}

	log.Infof("installing %s", r.Name)
	for _, k := range values {
		log.Debug(k)
	}
	if r.Timeout > 0 {
		values = append(values, fmt.Sprintf("preJob.timeout=%d", r.Timeout))
	}
	// raw := r.ValuesRaw()
	fmt.Printf("%+v", chartArgs)
	//err := c7nctx.Ctx.HelmClient.InstallRelease(values, "", chartArgs)
	return nil
}

// 执行 after Task，完成后更新任务状态，并执行 wg.done
func (r *Release) ExecuteAfterTasks(s *slaver.Slaver) error {

	log.Infof("%s: started, will execute required commands and requests", r.Name)
	return r.executeExternalFunc(r.AfterInstall, s)
}

func (r *Release) ExecutePreCommands(s *slaver.Slaver) error {
	log.Infof("%s: started, will execute required commands and requests", r.Name)
	err := r.executeExternalFunc(r.PreInstall, s)
	return err
}

func (r *Release) executeExternalFunc(c []ReleaseJob, s *slaver.Slaver) error {
	for _, pi := range c {
		if len(pi.Commands) > 0 {
			if err := pi.executeSql(r, "mysql", s); err != nil {
				return err
			}
		}
		if len(pi.Mysql) > 0 {
			if err := pi.executeSql(r, "mysql", s); err != nil {
				return err
			}
		}
		if len(pi.Psql) > 0 {
			if err := pi.executeSql(r, "postgres", s); err != nil {
				return err
			}
		}
		if pi.Request != nil {
			if err := pi.executeRequests(r, s); err != nil {
				return err
			}
		}
	}
	return nil
}

func (pi *ReleaseJob) executeSql(rls *Release, sqlType string, s *slaver.Slaver) error {

	task, err := c7nclient.GetTask(pi.Name)
	if err != nil {
		if std_errors.Is(err, c7nerrors.TaskInfoIsNotFoundError) {
			task = c7nclient.NewReleaseJobTask(pi.Name, consts.SqlTask, consts.Version)
		} else {
			return err
		}
	}
	defer c7nclient.SaveTask(*task)

	if task.Status == consts.SucceedStatus {
		log.Infof("task %s of %s had executed", pi.Name, rls.Name)
		return nil
	}

	log.Infof("executing %s , %s", rls.Name, pi.Name)
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
	rlsRef, err := c7nclient.GetTask(pi.InfraRef)
	if err != nil {
		return err
	}
	if err := s.ExecuteRemoteSql(sqlList, &rlsRef.Resource, pi.Database, sqlType); err != nil {
		task.Status = consts.FailedStatus
		task.Reason = err.Error()
		return err
	}
	task.Status = consts.SucceedStatus
	return nil
}

func (pi *ReleaseJob) executeRequests(rls *Release, s *slaver.Slaver) error {
	if pi.Request == nil {
		return nil
	}
	task, err := c7nclient.GetTask(pi.Name)
	if err != nil {
		if std_errors.Is(err, c7nerrors.TaskInfoIsNotFoundError) {
			task = &c7nclient.TaskInfo{
				Name:     pi.Name,
				RefName:  rls.Name,
				Type:     consts.TaskType,
				Status:   consts.UninitializedStatus,
				TaskType: consts.HttpGetTask,
				Version:  rls.Version,
			}
			task = c7nclient.NewReleaseJobTask(pi.Name, consts.HttpGetTask, consts.Version)
		} else {
			return err
		}
	}
	if task.Type == consts.SucceedStatus {
		log.Infof("task %s had executed", pi.Name)
		return nil
	}
	defer c7nclient.SaveTask(*task)

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

	_, err = s.ExecuteRemoteRequest(f)
	if err != nil {
		task.Status = consts.FailedStatus
		task.Reason = err.Error()
		return err
	}
	task.Status = consts.SucceedStatus
	return nil
}

// 将 config.yml 中的值合并到 Release.Resource
func (r *Release) mergerResource(uc *config.C7nConfig) {
	cnf := uc
	if res := cnf.GetResource(r.Name); res == nil {
		log.Debugf("There is no resource in config.yaml of Release %s", r.Name)
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

// convert yml format values template to yaml raw data
// 获取 resourcePath 路径下的 values 文件
func (r *Release) ValuesRaw(resourcePath, helmValue string) (string, error) {
	// values.yaml 与 r 名一致
	valuesFilepath := fmt.Sprintf(filepath.Join(resourcePath, helmValue, r.Name) + ".yaml")
	data, err := utils.GetResource(valuesFilepath)
	if err != nil {
		log.Debugf("load helm values file %s failed: %+v", valuesFilepath, err)
	}
	// 不存在配置文件的返回空字符串
	return string(data[:]), nil
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

func (r *Release) String() string {
	b, _ := json.MarshalIndent(*r, "\t", "\t")
	return string(b)
}
