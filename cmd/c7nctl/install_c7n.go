package main

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/action"
	"github.com/choerodon/c7nctl/pkg/config"
	c7nconsts "github.com/choerodon/c7nctl/pkg/consts"
	"github.com/choerodon/c7nctl/pkg/resource"
	c7nutils "github.com/choerodon/c7nctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
)

const installC7nDesc = `One-click installation choerodon, When your prepared k8s, helm and NFS.
To install choerodon, you must set up the choerodon install configuration file
and specify the file with "--c7n-config <install-c7n-config.yaml>".

Ensure you run this within server can vista k8s.
`

func newInstallC7nCmd(cfg *action.C7nConfiguration, out io.Writer) *cobra.Command {
	c := action.NewChoerodon(cfg, settings)

	cmd := &cobra.Command{
		Use:   "c7n",
		Short: "One-click installation choerodon",
		Long:  installC7nDesc,
		Run: func(_ *cobra.Command, args []string) {
			setUserConfig(c.SkipInput)
			if err := runInstallC7n(c); err != nil {
				log.Error(err)
				log.Error("Install Choerodon failed")
				c.Metrics.ErrorMsg[0] = err.Error()
			} else {
				log.Info("Install Choerodon succeed")
			}
			c.Metrics.Send()
		},
	}

	flags := cmd.PersistentFlags()
	addInstallFlags(flags, c)

	// set defaults from environment
	return cmd
}

func runInstallC7n(c *action.Choerodon) error {
	// 当 version 没有设置时，从 git repo 获取最新版本(本地的 config.yaml 也有配置 version ？)
	if c.Version == "" {
		c.Version = c7nutils.GetVersion(c7nconsts.DefaultGitBranch)
	}
	log.Infof("The current installing version is %s", c.Version)

	id, err := c.GetInstallDef(settings.ConfigFile, settings.ResourceFile)
	if err != nil {
		return errors.WithMessage(err, "Failed to get install configration file")
	}
	// TODO 当 repoUrl 优先级 flag -> config.yaml -> install.yaml -> default
	// 初始化 helmInstall
	// 只有 id 中用到了 RepoUrl
	if id.Spec.Basic.RepoURL != "" {
		c.RepoUrl = id.Spec.Basic.RepoURL
	} else {
		c.RepoUrl = c7nconsts.DefaultRepoUrl
	}
	c.DefaultAccessModes = id.DefaultAccessModes
	c.Slaver = &id.Spec.Basic.Slaver

	// 检查资源
	if err := c.CheckResource(&id.Spec.Resources); err != nil {
		return err
	}
	if err := c.CheckNamespace(c.Namespace); err != nil {
		return err
	}

	stopCh := make(chan struct{})
	_, err = c.PrepareSlaver(stopCh)
	if err != nil {
		return errors.WithMessage(err, "Create Slaver failed")
	}
	defer func() {
		stopCh <- struct{}{}
	}()

	// 渲染 Release
	if err := c.RenderReleases(id); err != nil {
		return nil
	}

	releaseGraph := resource.NewReleaseGraph(id.Spec.Release)
	installQueue := releaseGraph.TopoSortByKahn()

	for !installQueue.IsEmpty() {
		rls := installQueue.Dequeue()
		log.Infof("start install %s", rls.Name)
		// 获取的 values.yaml 必须经过渲染，只能放在 id 中
		vals, err := id.RenderHelmValues(rls, c.UserConfig)
		if err != nil {
			return err
		}
		if err = c.InstallRelease(rls, vals); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("Release %s install failed", rls.Name))
		}
	}
	// 等待所有 afterTask 执行完成。
	c.Wg.Wait()
	// c.SendMetrics(err)
	// 清理历史的job，cm，slaver 等
	return c.Clean()
}

func addInstallFlags(fs *pflag.FlagSet, client *action.Choerodon) {
	// moved to EnvSettings
	//fs.StringVarP(&client.ResourceFile, "resource-file", "r", "", "Resource file to read from, It provide which app should be installed")
	//fs.StringVarP(&client.ConfigFile, "c7n-config", "c", "", "User Config file to read from, User define config by this file")
	//fs.StringVarP(&client.Namespace, "namespace", "n", "c7n-system", "set namespace which install choerodon")

	fs.StringVar(&client.Version, "version", "", "specify a version")
	fs.StringVar(&client.Prefix, "prefix", "", "add prefix to all helm release")

	fs.BoolVar(&client.NoTimeout, "no-timeout", false, "disable resource job timeout")
	fs.BoolVar(&client.SkipInput, "skip-input", false, "use default username and password to avoid user input")
}

func setUserConfig(skipInput bool) {
	// 在 c7nctl.initConfig() 中 viper 获取了默认的配置文件
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Debug(err)
	}
	if !skipInput {
		// 当用户没有接受条款时，让其输入
		if !cfg.Terms.Accepted {
			c7nutils.AskAgreeTerms()
			mail := inputUserMail()
			cfg.Terms.Accepted = true
			cfg.OpsMail = mail
			viper.Set("terms", cfg.Terms)
			viper.Set("opsMail", cfg.OpsMail)
			viper.WriteConfig()
		}
	} else {
		log.Info("your are execute job by skip input option, so we think you had allowed we collect your information")
	}
}

func inputUserMail() string {
	mail, err := c7nutils.AcceptUserInput(c7nutils.Input{
		Password: false,
		Tip:      "请输入您的邮箱以便通知您重要的更新(Please enter your email address):  ",
		Regex:    "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	})
	c7nutils.CheckErr(err)
	return mail
}
