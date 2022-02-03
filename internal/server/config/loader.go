package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_config/mock_$GOFILE
// Loader represents config Loader
type Loader interface {
	Load() error
	Get() *Config
}

type configLoader struct {
	configFilePath string
	config         *Config
}

func (c *configLoader) Get() *Config {
	return c.config
}

// NewConfigLoader returns an instance of Loader
func NewConfigLoader(configPath string) Loader {
	cfg := &Config{}
	return &configLoader{
		configFilePath: configPath,
		config:         cfg,
	}
}

// check implementation
var _ Loader = (*configLoader)(nil)

// Load loads usual application config together with secret config if defined.
// This method loads the app and secret config into the same structure
// so, be careful if both configs have overlapped sections and fields
func (c *configLoader) Load() error {
	// load config
	if c.configFilePath == "" {
		return errors.New("config path is not specified")
	}

	if err := parseYamlFile(c.configFilePath, &c.config); err != nil {
		return errors.Wrap(err, "can`t parse config file")
	}

	return nil
}

func parseYamlFile(fullName string, data interface{}) error {
	bytes, err := ioutil.ReadFile(fullName)
	if err != nil {
		return errors.Wrap(err, "can`t read file")
	}

	return yaml.Unmarshal(bytes, data)
}
