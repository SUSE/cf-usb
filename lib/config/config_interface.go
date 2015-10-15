package config

import (
	"encoding/json"
	"io"

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
	Output              io.Writer
}

type Config struct {
	Crednetials      brokerapi.BrokerCredentials `json:"broker_credentials"`
	ServiceCatalog   []brokerapi.Service         `json:"services"`
	DriverConfigs    []DriverConfig              `json:"driver_configs"`
	Listen           string                      `json:"listen"`
	ManagementListen string                      `json:"management_listen"`
	StartMgmt        bool                        `json:"start_mgmt"`
	APIVersion       string                      `json:"api_version"`
	LogLevel         string                      `json:"logLevel"`
}

type ConfigProvider interface {
	LoadConfiguration() (*Config, error)
	GetDriverProperties(driverType string) (DriverProperties, error)
	GetDriverTypes() ([]string, error)
}
