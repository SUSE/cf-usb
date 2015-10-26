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

	api.UpdateDriverHandler = UpdateDriverHandlerFunc(func(params UpdateDriverParams) error {
		return errors.NotImplemented("operation updateDriver has not yet been implemented")
	})

	api.UploadDriverHandler = UploadDriverHandlerFunc(func(params UploadDriverParams) error {
		return errors.NotImplemented("operation uploadDriver has not yet been implemented")
	})

	api.DeleteServicePlanHandler = DeleteServicePlanHandlerFunc(func(params DeleteServicePlanParams) error {
		return errors.NotImplemented("operation deleteServicePlan has not yet been implemented")
	})

	api.GetDriverSchemaHandler = GetDriverSchemaHandlerFunc(func(params GetDriverSchemaParams) (*genmodel.DriverSchema, error) {
		return nil, errors.NotImplemented("operation getDriverSchema has not yet been implemented")
	})

	api.GetDriverHandler = GetDriverHandlerFunc(func(params GetDriverParams) (*genmodel.Driver, error) {
		return nil, errors.NotImplemented("operation getDriver has not yet been implemented")
	})

	api.GetDriversHandler = GetDriversHandlerFunc(func(params GetDriversParams) (*[]genmodel.Driver, error) {
		return nil, errors.NotImplemented("operation getDrivers has not yet been implemented")
	})

	api.CreateDriverHandler = CreateDriverHandlerFunc(func(params CreateDriverParams) error {
		return errors.NotImplemented("operation createDriver has not yet been implemented")
	})

	api.UpdateServiceHandler = UpdateServiceHandlerFunc(func(params UpdateServiceParams) error {
		return errors.NotImplemented("operation updateService has not yet been implemented")
	})

	api.GetDriverInstancesHandler = GetDriverInstancesHandlerFunc(func(params GetDriverInstancesParams) (*[]genmodel.DriverInstance, error) {
		return nil, errors.NotImplemented("operation getDriverInstances has not yet been implemented")
	})

	api.PingDriverInstanceHandler = PingDriverInstanceHandlerFunc(func(params PingDriverInstanceParams) error {
		return errors.NotImplemented("operation pingDriverInstance has not yet been implemented")
	})

	api.GetServicePlanHandler = GetServicePlanHandlerFunc(func(params GetServicePlanParams) (*genmodel.Plan, error) {
		return nil, errors.NotImplemented("operation getServicePlan has not yet been implemented")
	})

	api.CreateDialHandler = CreateDialHandlerFunc(func(params CreateDialParams) error {
		return errors.NotImplemented("operation createDial has not yet been implemented")
	})

	api.UpdateDialHandler = UpdateDialHandlerFunc(func(params UpdateDialParams) error {
		return errors.NotImplemented("operation updateDial has not yet been implemented")
	})

	api.DeleteServiceHandler = DeleteServiceHandlerFunc(func(params DeleteServiceParams) error {
		return errors.NotImplemented("operation deleteService has not yet been implemented")
	})

	api.GetDialSchemaHandler = GetDialSchemaHandlerFunc(func(params GetDialSchemaParams) (*genmodel.DialSchema, error) {
		return nil, errors.NotImplemented("operation getDialSchema has not yet been implemented")
	})

	api.GetServicesHandler = GetServicesHandlerFunc(func(params GetServicesParams) (*[]genmodel.Service, error) {
		return nil, errors.NotImplemented("operation getServices has not yet been implemented")
	})

	api.CreateDriverInstanceHandler = CreateDriverInstanceHandlerFunc(func(params CreateDriverInstanceParams) (*genmodel.DriverInstance, error) {
		return nil, errors.NotImplemented("operation createDriverInstance has not yet been implemented")
	})

	api.UpdateDriverInstanceHandler = UpdateDriverInstanceHandlerFunc(func(params UpdateDriverInstanceParams) error {
		return errors.NotImplemented("operation updateDriverInstance has not yet been implemented")
	})

	api.DeleteDriverHandler = DeleteDriverHandlerFunc(func(params DeleteDriverParams) error {
		return errors.NotImplemented("operation deleteDriver has not yet been implemented")
	})

	api.DeleteDriverInstanceHandler = DeleteDriverInstanceHandlerFunc(func(params DeleteDriverInstanceParams) error {
		return errors.NotImplemented("operation deleteDriverInstance has not yet been implemented")
	})

	api.GetAllDialsHandler = GetAllDialsHandlerFunc(func(params GetAllDialsParams) (*[]genmodel.Dial, error) {
		return nil, errors.NotImplemented("operation getAllDials has not yet been implemented")
	})

	api.UpdateServicePlanHandler = UpdateServicePlanHandlerFunc(func(params UpdateServicePlanParams) error {
		return errors.NotImplemented("operation updateServicePlan has not yet been implemented")
	})

	api.GetBrokerInfoHandler = GetBrokerInfoHandlerFunc(func(params GetBrokerInfoParams) (*genmodel.Broker, error) {
		return nil, errors.NotImplemented("operation getBrokerInfo has not yet been implemented")
	})

	api.GetServiceHandler = GetServiceHandlerFunc(func(params GetServiceParams) (*genmodel.Service, error) {
		return nil, errors.NotImplemented("operation getService has not yet been implemented")
	})

	api.UpdateBrokerInfoHandler = UpdateBrokerInfoHandlerFunc(func(params UpdateBrokerInfoParams) (*genmodel.Broker, error) {
		return nil, errors.NotImplemented("operation updateBrokerInfo has not yet been implemented")
	})

	api.DeleteDialHandler = DeleteDialHandlerFunc(func(params DeleteDialParams) error {
		return errors.NotImplemented("operation deleteDial has not yet been implemented")
	})

	api.GetInfoHandler = GetInfoHandlerFunc(func(params GetInfoParams) (*genmodel.Info, error) {
		return nil, errors.NotImplemented("operation getInfo has not yet been implemented")
	})

	api.GetServicePlansHandler = GetServicePlansHandlerFunc(func(params GetServicePlansParams) (*[]genmodel.Plan, error) {
		return nil, errors.NotImplemented("operation getServicePlans has not yet been implemented")
	})

	api.GetDialHandler = GetDialHandlerFunc(func(params GetDialParams) (*genmodel.Dial, error) {
		return nil, errors.NotImplemented("operation getDial has not yet been implemented")
	})

	api.GetDriverInstanceHandler = GetDriverInstanceHandlerFunc(func(params GetDriverInstanceParams) (*genmodel.DriverInstance, error) {
		return nil, errors.NotImplemented("operation getDriverInstance has not yet been implemented")
	})

	api.CreateServiceHandler = CreateServiceHandlerFunc(func(params CreateServiceParams) (*genmodel.Service, error) {
		return nil, errors.NotImplemented("operation createService has not yet been implemented")
	})

	api.CreateServicePlanHandler = CreateServicePlanHandlerFunc(func(params CreateServicePlanParams) error {
		return errors.NotImplemented("operation createServicePlan has not yet been implemented")
	})

	api.UpdateCatalogHandler = UpdateCatalogHandlerFunc(func(params UpdateCatalogParams) error {
		return errors.NotImplemented("operation updateCatalog has not yet been implemented")
	})

}
