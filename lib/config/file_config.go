package config

import "github.com/hpcloud/gocfbroker"

type fileConfig struct {
	Configuration,
	path string
}

func NewFileConfig(path string) Configuration {
	return &fileConfig{path: path}
}

func (configuration *fileConfig) LoadConfiguration() (Config, error) {
	var config Config

	err := gocfbroker.LoadConfig(configuration.path, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
