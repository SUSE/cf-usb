package config

type DriverState struct {
	Instances []*ServiceInstance `json:"instances"`
}

type ServiceInstance struct {
	ID          string               `json:"id"`
	Credentials []*ServiceCredential `json:"credentials"`
}

type ServiceCredential struct {
	ID string `json:"id"`
}

type CredentialsConfig struct {
	StaticConfig string `json:"static_config"`
}
