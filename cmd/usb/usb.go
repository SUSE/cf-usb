package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
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

	drivers, err := usb.startDrivers(configProvider)
	if err != nil {
		logger.Error("start-drivers", err)
		os.Exit(1)
	}
	logger.Info("run", lager.Data{"action": "starting drivers"})
	usbService := lib.NewUsbBroker(drivers, usb.config, logger)
	brokerAPI := brokerapi.New(usbService, logger, usb.config.Crednetials)

	http.Handle("/", brokerAPI)

	addr := usb.config.Listen

	logger.Info("run", lager.Data{"addr": addr})
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Fatal("error-listening", err)
	}

}

func (usb *UsbApp) startDrivers(configProvider config.ConfigProvider) ([]*lib.DriverProvider, error) {
	var drivers []*lib.DriverProvider
	driverTypes, err := configProvider.GetDriverTypes()
	if err != nil {
		return drivers, err
	}

	for _, driverType := range driverTypes {
		usb.logger.Info("start-driver ", lager.Data{"driver-type": driverType})

		driverProp, err := configProvider.GetDriverProperties(driverType)
		if err != nil {
			return drivers, err
		}

		usb.logger.Info("start-driver", lager.Data{"driver-type": driverType})
		driver, err := lib.NewDriverProvider(driverType, driverProp)

		if err != nil {
			return drivers, err
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil

}
