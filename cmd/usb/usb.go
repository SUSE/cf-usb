package main

import (
	"log"

	"github.com/hpcloud/cf-usb/lib"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/gocfbroker"
	"github.com/hpcloud/gocfbroker/boltdb"
)

type Usb interface {
	GetCommands() []CLICommandProvider
	Run(config.ConfigProvider)
}
type UsbApp struct {
	config config.Config
}

func NewUsbApp() Usb {
	return &UsbApp{config: config.Config{}}
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
		log.Fatalf("Unable to load configuration %s", err.Error())
	}

	//TODO: Use etcd if configured
	db, err := boltdb.New(usb.config.BoltFilename, usb.config.BoltBucket)
	if err != nil {
		log.Fatalln("failed to open database:", err)
	}

	drivers, err := usb.startDrivers(configProvider)
	if err != nil {
		log.Fatalln("Failed to start drivers", err)
	}
	log.Println("drivers started")
	usbService := lib.NewUsbBroker(drivers)
	broker, err := gocfbroker.New(usbService, db, usb.config.Options)
	if err != nil {
		log.Fatalln(err)
	}

	broker.Start()

}

func (usb *UsbApp) startDrivers(configProvider config.ConfigProvider) ([]lib.DriverProvider, error) {
	var drivers []lib.DriverProvider
	log.Println("Detecting drivers")
	driverTypes, err := configProvider.GetDriverTypes()
	if err != nil {
		return drivers, err
	}

	for _, driverType := range driverTypes {
		log.Println("Starting Driver: ", driverType)

		driverProp, err := configProvider.GetDriverProperties(driverType)
		if err != nil {
			return drivers, err
		}

		driver, err := lib.NewDriverProvider(driverType, driverProp)
		if err != nil {
			return drivers, err
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil

}
