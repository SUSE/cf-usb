package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/mgmt"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	"github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

type Usb interface {
	GetCommands() []CLICommandProvider
	Run(config.ConfigProvider)
}

var logger = lager.NewLogger("usb")

const (
	DEBUG = "debug"
	INFO  = "info"
	ERROR = "error"
	FATAL = "fatal"
)

func getLogLevel(config *config.Config) lager.LogLevel {
	var minLogLevel lager.LogLevel
	switch config.LogLevel {
	case DEBUG:
		minLogLevel = lager.DEBUG
	case INFO:
		minLogLevel = lager.INFO
	case ERROR:
		minLogLevel = lager.ERROR
	case FATAL:
		minLogLevel = lager.FATAL
	default:
		panic(fmt.Errorf("invalid log level: %s", config.LogLevel))
	}

	return minLogLevel
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
	}
}

func (usb *UsbApp) Run(configProvider config.ConfigProvider) {
	var err error
	usb.config, err = configProvider.LoadConfiguration()
	if err != nil {
		fmt.Println("Unable to load configuration", err.Error())
		os.Exit(1)
	}

	logger.RegisterSink(lager.NewWriterSink(os.Stdout, getLogLevel(usb.config)))

	usb.logger = logger

	drivers := usb.getDrivers(configProvider)

	logger.Info("run", lager.Data{"action": "starting drivers"})
	usbService := lib.NewUsbBroker(drivers, usb.config, logger)
	brokerAPI := brokerapi.New(usbService, logger, usb.config.BrokerAPI.Credentials)

	addr := usb.config.BrokerAPI.Listen

	if usb.config.ManagementAPI != nil {
		go func() {
			mgmtaddr := usb.config.ManagementAPI.Listen
			swaggerJSON, err := data.Asset("swagger-spec/api.json")
			if err != nil {
				logger.Fatal("error-start-mgmt-api", err)
			}

			swaggerSpec, err := spec.New(swaggerJSON, "")
			if err != nil {
				logger.Fatal("error-start-mgmt-api", err)
			}

			uaaAuthConfig, err := configProvider.GetUaaAuthConfig()
			if err != nil {
				logger.Error("error-start-mgmt-api", err)
			}

			auth, err := uaa.NewUaaAuth(uaaAuthConfig.PublicKey, uaaAuthConfig.Scope, usb.config.BrokerAPI.DevMode)
			if err != nil {
				logger.Error("error-start-mgmt-api", err)
			}

			mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)
			mgmt.ConfigureAPI(mgmtAPI, auth, usb.config)

			logger.Info("run", lager.Data{"mgmtadd": mgmtaddr})
			http.ListenAndServe(mgmtaddr, mgmtAPI.Serve())
		}()
	}
	logger.Info("run", lager.Data{"addr": addr})
	err = http.ListenAndServe(addr, brokerAPI)
	if err != nil {
		logger.Fatal("error-listening", err)
	}

}

func (usb *UsbApp) getDrivers(configProvider config.ConfigProvider) []*lib.DriverProvider {
	var drivers []*lib.DriverProvider
	for _, driver := range usb.config.Drivers {
		for _, driverInstance := range driver.DriverInstances {
			usb.logger.Info("start-driver ", lager.Data{"driver-type": driver.DriverType,
				"driver-instance": driverInstance.Name})
			instanceConfig, err := configProvider.GetDriverInstanceConfig(driverInstance.ID)
			if err != nil {
				logger.Error("failed-to-load-driver-config", err, lager.Data{"DriverInstance": driverInstance.ID})
			}
			driver := lib.NewDriverProvider(driver.DriverType, instanceConfig, usb.logger)
			err = driver.Validate()
			if err != nil {
				logger.Error("failed-to-validate-driver", err, lager.Data{"DriverInstance": driverInstance.ID})
			}
			drivers = append(drivers, driver)
		}
	}

	return drivers

}
