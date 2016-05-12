package config

import (
	"encoding/json"

	"github.com/frodenas/brokerapi"
)

//BrokerAPI provides the type for definition of broker API
type BrokerAPI struct {
	ExternalURL string                      `json:"external_url"`
	Listen      string                      `json:"listen"`
	Credentials brokerapi.BrokerCredentials `json:"credentials"`
}

//ManagementAPI provides the type for definition of management API
type ManagementAPI struct {
	Listen          string           `json:"listen"`
	DevMode         bool             `json:"dev_mode"`
	UaaClient       string           `json:"uaa_client"`
	UaaSecret       string           `json:"uaa_secret"`
	BrokerName      string           `json:"broker_name"`
	Authentication  *json.RawMessage `json:"authentication"`
	CloudController CloudController  `json:"cloud_controller"`
}

//Uaa provides the type to use for authentication and authorization
type Uaa struct {
	UaaAuth UaaAuth `json:"uaa"`
}

//UaaAuth provides authentication and authorization definition
type UaaAuth struct {
	Scope                    string `json:"adminscope"`
	PublicKey                string `json:"public_key"`
	SymmetricVerificationKey string `json:"symmetric_verification_key"`
}

//CloudController is the cloud controller definition
type CloudController struct {
	API               string `json:"api"`
	SkipTLSValidation bool   `json:"skip_tls_validation"`
}

//Dial is the dial definition. Dials should have a configuration and a serviceplan
type Dial struct {
	Configuration *json.RawMessage      `json:"configuration,omitempty"`
	Plan          brokerapi.ServicePlan `json:"plan"`
}

//Instance is definition of an Instance with the corresponding info
type Instance struct {
	TargetURL         string            `json:"target"`
	Name              string            `json:"name"`
	Dials             map[string]Dial   `json:"dials"`
	Service           brokerapi.Service `json:"service"`
	AuthenticationKey string            `json:"authentication_key"`
}

//RoutesRegister is the definition for RouteRegister
type RoutesRegister struct {
	NatsMembers      []string `json:"nats_members"`
	BrokerAPIHost    string   `json:"broker_api_host,omitempty"`
	ManagmentAPIHost string   `json:"management_api_host,omitempty"`
}

//Config is the configuration definition
type Config struct {
	APIVersion     string              `json:"api_version"`
	BrokerAPI      BrokerAPI           `json:"broker_api"`
	ManagementAPI  *ManagementAPI      `json:"management_api,omitempty"`
	Instances      map[string]Instance `json:"instances"`
	RoutesRegister *RoutesRegister     `json:"routes_register"`
}

//Provider is the definition for a config provider
type Provider interface {
	LoadConfiguration() (*Config, error)
	LoadDriverInstance(driverInstanceID string) (*Instance, error)
	GetUaaAuthConfig() (*UaaAuth, error)
	SetInstance(instanceid string, driverInstance Instance) error
	GetInstance(instanceid string) (instance *Instance, driverID string, err error)
	DeleteInstance(instanceid string) error
	SetService(instanceid string, service brokerapi.Service) error
	GetService(serviceid string) (service *brokerapi.Service, instanceid string, err error)
	DeleteService(instanceid string) error
	SetDial(instanceid string, dialid string, dial Dial) error
	GetDial(dialid string) (dial *Dial, instanceID string, err error)
	DeleteDial(dialid string) error
	InstanceNameExists(driverInstanceName string) (bool, error)
	GetPlan(plandid string) (plan *brokerapi.ServicePlan, dialid string, instanceid string, err error)
}
