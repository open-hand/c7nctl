package consts

const (
	StaticLogsCM        = "c7n-logs"
	StaticReleaseKey    = "release"
	StaticTaskKey       = "task"
	StaticPersistentKey = "persistent"
	PvType              = "pv"
	PvcType             = "pvc"
	CRDType             = "crd"
	ReleaseTYPE         = "helm"
	TaskType            = "task"
	UninitializedStatus = "uninitialized"
	SucceedStatus       = "succeed"
	FailedStatus        = "failed"
	InstalledStatus     = "installed"
	RenderedStatus      = "rendered"
	// if have after process while wait
	CreatedStatus      = "created"
	staticInstalledKey = "installed"
	staticExecutedKey  = "execute"
	SqlTask            = "sql"
	HttpGetTask        = "httpGet"
)
