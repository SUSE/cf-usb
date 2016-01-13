package config

import (
	"encoding/json"

	"github.com/frodenas/brokerapi"
)

type BrokerAPI struct {
	ExternalUrl string                      `json:"external_url"`
	Listen      string                      `json:"listen"`
	Credentials brokerapi.BrokerCredentials `json:"credentials"`
}

type ManagementAPI struct {
	Listen          string           `json:"listen"`
	DevMode         bool             `json:"dev_mode"`
	UaaClient       string           `json:"uaa_client"`
	UaaSecret       string           `json:"uaa_secret"`
	Authentication  *json.RawMessage `json:"authentication"`
	CloudController CloudController  `json:"cloud_controller"`
}

type Uaa struct {
	UaaAuth UaaAuth `json:"uaa"`
}

type UaaAuth struct {
	Scope     string `json:"adminscope"`
	PublicKey string `json:"public_key"`
}

type CloudController struct {
	Api               string `json:"api"`
	SkipTslValidation bool   `json:"skip_tsl_validation"`
}

type Dial struct {
	Configuration *json.RawMessage      `json:"configuration,omitempty"`
	Plan          brokerapi.ServicePlan `json:"plan"`
}

type DriverInstance struct {
	Name          string            `json:"name"`
	Configuration *json.RawMessage  `json:"configuration"`
	Dials         map[string]Dial   `json:"dials"`
	Service       brokerapi.Service `json:"service"`
}

type Driver struct {
	DriverType      string                    `json:"driver_type"`
	DriverInstances map[string]DriverInstance `json:"driver_instances,omitempty"`
	DriverName      string                    `json:"driver_name"`
}

type RoutesRegister struct {
	NatsMembers      []string `json:"nats_members"`
	BrokerAPIHost    string   `json:"broker_api_host,omitempty"`
	ManagmentAPIHost string   `json:"management_api_host,omitempty"`
}

type Config struct {
	APIVersion     string            `json:"api_version"`
	BrokerAPI      BrokerAPI         `json:"broker_api"`
	ManagementAPI  *ManagementAPI    `json:"management_api,omitempty"`
	Drivers        map[string]Driver `json:"drivers"`
	RoutesRegister *RoutesRegister   `json:"routes_register"`
}

type ConfigProvider interface {
	LoadConfiguration() (*Config, error)
	LoadDriverInstance(driverInstanceID string) (*DriverInstance, error)
	GetUaaAuthConfig() (*UaaAuth, error)
	SetDriver(string, Driver) error
	GetDriver(string) (*Driver, error)
	DeleteDriver(string) error
	SetDriverInstance(string, string, DriverInstance) error
	GetDriverInstance(string) (*DriverInstance, error)
	DeleteDriverInstance(string) error
	SetService(string, brokerapi.Service) error
	GetService(string) (*brokerapi.Service, error)
	DeleteService(string) error
	SetDial(string, string, Dial) error
	GetDial(string, string) (*Dial, error)
	DeleteDial(string, string) error
	ServiceNameExists(string) (bool, error)
	DriverTypeExists(string) (bool, error)
}
