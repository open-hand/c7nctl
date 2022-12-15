package client

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	liberrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"io"
	"os"
	"strings"
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

	if err := actionConfig.Init(kube.GetConfig(kubeconfig, "", namespace), namespace, helmDriver, log.Debugf); err != nil {
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

func (h *Helm3Client) newHelm3Template(cfg *action.Configuration) *action.Install {
	client := action.NewInstall(cfg)

	var validate bool
	var includeCrds bool
	var extraAPIs []string

	client.DryRun = true
	client.ReleaseName = "release-name"
	client.Replace = true // Skip the name check
	client.ClientOnly = !validate
	client.APIVersions = chartutil.VersionSet(extraAPIs)
	client.IncludeCRDs = includeCrds

	return client
}
func (h *Helm3Client) newHelm3Upgrade(cfg *action.Configuration, args ChartArgs) *action.Upgrade {
	upgrade := action.NewUpgrade(cfg)
	upgrade.ChartPathOptions = action.ChartPathOptions{
		CaFile:   args.CaFile,
		CertFile: args.CertFile,
		KeyFile:  args.KeyFile,
		Keyring:  args.Keyring,
		RepoURL:  args.RepoUrl,
		Verify:   args.Verify,
		Version:  args.Version,
	}
	// 默认更新或者安装
	upgrade.Install = true
	//upgrade.CreateNamespace = createNamespace
	//upgrade.DryRun = false
	//upgrade.DisableHooks = false
	//upgrade.SkipCRDs = false
	//upgrade.Timeout = xx
	//upgrade.Wait = false
	//upgrade.Devel = false
	upgrade.Namespace = args.Namespace
	//upgrade.Atomic = client.Atomic
	//upgrade.PostRenderer = client.PostRenderer
	//upgrade.DisableOpenAPIValidation = client.DisableOpenAPIValidation
	//upgrade.SubNotes = client.SubNotes
	return upgrade
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

func (h *Helm3Client) Upgrade(cArgs ChartArgs, vals map[string]interface{}, out io.Writer) (*release.Release, error) {
	client := h.newHelm3Upgrade(h.Configuration, cArgs)
	// Fixes #7002 - Support reading values from STDIN for `upgrade` command
	// Must load values AFTER determining if we have to call install so that values loaded from stdin are are not read twice
	if client.Install {
		// If a release does not exist, install it.
		histClient := action.NewHistory(h.Configuration)
		histClient.Max = 1
		if _, err := histClient.Run(cArgs.ReleaseName); err == driver.ErrReleaseNotFound {
			// Only print this to stdout for table output

			instClient := action.NewInstall(h.Configuration)
			instClient.CreateNamespace = true
			instClient.ChartPathOptions = client.ChartPathOptions
			instClient.DryRun = client.DryRun
			instClient.DisableHooks = client.DisableHooks
			instClient.SkipCRDs = client.SkipCRDs
			instClient.Timeout = client.Timeout
			instClient.Wait = client.Wait
			instClient.Devel = client.Devel
			instClient.Namespace = client.Namespace
			instClient.Atomic = client.Atomic
			instClient.PostRenderer = client.PostRenderer
			instClient.DisableOpenAPIValidation = client.DisableOpenAPIValidation
			instClient.SubNotes = client.SubNotes

			rel, err := h.Install(cArgs, vals, out)
			if err != nil {
				return nil, err
			}
			return rel, nil
		} else if err != nil {
			return nil, err
		}
	}

	log.Debugf("Original chart version: %q", client.Version)
	if client.Version == "" && client.Devel {
		log.Debugf("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}
	settings := cli.New()

	chartPath, err := client.ChartPathOptions.LocateChart(cArgs.ChartName, settings)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	ch, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return nil, err
		}
	}

	if ch.Metadata.Deprecated {
		fmt.Fprintln(out, "WARNING: This chart is deprecated")
	}

	rel, err := client.Run(cArgs.ReleaseName, ch, vals)
	if err != nil {
		return nil, liberrors.Wrap(err, "UPGRADE FAILED")
	}

	return rel, nil
}

func (h *Helm3Client) Template(chartFile string, out io.Writer) (string, error) {
	client := h.newHelm3Template(h.Configuration)
	valueOpts := &values.Options{}

	rel, err := runInstall([]string{chartFile}, client, valueOpts, out)

	if err != nil {
		return "", err
	}

	// We ignore a potential error here because, when the --debug flag was specified,
	// we always want to print the YAML, even if it is not valid. The error is still returned afterwards.
	var manifests bytes.Buffer
	if rel != nil {

		fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))

		if !client.DisableHooks {
			for _, m := range rel.Hooks {
				fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", m.Path, m.Manifest)
			}
		}
	}

	return manifests.String(), err
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

func runInstall(args []string, client *action.Install, valueOpts *values.Options, out io.Writer) (*release.Release, error) {
	log.Debugf("Original chart version: %q", client.Version)
	if client.Version == "" && client.Devel {
		log.Debugf("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	name, chart, err := client.NameAndChart(args)
	if err != nil {
		return nil, err
	}
	client.ReleaseName = name
	settings := cli.New()
	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return nil, err
	}

	log.Debug("CHART PATH: %s\n", cp)

	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		log.Warn("This chart is deprecated")
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
					return nil, errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	client.Namespace = settings.Namespace()
	return client.Run(chartRequested, vals)
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
