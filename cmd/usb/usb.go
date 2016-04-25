package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/frodenas/brokerapi"
	loads "github.com/go-openapi/loads"

	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/mgmt"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/uaaapi"
	"github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"
)

type Usb interface {
	GetCommands() []CLICommandProvider
	Run(config.ConfigProvider, lager.Logger)
}

type UsbApp struct {
	config *config.Config
	logger lager.Logger
}

func NewUsbApp() Usb {
	return &UsbApp{config: &config.Config{}}
}

func (usb *UsbApp) GetCommands() []CLICommandProvider {
	return []CLICommandProvider{
		&FileConfigProvider{},
		&ConsulConfigProvider{},
		&RedisConfigProvider{},
	}
}

func (usb *UsbApp) Run(configProvider config.ConfigProvider, logger lager.Logger) {
	var err error
	usb.logger = logger
	usb.config, err = configProvider.LoadConfiguration()
	if err != nil {
		fmt.Println("Unable to load configuration", err.Error())
		os.Exit(1)
	}

	usb.logger.Info("initializing-drivers")

	usbService := lib.NewUsbBroker(configProvider, usb.logger)

	usb.logger.Info("initializing-brokerapi")

	brokerAPI := brokerapi.New(usbService, usb.logger, usb.config.BrokerAPI.Credentials)

	addr := usb.config.BrokerAPI.Listen

	if usb.config.ManagementAPI != nil {
		go func() {
			logger := usb.logger.Session("management-api")

			logger.Info("starting")

			mgmtaddr := usb.config.ManagementAPI.Listen
			swaggerJSON, err := data.Asset("swagger-spec/api.json")
			if err != nil {
				logger.Fatal("loading-swagger-asset-failed", err)
			}

			swaggerSpec, err := loads.Analyzed(swaggerJSON, "")
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

			client := httpclient.NewHttpClient(usb.config.ManagementAPI.CloudController.SkipTlsValidation)
			info := ccapi.NewGetInfo(usb.config.ManagementAPI.CloudController.Api, client, logger)
			tokenUrl, err := info.GetTokenEndpoint()
			if err != nil {
				logger.Fatal("retrieving-uaa-endpoint-failed", err)
			}

			tokenGenerator := uaaapi.NewTokenGenerator(tokenUrl, usb.config.ManagementAPI.UaaClient, usb.config.ManagementAPI.UaaSecret, client, logger)

			ccServiceBroker := ccapi.NewServiceBroker(client, tokenGenerator, usb.config.ManagementAPI.CloudController.Api, logger)

			mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)
			mgmt.ConfigureAPI(mgmtAPI, auth, configProvider, ccServiceBroker, logger, version)

			logger.Info("start-listening", lager.Data{"address": mgmtaddr})
			err = http.ListenAndServe(mgmtaddr, mgmtAPI.Serve(nil))
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
