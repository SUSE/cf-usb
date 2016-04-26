package lib

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/pivotal-golang/lager"

<<<<<<< HEAD
	httptransport "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
=======
	httptransport "github.com/go-swagger/go-swagger/httpkit/client"

	strfmt "github.com/go-swagger/go-swagger/strfmt"
>>>>>>> f998b3c... [HCFRO-193] Use rest for calling drivers

	"github.com/hpcloud/cf-usb/lib/servicemgr"
	models "github.com/hpcloud/cf-usb/lib/servicemgr/models"
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

	for _, instance := range config.DriverInstances {
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

	driver, err := broker.getDriver(serviceDetails.ServiceID)
	if err != nil {
		return brokerapi.ProvisioningResponse{}, false, err
	}
	if driver == nil {
		return brokerapi.ProvisioningResponse{}, false, errors.New(fmt.Sprintf("Cannot find driver for %s", serviceDetails.ServiceID))
	}

	instance, _ := driver.GetWorkspace(instanceID)
	if instance.Status != nil {
		statusInfo := instance.Status
		if *statusInfo == "failed" {
			return brokerapi.ProvisioningResponse{}, false, brokerapi.ErrInstanceAlreadyExists
		}
	}
	request := models.ServiceManagerWorkspaceCreateRequest{}
	request.WorkspaceID = &instanceID
	request.Details = serviceDetails.Parameters
	instance, errorDetails := driver.CreateWorkspace(request)
	if errorDetails.Message != nil {
<<<<<<< HEAD
		broker.logger.Error("provision-instance", nil, lager.Data{"message": errorDetails.Message, "code": errorDetails.Code})
=======
>>>>>>> f998b3c... [HCFRO-193] Use rest for calling drivers
		return brokerapi.ProvisioningResponse{}, false, errors.New(*errorDetails.Message)
	}

	broker.logger.Info("provision-instance-request-completed", lager.Data{"instance-id": instanceID, "status": instance.Status, "details": instance.Details})
	if instance.Status != nil {
		statusInfo := *instance.Status
		if statusInfo == "successful" {
			return brokerapi.ProvisioningResponse{}, false, nil
		}
	}
	if instance.Status != nil {
		statusInfo := *instance.Status
		if statusInfo == "unknown" {
			if !acceptsIncomplete {
				// TODO: clean up instance
				// driver.DeprovisionInstance(instanceID)
				return brokerapi.ProvisioningResponse{}, false, brokerapi.ErrAsyncRequired
			}

			return brokerapi.ProvisioningResponse{}, true, nil
		}
	}

	return brokerapi.ProvisioningResponse{}, false, errors.New("Unknown instance state")
}

func (broker *UsbBroker) Update(instanceID string, details brokerapi.UpdateDetails, acceptsIncomplete bool) (bool, error) {
	return false, brokerapi.ErrInstanceNotUpdateable
}

func (broker *UsbBroker) Deprovision(instanceID string, deprovisionDetails brokerapi.DeprovisionDetails, acceptsIncomplete bool) (bool, error) {
	broker.logger.Info("deprovision-instance-request", lager.Data{"instance-id": instanceID, "service-id": deprovisionDetails.ServiceID, "accept-incomplete": acceptsIncomplete})

	driver, err := broker.getDriver(deprovisionDetails.ServiceID)
	if err != nil {
		return false, err
	}
	if driver == nil {
		return false, errors.New(fmt.Sprintf("Cannot find driver for %s", deprovisionDetails.ServiceID))
	}

	instance, _ := driver.GetWorkspace(instanceID)
	if instance.Status != nil {
		statusInfo := instance.Status
		if *statusInfo == "none" {
			return false, brokerapi.ErrInstanceDoesNotExist
		}
	}

	errorDetails := driver.DeleteWorkspace(instanceID)
	if errorDetails.Message != nil {
<<<<<<< HEAD
		broker.logger.Error("deprovision-instance", nil, lager.Data{"message": errorDetails.Message, "code": errorDetails.Code})

=======
>>>>>>> f998b3c... [HCFRO-193] Use rest for calling drivers
		return false, errors.New(*errorDetails.Message)
	} else {
		return false, nil
	}

	broker.logger.Info("deprovision-instance-request-completed", lager.Data{"instance-id": instanceID, "status": instance.Status, "details": instance.Details})

	//  TODO:
	//  Is InProgress applicabale to deprovision?
	//  Should there be another state, e.g. status.DeprovisionInProgress vs status.ProvisionInProgrss
	if instance.Status != nil {
		statusInfo := *instance.Status
		if statusInfo == "unknown" {
			return true, nil
		}
	}

	return false, errors.New("Unknown instance state")
}

func (broker *UsbBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.BindingResponse, error) {
	broker.logger.Info("generate-credentials-request", lager.Data{"instance-id": instanceID, "binding-id": bindingID, "service-id": details.ServiceID})

	var response brokerapi.BindingResponse

	driver, err := broker.getDriver(details.ServiceID)
	if err != nil {
		return response, err
	}
	request := models.ServiceManagerConnectionCreateRequest{}
	request.ConnectionID = &bindingID
	request.Details = structs.Map(details)
	credentials, errorDetails := driver.CreateWorkspaceConnection(instanceID, request)
	if errorDetails.Message != nil {
<<<<<<< HEAD
		broker.logger.Error("bind", nil, lager.Data{"message": errorDetails.Message, "code": errorDetails.Code})

=======
>>>>>>> f998b3c... [HCFRO-193] Use rest for calling drivers
		return response, errors.New(*errorDetails.Message)
	}
	if *credentials.Status == "failed" {
		return response, brokerapi.ErrBindingAlreadyExists
	}

	response.Credentials = credentials.Details

	broker.logger.Info("generate-credentials-request-completed", lager.Data{"instance-id": instanceID})

	return response, nil
}

