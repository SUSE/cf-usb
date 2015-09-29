package config

import (
	"encoding/json"

	"github.com/hpcloud/gocfbroker"
)

type DriverConfig struct {
	DriverType    string          `json:"driver_type"`
	Configuration json.RawMessage `json:"configuration"`
	ServiceIDs    []string        `json:"service_ids"`
}

type DriverProperties struct {
	DriverConfiguration json.RawMessage
	Services            []gocfbroker.Service
}

type Config struct {
	BoltFilename  string         `json:"bolt_filename"`
	BoltBucket    string         `json:"bolt_bucket"`
	DriverConfigs []DriverConfig `json:"driver_configs"`

	gocfbroker.Options
}

type ConfigProvider interface {
	LoadConfiguration() (Config, error)
	GetDriverProperties(driverType string) (DriverProperties, error)
	GetDriverTypes() ([]string, error)
}
