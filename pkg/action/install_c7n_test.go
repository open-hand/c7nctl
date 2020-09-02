package action

import (
	"testing"
)

func TestGetInstallDef(t *testing.T) {
	/*
		install := initInstall()
		userConf := GetUserConfig("../../example/config.yml")

		install.GetInstallDef(userConf)
	*/
}

func initInstall() *Choerodon {
	/*
		cfg := NewCfg()
		setSettings(cfg.HelmInstall.Settings())
		_ = cfg.HelmInstall.SetupConnection()
		defer cfg.HelmInstall.Teardown()

		cfg.InitCfg()
		i := NewChoerodon(cfg)
		i.Version = "release-0.21"

		context.Ctx.Metrics.Mail = i.Mail
		context.Ctx = context.Context{
			// also init i.cfg
			HelmClient:   i.Cfg.HelmInstall,
			KubeClient:   i.Cfg.KubeClient,
			CommonLabels: i.CommonLabels,
			Metrics:      context.Ctx.Metrics,
			SkipInput:    i.SkipInput,
			Prefix:       i.Prefix,
		}

		return i
	*/
}

func setSettings(settings *helm_env.EnvSettings) {
	/*
		if settings.TLSCaCertFile == helm_env.DefaultTLSCaCert || settings.TLSCaCertFile == "" {
			settings.TLSCaCertFile = settings.Home.TLSCaCert()
		} else {
			settings.TLSCaCertFile = os.ExpandEnv(settings.TLSCaCertFile)
		}
		if settings.TLSCertFile == helm_env.DefaultTLSCert || settings.TLSCertFile == "" {
			settings.TLSCertFile = settings.Home.TLSCert()
		} else {
			settings.TLSCertFile = os.ExpandEnv(settings.TLSCertFile)
		}
		if settings.TLSKeyFile == helm_env.DefaultTLSKeyFile || settings.TLSKeyFile == "" {
			settings.TLSKeyFile = settings.Home.TLSKey()
		} else {
			settings.TLSKeyFile = os.ExpandEnv(settings.TLSKeyFile)
		}
	*/
}

func TestNewInstallQueue(t *testing.T) {
	/*
		install := initInstall()
		userConf := GetUserConfig("../../example/config.yml")
		install.ResourceFile = "../../manifests/install.yml"
		context.Ctx.UserConfig = userConf
		context.Ctx.Namespace = userConf.Metadata.Namespace

		id := install.GetInstallDef(userConf)
		graph := resource.NewReleaseGraph(id.Spec.Release)
		queue := graph.TopoSortByKahn()
		for !queue.IsEmpty() {
			rls := queue.Dequeue()
			t.Log(rls.Name)
		}
	*/
}

func TestRenderRelease(t *testing.T) {
	/*
		install := initInstall()
		userConf := GetUserConfig("../../example/config.yml")
		install.ResourceFile = "../../manifests/install.yml"
		context.Ctx.UserConfig = userConf
		context.Ctx.Namespace = userConf.Metadata.Namespace

		id := install.GetInstallDef(userConf)

		// 渲染 Release
		for _, rls := range id.Spec.Release {
			// 传入参数的是 *Release
			RenderRelease_old(rls)
		}
		for _, rls := range id.Spec.Release {
			t.Log(rls.Name)
		}
	*/
}
