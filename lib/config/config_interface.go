package config

import (
	"encoding/json"

	"github.com/pivotal-cf/brokerapi"
)

type BrokerAPI struct {
	Listen      string                      `json:"listen"`
	Credentials brokerapi.BrokerCredentials `json:"credentials"`
}

type ManagementAPI struct {
	Listen         string           `json:"listen"`
	UaaClient      string           `json:"uaa_client"`
	UaaSecret      string           `json:"uaa_secret"`
	Authentication *json.RawMessage `json:"authentication"`
}

type Dial struct {
	ID            string                `json:"id"`
	Configuration *json.RawMessage      `json:"configuration,omitempty"`
	Plan          brokerapi.ServicePlan `json:"plan"`
}
type DriverInstance struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Configuration *json.RawMessage  `json:"configuration"`
	Dials         []Dial            `json:"dials"`
	Service       brokerapi.Service `json:"service"`
}

type Driver struct {
	ID              string            `json:"id"`
	DriverType      string            `json:"driver_type"`
	DriverInstances []*DriverInstance `json:"driver_instances,omitempty"`
}

type Config struct {
	APIVersion    string         `json:"api_version"`
	LogLevel      string         `json:"logLevel"`
	BrokerAPI     BrokerAPI      `json:"broker_api"`
	ManagementAPI *ManagementAPI `json:"management_api,omitempty"`
	Drivers       []Driver       `json:"drivers"`
}

type ConfigProvider interface {
	LoadConfiguration() (*Config, error)
	GetDriverInstanceConfig(driverInstanceID string) (*DriverInstance, error)
}
