package mgmt

import (
	"fmt"

	errors "github.com/go-openapi/errors"
	swaggerruntime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"
)

func ConfigureAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface,
	configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface,
	logger lager.Logger, usbVersion string) {
	log := logger.Session("usb-mgmt")

	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = swaggerruntime.JSONConsumer()

	api.JSONProducer = swaggerruntime.JSONProducer()

	api.AuthorizationAuth = func(token string) (interface{}, error) {
		err := auth.IsAuthenticated(token)

		if err != nil {
			return nil, err
		}

		log.Debug("authentication-succeeded")

		return true, nil
	}

	ConfigureDriverAPI(api, auth, configProvider, ccServiceBroker, logger)

	ConfigureInstanceAPI(api, auth, configProvider, ccServiceBroker, logger)

	ConfigureDialAPI(api, auth, configProvider, ccServiceBroker, logger)

	ConfigureServicePlanAPI(api, auth, configProvider, ccServiceBroker, logger)

	ConfigureServiceAPI(api, auth, configProvider, ccServiceBroker, logger)

	api.GetInfoHandler = GetInfoHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("get-info")
		log.Info("request")

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetInfoInternalServerError{
				Payload: err.Error(),
			}
		}

		info := &genmodel.Info{
			BrokerAPIVersion: &config.APIVersion,
			UsbVersion:       &usbVersion,
		}

		return &GetInfoOK{
			Payload: info,
		}
	})

	api.UpdateCatalogHandler = UpdateCatalogHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("update-catalog")
		log.Info("request")

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &UpdateCatalogInternalServerError{Payload: err.Error()}
		}

		brokerName := defaultBrokerName
		if len(config.ManagementAPI.BrokerName) > 0 {
			brokerName = config.ManagementAPI.BrokerName
		}

		guid, err := ccServiceBroker.GetServiceBrokerGuidByName(brokerName)
		if err != nil {
			log.Error("get-service-broker-failed", err)
			return &UpdateCatalogInternalServerError{Payload: err.Error()}
		}

		if guid == "" {
			return &UpdateCatalogInternalServerError{Payload: fmt.Sprintf("Broker %s guid not found", brokerName)}
		} else {
			err = ccServiceBroker.Update(guid, brokerName, config.BrokerAPI.ExternalUrl, config.BrokerAPI.Credentials.Username, config.BrokerAPI.Credentials.Password)
			if err != nil {
				return &UpdateCatalogInternalServerError{Payload: err.Error()}
			}
		}

		for _, driver := range config.Drivers {
			for _, instance := range driver.DriverInstances {
				err = ccServiceBroker.EnableServiceAccess(instance.Service.Name)
				if err != nil {
					log.Error("enable-service-access-failed", err)
					return &UpdateCatalogInternalServerError{Payload: err.Error()}
				}
			}
		}

		return &UpdateCatalogOK{}
	})
}
