package main

import (
	"github.com/choerodon/c7nctl/pkg/resource"
	"github.com/spf13/pflag"
)

func addResourceClientFlags(fs *pflag.FlagSet, client *resource.Client) {
	fs.BoolVar(&client.Business, "biz", false, "enable install business choerodon")
	fs.StringVar(&client.Username, "auth-user", "", "The authenticated user of the installation resource")
	fs.StringVar(&client.Password, "auth-pass", "", "The authenticated password of the installation resource")
}
