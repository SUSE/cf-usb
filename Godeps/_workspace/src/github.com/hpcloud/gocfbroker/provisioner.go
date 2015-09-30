package gocfbroker

// Provisioner represents a service managed by a v2 broker.
type Provisioner interface {
	// Provision a service instance.
	Provision(instanceID string, pr ProvisionRequest) (ProvisionResponse, error)
	// Deprovision a service instance.
	Deprovision(instanceID, serviceID, planID string) error
	// Update an instance of the service.
	Update(instanceID string, upr UpdateProvisionRequest) error
	// Bind an instance of the service.
	Bind(instanceID, bindingID string, br BindingRequest) (BindingResponse, error)
	// Unbind an instance the service.
	Unbind(instanceID, bindingID, serviceID, planID string) error
}
