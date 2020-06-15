package client

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/vinkdong/gox/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/repo"
	"k8s.io/helm/pkg/strvals"
	"k8s.io/helm/pkg/tlsutil"
	"os"
	"path/filepath"
	"strings"
)

type HelmClient struct {

	// HelmClient is a client for working with helm
	helmClient *helm.Interface

	// setting stores setting of helm client
	settings *helm_env.EnvSettings

	tillerTunnel *kube.Tunnel
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

func (hc *HelmClient) SetupConnection() error {
	if hc.settings.TillerHost == "" {
		config, client, err := getKubeClient(hc.settings.KubeContext, hc.settings.KubeConfig)
		if err != nil {
			return err
		}

		hc.tillerTunnel, err = portforwarder.New(hc.settings.TillerNamespace, client, config)
		if err != nil {
			return err
		}

		hc.settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", hc.tillerTunnel.Local)
		// debug("Created tunnel using local port: '%d'\n", tillerTunnel.Local)
	}

	// Set up the gRPC config.
	// debug("SERVER: %q\n", settings.TillerHost)

	// Plugin support.
	return nil
}

func (hc *HelmClient) InitSettings() {
	settings := hc.settings
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
}

func (hc *HelmClient) Settings() *helm_env.EnvSettings {
	if hc.settings == nil {
		hc.settings = new(helm_env.EnvSettings)
	}
	return hc.settings
}

func (hc *HelmClient) InstallRelease(values []string, raw string, chartArgs ChartArgs) error {
	// TODO namespace check
	cp, err := hc.locateChartPath(chartArgs)
	if err != nil {
		return err
	}
	// merge values and valueFile
	rawVals, err := vals(values, raw)
	if err != nil {
		log.Error(err)
	}
	//log.Debug("\n" + string(rawVals))
	// TODO release name genarge and check

	chartRequested, err := chartutil.Load(cp)

	// todo: add dependencies
	//chartRequested, err := chartutil.Load(cp)
	//if err != nil {
	//	log.Error(err)
	//}
	// todo: get default ns form kube-context

	res, err := (*hc.helmClient).InstallReleaseFromChart(
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

func (hc *HelmClient) Teardown() {
	if hc.tillerTunnel != nil {
		hc.tillerTunnel.Close()
	}
}

// configForContext creates a Kubernetes REST client helmClient for a given kubeconfig context.
func configForContext(context string, kubeconfig string) (*rest.Config, error) {
	config, err := kube.GetConfig(context, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(context string, kubeconfig string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context, kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

// ensureHelmClient returns a new helm client impl. if h is not nil.
/*func ensureHelmClient(h helm.Interface) helm.Interface {
	if h != nil {
		return h
	}
	return NewClient()
}
*/

func (hc *HelmClient) locateChartPath(chartArgs ChartArgs) (string, error) {
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

	crepo := filepath.Join(hc.settings.Home.Repository(), name)
	if _, err := os.Stat(crepo); err == nil {
		return filepath.Abs(crepo)
	}

	dl := downloader.ChartDownloader{
		HelmHome: hc.settings.Home,
		Out:      os.Stdout,
		Keyring:  chartArgs.Keyring,
		Getters:  getter.All(*hc.settings),
		// TODO Keyring, Username, Password
	}
	if chartArgs.Verify {
		dl.Verify = downloader.VerifyAlways
	}
	if chartArgs.RepoUrl != "" {
		chartURL, err := repo.FindChartInRepoURL(chartArgs.RepoUrl, name, version,
			chartArgs.CertFile, chartArgs.KeyFile, chartArgs.CaFile, getter.All(*hc.settings))
		if err != nil {
			return "", err
		}
		name = chartURL
	}

	if _, err := os.Stat(hc.settings.Home.Archive()); os.IsNotExist(err) {
		os.MkdirAll(hc.settings.Home.Archive(), 0744)
	}

	filename, _, err := dl.DownloadTo(name, version, hc.settings.Home.Archive())
	if err == nil {
		lname, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		return lname, nil
	} else if hc.settings.Debug {
		return filename, err
	}

	return filename, fmt.Errorf("failed to download %q (hint: running `helm repo update` may help)", name)
}

func GetHelmClient(hc *HelmClient) *HelmClient {
	if hc.helmClient == nil {
		c := NewHelmClient(hc.settings)
		hc.helmClient = &c
	}
	return hc
}

func NewHelmClient(settings *helm_env.EnvSettings) helm.Interface {

	options := []helm.Option{helm.Host(settings.TillerHost), helm.ConnectTimeout(settings.TillerConnectionTimeout)}

	if settings.TLSVerify || settings.TLSEnable {
		// debug("Host=%q, Key=%q, Cert=%q, CA=%q\n", settings.TLSServerName, settings.TLSKeyFile, settings.TLSCertFile, settings.TLSCaCertFile)
		tlsopts := tlsutil.Options{
			ServerName:         settings.TLSServerName,
			KeyFile:            settings.TLSKeyFile,
			CertFile:           settings.TLSCertFile,
			InsecureSkipVerify: true,
		}
		if settings.TLSVerify {
			tlsopts.CaCertFile = settings.TLSCaCertFile
			tlsopts.InsecureSkipVerify = false
		}
		tlscfg, err := tlsutil.ClientConfig(tlsopts)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		options = append(options, helm.WithTLS(tlscfg))
	}
	return helm.NewClient(options...)
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

// Merges source and destination map, preferring values from the source map
func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}

// vals merges values from files specified via -f/--values and
// directly via --set or --set-string or --set-file, marshaling them to YAML
func vals(values []string, fileValues string) ([]byte, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	currentMap := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(fileValues), &currentMap); err != nil {
		return []byte{}, fmt.Errorf("failed to parse %s", err)
	}
	// Merge with the previous map
	base = mergeValues(base, currentMap)

	// User specified a value via --set
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return []byte{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return yaml.Marshal(base)
}
