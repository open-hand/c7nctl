package client

import (
	"fmt"
	"github.com/ghodss/yaml"
	liberrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"
	"io"
	"os"
)

var helmClient *Helm3Client

// helmClient 使用单例模式，
type Helm3Client struct {
	*action.Configuration
}

type ChartArgs struct {
	RepoUrl     string
	Version     string
	Namespace   string
	ReleaseName string
	Verify      bool
	Keyring     string
	CertFile    string
	KeyFile     string
	CaFile      string
	ChartName   string
}

func NewHelm3Client(cfg *action.Configuration) *Helm3Client {
	return &Helm3Client{
		Configuration: cfg,
	}
}

func InitConfiguration(kubeconfig, namespace string) *action.Configuration {
	actionConfig := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")
	// TODO 是否
	if err := actionConfig.Init(kube.GetConfig(kubeconfig, "", namespace), namespace, helmDriver, func(format string, v ...interface{}) {
		log.Warnf(format, v)
	}); err != nil {
		log.Fatal(err)
	}
	return actionConfig
}

func (h *Helm3Client) Install(cArgs ChartArgs, vals map[string]interface{}, out io.Writer) (*release.Release, error) {
	client := h.newHelm3Install(h.Configuration, cArgs)

	log.Debugf("Original chart version: %q", client.Version)
	if client.Version == "" && client.Devel {
		log.Debugf("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}
	// TODO 移动到 helm3Client
	os.Setenv("HELM_NAMESPACE", cArgs.Namespace)
	settings := cli.New()
	cp, err := client.ChartPathOptions.LocateChart(cArgs.ChartName, settings)
	if err != nil {
		return nil, err
	}

	log.Debugf("CHART PATH: %s\n", cp)

	p := getter.All(settings)

	// 直接传入 vals。因为 values.yaml 是模版配置文件，所以需要提前渲染
	/*
		vals, err := valueOpts.MergeValues(p)
		if err != nil {
			return nil, err
		}
	*/

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		fmt.Fprintln(out, "WARNING: This chart is deprecated")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              out,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, liberrors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	return client.Run(chartRequested, vals)
}

func (h *Helm3Client) newHelm3Install(cfg *action.Configuration, args ChartArgs) *action.Install {
	install := action.NewInstall(cfg)
	install.ChartPathOptions = action.ChartPathOptions{
		CaFile:   args.CaFile,
		CertFile: args.CertFile,
		KeyFile:  args.KeyFile,
		Keyring:  args.Keyring,
		RepoURL:  args.RepoUrl,
		Verify:   args.Verify,
		Version:  args.Version,
	}
	install.ReleaseName = args.ReleaseName
	install.Namespace = args.Namespace

	return install
}

func RunHelmInstall(client *action.Install, chart string, vals map[string]interface{}, out io.Writer) (*release.Release, error) {
	log.Debugf("Original chart version: %q", client.Version)
	if client.Version == "" && client.Devel {
		log.Debugf("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}
	settings := cli.New()
	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return nil, err
	}

	log.Debugf("CHART PATH: %s\n", cp)

	p := getter.All(settings)
	// 直接传入 vals

	/*
		vals, err := valueOpts.MergeValues(p)
		if err != nil {
			return nil, err
		}
	*/

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		fmt.Fprintln(out, "WARNING: This chart is deprecated")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              out,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, liberrors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}
	client.Namespace = settings.Namespace()
	return client.Run(chartRequested, vals)
}

// isChartInstallable validates if a chart can be installed
//
// Application chart type is only installable
func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, liberrors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

// TODO 移动到 渲染 vals 的地方
func prepareValues(values []string) ([]byte, error) {
	// User specified a value via --set
	base := map[string]interface{}{}
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return []byte{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}
	return yaml.Marshal(base)
}
