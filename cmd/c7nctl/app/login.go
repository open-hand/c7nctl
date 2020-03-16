package app

import (
	"fmt"
	"github.com/choerodon/c7nctl/pkg/authorize"
	"github.com/choerodon/c7nctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/url"
)

// todo: remove it
func Login(cmd *cobra.Command, args []string) error {

	cfg, err := utils.GetConfig()
	if err != nil {
		return err
	}
	auth := authorize.DefaultAuthorization(cfg)

	if len(args) == 1 {
		uri, err := url.Parse(args[0])
		if err != nil {
			return err
		}
		auth.ServerUrl = fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)
	}

	if c := cfg.FindNamedClusterByServer(auth.ServerUrl); c != nil {
		auth.ClusterName = c.Name
	}
	if clusterName, err := cmd.Flags().GetString("name"); err != nil {
		return err
	} else if clusterName != "" {
		auth.ClusterName = clusterName
	}
	if auth.ClusterName == "" {
		auth.ClusterName = utils.RandomString()
	}
	if auth.ServerUrl == "" {
		return errors.New("Mush specify a Choerodon api server ")
	}
	return auth.Login()
}
