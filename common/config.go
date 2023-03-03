package common

import (
	"os"
	"time"

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
		CORS    struct {
			ALLOWED_ORIGINS   []string      `yaml:"allowed_origins"`
			ALLOWED_METHOS    []string      `yaml:"allowed_methods"`
			ALLOWED_HEADERS   []string      `yaml:"allowed_headers"`
			ALLOW_CREDENTIALS bool          `yaml:"allow_credentials"`
			EXPOSED_HEADERS   []string      `yaml:"exposed_headers"`
			MAX_AGE           time.Duration `yaml:"max_duration"`
		}
	}
}

// Load the yaml file into the Config struct
func LoadConfig(cfgFile string) (conf *Config, err error) {

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
