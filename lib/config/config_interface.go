package config

import (
	"encoding/json"

	"github.com/pivotal-cf/brokerapi"
)

type DriverConfig struct {
	DriverType    string           `json:"driver_type"`
	Configuration *json.RawMessage `json:"configuration"`
	ServiceIDs    []string         `json:"service_ids"`
}

type DriverProperties struct {
	DriverConfiguration *json.RawMessage
	Services            []brokerapi.Service
}

type Config struct {
	Crednetials    brokerapi.BrokerCredentials `json:"broker_credentials"`
	ServiceCatalog []brokerapi.Service         `json:"services"`
	DriverConfigs  []DriverConfig              `json:"driver_configs"`
	Listen         string                      `json:"listen"`
	APIVersion     string                      `json:"api_version"`
	LogLevel       string                      `json:logLevel`
}

type ConfigProvider interface {
	LoadConfiguration() (*Config, error)
	GetDriverProperties(driverType string) (DriverProperties, error)
	GetDriverTypes() ([]string, error)
}
