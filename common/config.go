package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env  string `yaml:"env"`
	Port string `yaml:"port"`
	Bind string `yaml:"bind"`

	GoBookStore struct {
		DB      string `yaml:"db"`
		URI     string `yaml:"dsn"`
		LOGPATH string `yaml:"logpath"`
	} `yaml:"gobookstore"`
}

// Load the yaml file into the Config struct
func LoadConfig(cfgFile string) (conf *Config, err error) {

	conf = &Config{}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		return
	}
	return
}
