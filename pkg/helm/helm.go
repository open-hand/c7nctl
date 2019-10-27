package helm

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/vinkdong/gox/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/repo"
	"k8s.io/helm/pkg/strvals"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	Client     *helm.Client
	Settings   helm_env.EnvSettings
	Tunnel     *kube.Tunnel
	KubeClient kubernetes.Interface
}

var (
	settings helm_env.EnvSettings
)

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

func (client *Client) locateChartPath(chartArgs ChartArgs) (string, error) {
	name := strings.TrimSpace(chartArgs.ChartName)
	version := strings.TrimSpace(chartArgs.Version)
	if fi, err := os.Stat(name); err == nil {
		abs, err := filepath.Abs(name)
		if err != nil {
			return abs, err
		}
		if chartArgs.Verify {
			if fi.IsDir() {
				return "", errors.New("cannot verify a directory")
			}
			if _, err := downloader.VerifyChart(abs, chartArgs.Keyring); err != nil {
				return "", err
			}
		}
		return abs, nil
	}
	if filepath.IsAbs(name) || strings.HasPrefix(name, ".") {
		return name, fmt.Errorf("path %q not found", name)
	}

	settings := client.Settings

	crepo := filepath.Join(settings.Home.Repository(), name)
	if _, err := os.Stat(crepo); err == nil {
		return filepath.Abs(crepo)
	}

	dl := downloader.ChartDownloader{
		HelmHome: settings.Home,
		Out:      os.Stdout,
		Keyring:  chartArgs.Keyring,
		Getters:  getter.All(settings),
	}
	if chartArgs.Verify {
		dl.Verify = downloader.VerifyAlways
	}
	if chartArgs.RepoUrl != "" {
		chartURL, err := repo.FindChartInRepoURL(chartArgs.RepoUrl, name, version,
			chartArgs.CertFile, chartArgs.KeyFile, chartArgs.CaFile, getter.All(settings))
		if err != nil {
			return "", err
		}
		name = chartURL
	}

	if _, err := os.Stat(settings.Home.Archive()); os.IsNotExist(err) {
		os.MkdirAll(settings.Home.Archive(), 0744)
	}

	filename, _, err := dl.DownloadTo(name, version, settings.Home.Archive())
	if err == nil {
		lname, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		return lname, nil
	} else if settings.Debug {
		return filename, err
	}

	return filename, fmt.Errorf("failed to download %q (hint: running `helm repo update` may help)", name)
}

// init helm client
func (client *Client) InitClient() {
	settings := &client.Settings
	if settings.TillerHost == "" && client.Tunnel != nil {
		settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", client.Tunnel.Local)
		settings.TillerConnectionTimeout = 86400
	}

	options := []helm.Option{helm.Host(settings.TillerHost), helm.ConnectTimeout(settings.TillerConnectionTimeout)}
	client.Client = helm.NewClient(options...)
	var DefaultHelmHome = filepath.Join(homedir.HomeDir(), ".helm")
	p := (*string)(&settings.Home)
	*p = DefaultHelmHome

}

func (client *Client) InstallRelease(values []string, raw string, chartArgs ChartArgs) error {

	cp, err := client.locateChartPath(chartArgs)
	if err != nil {
		return err
	}

	var rawVals []byte
	if raw == "" {
		rawVals, err = prepareValues(values)
	} else {
		rawVals = []byte(raw)
	}
	log.Debug("\n" + string(rawVals))

	chartRequested, err := chartutil.Load(cp)

	// todo: add dependencies
	//chartRequested, err := chartutil.Load(cp)
	//if err != nil {
	//	log.Error(err)
	//}
	// todo: get default ns form kube-context
	res, err := client.Client.InstallReleaseFromChart(
		chartRequested,
		chartArgs.Namespace,
		helm.ValueOverrides(rawVals),
		helm.ReleaseName(chartArgs.ReleaseName),
		helm.InstallDryRun(false),
		helm.InstallReuseName(false),
		helm.InstallDisableHooks(false),
		helm.InstallTimeout(86400),
		helm.InstallWait(false))
	if err != nil {
		return err
	}

	rel := res.GetRelease()
	if rel == nil {
		log.Errorf("can't get release name of chart %s", rel.Name)
	}
	log.Successf("installed %s", rel.Name)
	return nil
}

func (client *Client) UpgradeRelease(values []byte, chartArgs ChartArgs) error {

	cp, err := client.locateChartPath(chartArgs)
	if err != nil {
		return err
	}
	chartRequested, err := chartutil.Load(cp)

	// todo: add dependencies
	//chartRequested, err := chartutil.Load(cp)
	//if err != nil {
	//	log.Error(err)
	//}
	// todo: get default ns form kube-context
	res, err := client.Client.UpdateReleaseFromChart(
		chartArgs.ReleaseName,
		chartRequested,
		helm.UpdateValueOverrides(values),
		helm.ResetValues(true),
		helm.ReuseValues(true),
		helm.UpgradeDryRun(false),
		helm.UpgradeForce(true),
		helm.UpgradeTimeout(86400))
	if err != nil {
		return err
	}

	rel := res.GetRelease()
	if rel == nil {
		log.Errorf("can't get release name of chart %s", rel.Name)
	}
	log.Successf("Upgrade %s success", rel.Name)
	return nil
}

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

func List() {

}
