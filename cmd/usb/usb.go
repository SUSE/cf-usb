package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/frodenas/brokerapi"
	loads "github.com/go-openapi/loads"

	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/csm"
	"github.com/hpcloud/cf-usb/lib/mgmt"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/hpcloud/cf-usb/lib/mgmt/operations"
	"github.com/pivotal-golang/lager"
)

//Usb is the base type for Universal Service Broker
type Usb interface {
	GetCommands() []CLICommandProvider
	Run(config.Provider, lager.Logger)
}

//UsbApp is the base type to be used by applications
type UsbApp struct {
	config *config.Config
	logger lager.Logger
}

//NewUsbApp creates an instance of UsbApp and returns it's pointer address
func NewUsbApp() Usb {
	return &UsbApp{config: &config.Config{}}
}

//GetCommands returns the available config provider commands
func (usb *UsbApp) GetCommands() []CLICommandProvider {
	return []CLICommandProvider{
		&FileConfigProvider{},
		&ConsulConfigProvider{},
		&RedisConfigProvider{},
	}
}

//Run starts and runs an UsbApp based on the configProvider passed as a param
func (usb *UsbApp) Run(configProvider config.Provider, logger lager.Logger) {
	var err error
	usb.logger = logger
	usb.config, err = configProvider.LoadConfiguration()
	if err != nil {
		fmt.Println("Unable to load configuration", err.Error())
		os.Exit(1)
	}

	usb.logger.Info("initializing-drivers")

	csmClient := csm.NewCSMClient(usb.logger)

	usbService := lib.NewUsbBroker(configProvider, usb.logger, csmClient)

	usb.logger.Info("initializing-brokerapi")

	brokerAPI := brokerapi.New(usbService, usb.logger, usb.config.BrokerAPI.Credentials)

	addr := usb.config.BrokerAPI.Listen

	if usb.config.ManagementAPI != nil {
		go func() {
			logger := usb.logger.Session("management-api")

			logger.Info("starting")

			mgmtaddr := usb.config.ManagementAPI.Listen

			swaggerSpec, err := loads.Analyzed(mgmt.SwaggerJSON, "")
			if err != nil {
				logger.Fatal("initializing-swagger-failed", err)
			}

			uaaAuthConfig, err := configProvider.GetUaaAuthConfig()
			if err != nil {
				logger.Error("initializing-uaa-config-failed", err)
			}

			auth, err := uaa.NewUaaAuth(
				uaaAuthConfig.PublicKey,
				uaaAuthConfig.SymmetricVerificationKey,
				uaaAuthConfig.Scope,
				usb.config.ManagementAPI.DevMode,
				logger)
			if err != nil {
				logger.Fatal("initializing-uaa-auth-failed", err)
			}

			client := httpclient.NewHTTPClient(usb.config.ManagementAPI.CloudController.SkipTLSValidation)
			info := ccapi.NewGetInfo(usb.config.ManagementAPI.CloudController.API, client, logger)
			tokenURL, err := info.GetTokenEndpoint()
			if err != nil {
				logger.Fatal("retrieving-uaa-endpoint-failed", err)
			}

			tokenGenerator := uaaapi.NewTokenGenerator(tokenURL, usb.config.ManagementAPI.UaaClient, usb.config.ManagementAPI.UaaSecret, client, logger)

			ccServiceBroker := ccapi.NewServiceBroker(client, tokenGenerator, usb.config.ManagementAPI.CloudController.API, logger)

			mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)
			api := mgmt.ConfigureAPI(mgmtAPI, auth, configProvider, ccServiceBroker, logger, version)

			logger.Info("start-listening", lager.Data{"address": mgmtaddr})
			err = http.ListenAndServe(mgmtaddr, api)
			if err != nil {
				logger.Fatal("listening-failed", err)
			}
		}()
	}

	if usb.config.RoutesRegister != nil {
		go usb.StartRouteRegistration(usb.config, usb.logger)
	}

	usb.logger.Info("start-listening-brokerapi", lager.Data{"address": addr})
	err = http.ListenAndServe(addr, brokerAPI)
	if err != nil {
		usb.logger.Fatal("listening-brokerapi-failed", err)
	}
}
