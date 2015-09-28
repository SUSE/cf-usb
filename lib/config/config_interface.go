package config

import "github.com/hpcloud/gocfbroker"

type Config struct {
	BoltFilename string `json:"bolt_filename"`
	BoltBucket   string `json:"bolt_bucket"`

	gocfbroker.Options
}

type Configuration interface {
	LoadConfiguration() (Config, error)
}
