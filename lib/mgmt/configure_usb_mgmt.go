package mgmt

import (
	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"

	"github.com/hpcloud/cf-usb/lib/genmodel"
	. "github.com/hpcloud/cf-usb/lib/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

func ConfigureAPI(api *UsbMgmtAPI) {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.DeleteDriverConfigHandler = DeleteDriverConfigHandlerFunc(func() error {
		return errors.NotImplemented("operation deleteDriverConfig has not yet been implemented")
	})

	api.GetServicePlanByIDHandler = GetServicePlanByIDHandlerFunc(func() (*genmodel.Plan, error) {
		return nil, errors.NotImplemented("operation getServicePlanByID has not yet been implemented")
	})

	api.CreateDriverConfigHandler = CreateDriverConfigHandlerFunc(func(params CreateDriverConfigParams) error {
		return errors.NotImplemented("operation createDriverConfig has not yet been implemented")
	})

	api.DeleteServicePlanHandler = DeleteServicePlanHandlerFunc(func() error {
		return errors.NotImplemented("operation deleteServicePlan has not yet been implemented")
	})

	api.GetServicePlansHandler = GetServicePlansHandlerFunc(func() (*[]genmodel.Plan, error) {
		return nil, errors.NotImplemented("operation getServicePlans has not yet been implemented")
	})

	api.GetDriversHandler = GetDriversHandlerFunc(func() (*[]string, error) {
		return nil, errors.NotImplemented("operation getDrivers has not yet been implemented")
	})

	api.UpdateServicePlanHandler = UpdateServicePlanHandlerFunc(func(params UpdateServicePlanParams) error {
		return errors.NotImplemented("operation updateServicePlan has not yet been implemented")
	})

	api.GetServicesHandler = GetServicesHandlerFunc(func() (*[]genmodel.Service, error) {
		return nil, errors.NotImplemented("operation getServices has not yet been implemented")
	})

	api.GetInfoHandler = GetInfoHandlerFunc(func() (*genmodel.Info, error) {

		return &genmodel.Info{Version: "0.0.1"}, nil
	})

	api.GetDriverConfigByIDHandler = GetDriverConfigByIDHandlerFunc(func() (*genmodel.DriverConfig, error) {
		return nil, errors.NotImplemented("operation getDriverConfigByID has not yet been implemented")
	})

	api.CreateServicePlanHandler = CreateServicePlanHandlerFunc(func(params CreateServicePlanParams) error {
		return errors.NotImplemented("operation createServicePlan has not yet been implemented")
	})

	api.UpdateServiceHandler = UpdateServiceHandlerFunc(func(params UpdateServiceParams) error {
		return errors.NotImplemented("operation updateService has not yet been implemented")
	})

	api.DeleteServiceHandler = DeleteServiceHandlerFunc(func() error {
		return errors.NotImplemented("operation deleteService has not yet been implemented")
	})

	api.GetServiceByIDHandler = GetServiceByIDHandlerFunc(func() (*genmodel.Service, error) {
		return nil, errors.NotImplemented("operation getServiceByID has not yet been implemented")
	})

	api.GetDriverConfigsHandler = GetDriverConfigsHandlerFunc(func() (*[]genmodel.DriverConfig, error) {
		return nil, errors.NotImplemented("operation getDriverConfigs has not yet been implemented")
	})

	api.CreateServiceHandler = CreateServiceHandlerFunc(func(params CreateServiceParams) error {
		return errors.NotImplemented("operation createService has not yet been implemented")
	})

	api.UpdateDriverConfigHandler = UpdateDriverConfigHandlerFunc(func(params UpdateDriverConfigParams) error {
		return errors.NotImplemented("operation updateDriverConfig has not yet been implemented")
	})

}