func (broker *UsbBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	broker.logger.Info("revoke-credentials-request", lager.Data{"instance-id": instanceID, "binding-id": bindingID, "service-id": details.ServiceID})

	driver, err := broker.getDriver(details.ServiceID)
	if err != nil {
		return err
	}

	errorDetails := driver.DeleteWorkspaceConnection(instanceID, bindingID)
	if errorDetails.Message != nil {
<<<<<<< HEAD
		broker.logger.Error("unbind", nil, lager.Data{"message": errorDetails.Message, "code": errorDetails.Code})

=======
>>>>>>> f998b3c... [HCFRO-193] Use rest for calling drivers
		return errors.New(*errorDetails.Message)
	}

	broker.logger.Info("revoke-credentials-request-completed", lager.Data{"driver-response": errorDetails.Message})

	return nil
}

func (broker *UsbBroker) LastOperation(instanceID string) (brokerapi.LastOperationResponse, error) {
	broker.logger.Info("last-operation-request", lager.Data{"instance-id": instanceID})

	// TODO: how to get the driver for a instanceID. NOTE: the broker API does not require
	// the client to inclide the serviceID in the request
	driver, driverFound, err := broker.getDriverForServiceInstanceId(instanceID)
	if err != nil {
		return brokerapi.LastOperationResponse{}, err
	}
	if !driverFound {
		return brokerapi.LastOperationResponse{}, brokerapi.ErrInstanceDoesNotExist
	}

	instance, errorDetails := driver.GetWorkspace(instanceID)
	if errorDetails.Message != nil {
<<<<<<< HEAD
		broker.logger.Error("last-operation", nil, lager.Data{"message": errorDetails.Message, "code": errorDetails.Code})

=======
>>>>>>> f998b3c... [HCFRO-193] Use rest for calling drivers
		return brokerapi.LastOperationResponse{}, errors.New(*errorDetails.Message)
	}
	if instance.Status != nil {
		statusInfo := *instance.Status
		if statusInfo == "none" {
			return brokerapi.LastOperationResponse{}, brokerapi.ErrInstanceDoesNotExist
		}
	}
	if instance.Status != nil {
		statusInfo := *instance.Status
		if statusInfo == "successful" {
			return brokerapi.LastOperationResponse{State: brokerapi.LastOperationSucceeded}, nil
		}
	}
	if instance.Status != nil {
		statusInfo := *instance.Status
		if statusInfo == "unknown" {
			return brokerapi.LastOperationResponse{State: brokerapi.LastOperationInProgress}, nil
		}
	}
	if instance.Status != nil {
		statusInfo := *instance.Status
		if statusInfo == "failed" {
			return brokerapi.LastOperationResponse{State: brokerapi.LastOperationFailed}, nil
		}
	}
	// TODO: what about instance.Status == status.Deleted ?
	return brokerapi.LastOperationResponse{}, errors.New("Unknown instance state")
}

func (broker *UsbBroker) getDriver(serviceID string) (servicemgr.ServiceManagerInterface, error) {
	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	for _, driverInstance := range config.DriverInstances {
		if driverInstance.Service.ID == serviceID {
			if driverInstance.TargetURL != "" {
				u, err := url.Parse(driverInstance.TargetURL)
				if err != nil {
					return nil, err
				}

				transport := httptransport.New(u.Host, "/", []string{u.Scheme})

				debug, _ := strconv.ParseBool(os.Getenv("CF_TRACE"))

				transport.Debug = debug

				serviceManager := servicemgr.NewServiceManager(transport, strfmt.Default, broker.logger)
				return serviceManager, nil
			}
		}

	}

	return nil, errors.New("Driver not found")
}

func (broker *UsbBroker) getDriverForServiceInstanceId(instanceID string) (servicemgr.ServiceManagerInterface, bool, error) {
	config, err := broker.configProvider.LoadConfiguration()
	if err != nil {
		return nil, false, err
	}

	for _, driverInstance := range config.DriverInstances {
		u, err := url.Parse(driverInstance.TargetURL)
		if err != nil {
			return nil, false, err
		}

		transport := httptransport.New(u.Host, "/", []string{u.Scheme})

		debug, _ := strconv.ParseBool(os.Getenv("CF_TRACE"))

		transport.Debug = debug

		serviceManager := servicemgr.NewServiceManager(transport, strfmt.Default, broker.logger)

		instance, errorDetails := serviceManager.GetWorkspace(instanceID)
		if errorDetails.Message != nil {
			return nil, false, errors.New(*errorDetails.Message)
		}
		if instance.Status != nil {
			statusInfo := *instance.Status
			if statusInfo == "none" {
				continue
			}
		}

		return serviceManager, true, nil
	}

	return nil, false, nil
}
