package config

import (
	"encoding/json"

	"github.com/pivotal-cf/brokerapi"
)

type BrokerAPI struct {
	ExternalUrl string                      `json:"external_url"`
	Listen      string                      `json:"listen"`
	DevMode     bool                        `json:"dev_mode"`
	Credentials brokerapi.BrokerCredentials `json:"credentials"`
}

type ManagementAPI struct {
	Listen          string           `json:"listen"`
	UaaClient       string           `json:"uaa_client"`
	UaaSecret       string           `json:"uaa_secret"`
	Authentication  *json.RawMessage `json:"authentication"`
	CloudController *json.RawMessage `json:"cloud_controller"`
}

type Uaa struct {
	UaaAuth UaaAuth `json:"uaa"`
}

type UaaAuth struct {
	Scope     string `json:"adminscope"`
	PublicKey string `json:"public_key"`
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
	GetUaaAuthConfig() (*UaaAuth, error)
	SetDriver(Driver) error
	GetDriver(string) (Driver, error)
	SetDriverInstance(string, DriverInstance) error
	GetDriverInstance(string) (DriverInstance, error)
	SetService(string, brokerapi.Service) error
	GetService(string) (brokerapi.Service, error)
	SetDial(string, Dial) error
	GetDial(string, string) (Dial, error)
}
