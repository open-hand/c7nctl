package utils

import (
	"k8s.io/client-go/util/homedir"
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	configFileName = "config.yml"
)

type Config struct {
	Version         string
	Terms           Terms
	Clusters        []*NamedCluster
	SelectedCluster string `yaml:"current-cluster"`
	Users           []*NamedUser
	OpsMail         string `yaml:"opsMail"`
}

type NamedUser struct {
	Name string `yaml:"name"`
	User *User  `yaml:"user"`
}

type NamedCluster struct {
	Name    string   `yaml:"name"`
	Cluster *Cluster `yaml:"cluster"`
}

type Cluster struct {
	Server       string
	User         *User  `yaml:"-"`
	SelectedUser string `yaml:"user"`
}

type Terms struct {
	Accepted bool
}

type User struct {
	Mail       string
	Token      string
	Name       string
	ValidUntil int64
}

func newConfig() *Config {
	return &Config{
		Version: "v1",
		Terms: Terms{
			Accepted: false,
		},
	}
}

func getPath() string {
	homePath := homedir.HomeDir()
	configDir := fmt.Sprintf("%s%c.c7n%c", homePath, os.PathSeparator, os.PathSeparator)
	// check dir is exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, os.FileMode(0755))
	}
	configFile := configDir + configFileName
	return configFile
}

func GetConfig() (*Config, error) {

	c := newConfig()

	configFile := getPath()
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		c.Write()
		return c, nil
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return c, err
	}

	yaml.Unmarshal(data, c)
	return c, nil
}

func (c *Config) Write() error {
	configFile := getPath()
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, data, os.FileMode(0755))
}

func (c *Config) CurrentCluster() *Cluster {
	for _, v := range c.Clusters {
		if v.Name == c.SelectedCluster {
			v.Cluster.User = c.findUserByName(v.Cluster.SelectedUser)
			return v.Cluster
		}
	}
	return nil
}

func (c *Config) findUserByName(name string) *User {
	for _, v := range c.Users {
		if v.Name == name {
			return v.User
		}
	}
	return nil
}

func (c *Config) FindNamedClusterByServer(serverUrl string) *NamedCluster {
	for _, v := range c.Clusters {
		if v.Cluster.Server == serverUrl {
			v.Cluster.User = c.findUserByName(v.Cluster.SelectedUser)
			return v
		}
	}
	return nil
}

func (c *Config) CurrentUser() *User {
	if c.CurrentCluster() == nil || c.CurrentCluster().User == nil {
		return nil
	}
	return c.CurrentCluster().User
}

func (c *Config) CurrentServer() string {
	if c.CurrentCluster() == nil {
		return ""
	}
	return c.CurrentCluster().Server
}
