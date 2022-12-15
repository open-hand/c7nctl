package main

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/utils"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/yidaqiang/go-chartmuseum"
	"os/exec"
	"strings"

	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

const packageDesc = `Generate a Choerodon offline installation package`

type packageOption struct {
	configFile string
}

// upgradeCmd represents the upgrade command
func newPackageCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	var pkgOpt = packageOption{}
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Generate a Choerodon offline installation package",
		Long:  packageDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			cvm, err := getPackageConfig(pkgOpt.configFile)
			if err != nil {
				return err
			}
			chartPath := fmt.Sprintf("./choerodon-offline-%s/chart", cvm.Spec.VersionRegexp)
			imagePath := fmt.Sprintf("./choerodon-offline-%s/image/", cvm.Spec.VersionRegexp)

			chartClient, _ := chartmuseum.NewClient(chartmuseum.WithBaseURL(cvm.Spec.Chart.DefaultSource.Url))
			chart := cvm.Spec.Chart

			sureFilePath(cvm.Spec.VersionRegexp)

			for _, c := range cvm.Spec.Chart.Component {
				if c.Version == "" {
					c.Version, _ = utils.GetReleaseTag(chart.DefaultSource.Url+"/"+chart.DefaultSource.Repo, c.Name, cvm.Spec.VersionRegexp)
					logrus.Debugf("Chart %s version is %s\n", c.Name, c.Version)
				}
				chartVersionOption := chartmuseum.NewChartVersionOption(c.Name, c.Version)
				_, err = chartClient.Charts.DownloadChart(chart.DefaultSource.Repo, chartPath, chartVersionOption)
				if err != nil {
					logrus.Error(err)
				}
			}

			imageSet := mapset.NewSet[string]()
			complieRegex := regexp.MustCompile("image: (.*?)\n")
			files, _ := ioutil.ReadDir(chartPath)
			for _, fi := range files {
				if fi.IsDir() {
					//listAll(path + "/" + fi.Name())
					logrus.Debug("skip up dir")
				} else {
					chartfile := chartPath + "/" + fi.Name()

					template, err := cfg.HelmClient.Template(chartfile, out)
					if err != nil {
						return err
					}
					matchArr := complieRegex.FindAllStringSubmatch(template, -1)
					for _, image := range matchArr {
						imageSet.Add(image[1])
					}
				}
			}
			imageSet.Each(func(image string) bool {
				cmd := exec.Command("docker", "pull", image)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()

				if err != nil {
					logrus.Fatalf("cmd.Run() failed with %s\n", err)
				}
				imageArr := strings.Split(image, "/")
				imageNameAndVersion := strings.Split(imageArr[len(imageArr)-1], ":")
				imagePathOne := imagePath + imageNameAndVersion[0] + "-" + imageNameAndVersion[1] + ".tar"
				_, exist := utils.IsFileExist(imagePathOne)
				if exist {
					logrus.Infof("skip image %s", image)
				} else {
					cmd = exec.Command("docker", "save", image, "-o", imagePathOne)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					if err != nil {
						logrus.Fatalf("cmd.Run() failed with %s\n", err)
					}
				}

				return false
			})
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringVarP(&pkgOpt.configFile, "config", "c", "config.yaml", "离线包定义文件")

	return cmd
}

func sureFilePath(version string) {
	err := os.MkdirAll(fmt.Sprintf("./choerodon-offline-%s/chart", version), 0766)
	if err != nil {
		return
	}
	err = os.MkdirAll(fmt.Sprintf("./choerodon-offline-%s/image", version), 0766)
	if err != nil {
		return
	}
}

func getPackageConfig(file string) (*config.ChoerodonVersion, error) {
	if file == "" {
		return nil, errors.New("--config 不能为空")
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, errors.New("配置文件不存在")
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.WithMessagef(err, "读取配置文件 %s 失败: ", file)
	}

	cvm := &config.ChoerodonVersion{}
	if err = yaml.Unmarshal(data, cvm); err != nil {
		return nil, errors.WithMessagef(err, "序列化配置文件 %s 失败: ", file)
	}
	logrus.Infof("成功读取配置文件 %s", file)
	return cvm, nil
}
