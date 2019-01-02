package common

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
	Version string
	Terms   Terms
	User    User
}

type Terms struct {
	Accepted bool
}

type User struct {
	Mail string
}

func newConfig() *Config {
	return &Config{
		Version: "v1",
		Terms: Terms{
			Accepted: false,
		},
		User: User{
			Mail: "",
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
