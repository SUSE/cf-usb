package driver

type Driver interface {
	Provision(interface{}, *interface{}) error
	Deprovision(string, *string) error
	Bind(string, *string) error
	Unbind(string, *string) error
	Update(string, *string) error
	GetCatalog(string, *string) error
	GetInstances(string, *string) error
}
