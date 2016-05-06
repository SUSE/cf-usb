package lib

import (
	"errors"
	"net/url"

	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/pivotal-golang/lager"
)

var usbBroker UsbBroker

type UsbBroker struct {
	configProvider config.ConfigProvider
	logger         lager.Logger
}

func NewUsbBroker(configProvider config.ConfigProvider, logger lager.Logger) *UsbBroker {
	return &UsbBroker{configProvider: configProvider, logger: logger.Session("usb-broker")}
}

func (broker *UsbBroker) Services() brokerapi.CatalogResponse {
	broker.logger.Info("get-catalog-request", lager.Data{})

	var catalog []brokerapi.Service

	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		broker.logger.Fatal("retrive-configuration-failed", err)
	}

	for _, instance := range config.Instances {
		service := instance.Service

		for _, dial := range instance.Dials {
			service.Plans = append(service.Plans, dial.Plan)
		}

		catalog = append(catalog, service)
	}

	broker.logger.Info("get-catalog-request-completed", lager.Data{"services-found": len(catalog)})

	return brokerapi.CatalogResponse{Services: catalog}
}

func (broker *UsbBroker) Provision(instanceID string, serviceDetails brokerapi.ProvisionDetails, acceptsIncomplete bool) (brokerapi.ProvisioningResponse, bool, error) {
	broker.logger.Info("provision-instance-request", lager.Data{"instance-id": instanceID, "service-id": serviceDetails.ServiceID, "accept-incomplete": acceptsIncomplete})

	cmsClient, err := broker.getCSMClient(serviceDetails.ServiceID)
	if err != nil {
		return brokerapi.ProvisioningResponse{}, false, err
	}

	exists, err := cmsClient.WorkspaceExists(instanceID)
	if err != nil {
		return brokerapi.ProvisioningResponse{}, false, err
	}

	if exists {
		return brokerapi.ProvisioningResponse{}, false, brokerapi.ErrInstanceAlreadyExists
	}

	err = cmsClient.CreateWorkspace(instanceID)
	if err != nil {
		return brokerapi.ProvisioningResponse{}, false, err
	}
	broker.logger.Info("provision-instance-request-completed", lager.Data{"instance-id": instanceID})

	//TODO: wait for async operations in CSM
	return brokerapi.ProvisioningResponse{}, false, nil
}

func (broker *UsbBroker) Update(instanceID string, details brokerapi.UpdateDetails, acceptsIncomplete bool) (bool, error) {
	return false, brokerapi.ErrInstanceNotUpdateable
}

func (broker *UsbBroker) Deprovision(instanceID string, deprovisionDetails brokerapi.DeprovisionDetails, acceptsIncomplete bool) (bool, error) {
	broker.logger.Info("deprovision-instance-request", lager.Data{"instance-id": instanceID, "service-id": deprovisionDetails.ServiceID, "accept-incomplete": acceptsIncomplete})

	csmClient, err := broker.getCSMClient(deprovisionDetails.ServiceID)
	if err != nil {
		return false, err
	}

	exists, err := csmClient.WorkspaceExists(instanceID)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, brokerapi.ErrInstanceDoesNotExist
	}

	err = csmClient.DeleteWorkspace(instanceID)
	if err != nil {
		return false, err
	}

	//TODO: implement async
	return false, nil
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.BindingResponse, error) {
	broker.logger.Info("generate-credentials-request", lager.Data{"instance-id": instanceID, "binding-id": bindingID, "service-id": details.ServiceID})

	var response brokerapi.BindingResponse

	csmClient, err := broker.getCSMClient(details.ServiceID)
	if err != nil {
		return response, err
	}

	exists, err := csmClient.ConnectionExists(instanceID, bindingID)
	if err != nil {
		return response, err
	}

	if exists {
		return response, brokerapi.ErrBindingAlreadyExists
	}

	credentials, err := csmClient.CreateConnection(instanceID, bindingID)
	if err != nil {
		return response, err
	}

	response.Credentials = credentials

	broker.logger.Info("generate-credentials-request-completed", lager.Data{"instance-id": instanceID})

	return response, nil
}

func (broker *UsbBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	broker.logger.Info("revoke-credentials-request", lager.Data{"instance-id": instanceID, "binding-id": bindingID, "service-id": details.ServiceID})

	csmClient, err := broker.getCSMClient(details.ServiceID)
	if err != nil {
		return err
	}

	exists, err := csmClient.ConnectionExists(instanceID, bindingID)
	if err != nil {
		return err
	}

	if !exists {
		return brokerapi.ErrBindingDoesNotExist
	}

	broker.logger.Info("revoke-credentials-request-completed", lager.Data{})

	return nil
}

func (broker *UsbBroker) LastOperation(instanceID string) (brokerapi.LastOperationResponse, error) {
	broker.logger.Info("last-operation-request", lager.Data{"instance-id": instanceID})

	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		return brokerapi.LastOperationResponse{}, err
	}
	//TODO: add async

	instance := config.Instances[instanceID]

	instanceURL, err := url.Parse(instance.TargetURL)
	if err != nil {
		return brokerapi.LastOperationResponse{}, err
	}

	csmClient := csm.NewCSMClient(broker.logger, *instanceURL, instance.AuthenticationKey)

	exists, err := csmClient.WorkspaceExists(instanceID)
	if err != nil {
		return brokerapi.LastOperationResponse{}, err
	}

	if exists {
		return brokerapi.LastOperationResponse{State: brokerapi.LastOperationSucceeded}, nil
	}
	return brokerapi.LastOperationResponse{}, errors.New("Unknown instance state")

}

func (broker *UsbBroker) getCSMClient(serviceID string) (csm.CSMInterface, error) {
	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	for _, driverInstance := range config.Instances {
		if driverInstance.Service.ID == serviceID {
			if driverInstance.TargetURL != "" {
				u, err := url.Parse(driverInstance.TargetURL)
				if err != nil {
					return nil, err
				}
				client := csm.NewCSMClient(broker.logger, *u, driverInstance.AuthenticationKey)
				return client, nil
			}
		}

	}

	return nil, errors.New("Instance not found")
}
