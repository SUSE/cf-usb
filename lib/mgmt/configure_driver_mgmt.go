package mgmt

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/genmodel"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/ccapi"
	. "github.com/hpcloud/cf-usb/lib/operations"
	"github.com/pivotal-golang/lager"
	//	"github.com/xeipuuv/gojsonschema"

	uuid "github.com/satori/go.uuid"
)

func ConfigureDriverAPI(api *UsbMgmtAPI, auth authentication.AuthenticationInterface,
	configProvider config.ConfigProvider, ccServiceBroker ccapi.ServiceBrokerInterface,
	logger lager.Logger) {

	log := logger.Session("usb-mgmt-driver")

	api.UpdateDriverHandler = UpdateDriverHandlerFunc(func(params UpdateDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("update-driver")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		exists, err := configProvider.DriverExists(params.DriverID)
		if err != nil {
			return &UpdateDriverInternalServerError{Payload: err.Error()}
		}
		if exists == false {
			log.Debug("update-driver-does-not-exist", lager.Data{"driver-id-does-not-exit": params.DriverID})
			return &UpdateDriverNotFound{}
		}

		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &UpdateDriverInternalServerError{Payload: err.Error()}
		}
		if driver == nil {
			return &UpdateDriverNotFound{}
		}

		driver.DriverName = *params.Driver.Name

		err = configProvider.SetDriver(params.DriverID, *driver)
		if err != nil {
			return &UpdateDriverInternalServerError{Payload: err.Error()}
		}

		return &UpdateDriverOK{Payload: params.Driver}
	})

	api.UploadDriverHandler = UploadDriverHandlerFunc(func(params UploadDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("upload-driver-bits")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &UploadDriverInternalServerError{Payload: err.Error()}
		}
		if driver == nil {
			return &UploadDriverNotFound{}
		}

		driverType := driver.DriverType

		driverPath, err := configProvider.GetDriversPath()
		if err != nil {
			return &UploadDriverInternalServerError{Payload: err.Error()}
		}

		driverPath = filepath.Join(driverPath, driverType)
		if runtime.GOOS == "windows" {
			driverPath = driverPath + ".exe"
		}

		defer params.File.Data.Close()

		f, err := os.OpenFile(driverPath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return &UploadDriverInternalServerError{Payload: err.Error()}
		}

		defer f.Close()

		reader := bufio.NewReader(params.File.Data)

		sha1 := sha1.New()
		_, err = io.Copy(f, io.TeeReader(reader, sha1))
		if err != nil {
			return &UploadDriverInternalServerError{Payload: err.Error()}
		}

		sha := base64.StdEncoding.EncodeToString(sha1.Sum(nil))
		if sha != params.Sha {
			f.Close()
			os.Remove(driverPath)

			return &UploadDriverInternalServerError{Payload: fmt.Sprintf("Checksum mismatch. Actual file SHA1 checksum: %s. Expected file SHA1 checksum: %s", sha, params.Sha)}
		}

		return &UploadDriverOK{}
	})

	api.GetDriverSchemaHandler = GetDriverSchemaHandlerFunc(func(params GetDriverSchemaParams, principal interface{}) middleware.Responder {
		/*log := log.Session("get-driver-schema")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		path, err := configProvider.GetDriversPath()
		if err != nil {
			return &GetDriverSchemaInternalServerError{Payload: err.Error()}
		}
		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDriverSchemaNotFound{}
		}
		schema, err := lib.GetConfigSchema(path, driver.DriverType, logger)
		if err != nil {
			return &GetDriverSchemaInternalServerError{Payload: err.Error()}
		}
		*/
		return &GetDriverSchemaOK{}
	})

	api.GetDriverHandler = GetDriverHandlerFunc(func(params GetDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		d, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDriverInstanceInternalServerError{Payload: err.Error()}
		}
		if d == nil {
			return &GetDriverNotFound{}
		}

		var instances = make([]string, 0)
		for instanceID, _ := range d.DriverInstances {
			instances = append(instances, instanceID)
		}

		driver := &genmodel.Driver{
			DriverType:      &d.DriverType,
			DriverInstances: instances,
			ID:              params.DriverID,
			Name:            &d.DriverName,
		}

		return &GetDriverOK{Payload: driver}
	})

	api.GetDriversHandler = GetDriversHandlerFunc(func(principal interface{}) middleware.Responder {
		log := log.Session("get-drivers")
		log.Info("request")

		var drivers = make([]*genmodel.Driver, 0)

		config, err := configProvider.LoadConfiguration()
		if err != nil {
			return &GetDriversInternalServerError{Payload: err.Error()}
		}

		for dId, d := range config.Drivers {
			var instances = make([]string, 0)
			for instanceID, _ := range d.DriverInstances {
				instances = append(instances, instanceID)
			}

			var driverID string
			driverID = dId
			var name string
			name = d.DriverName
			var dtype string
			dtype = d.DriverType

			driver := &genmodel.Driver{
				ID:              driverID,
				DriverType:      &dtype,
				Name:            &name,
				DriverInstances: instances,
			}

			drivers = append(drivers, driver)
		}

		log.Debug("", lager.Data{"drivers-found": len(drivers)})

		return &GetDriversOK{Payload: drivers}
	})

	api.CreateDriverHandler = CreateDriverHandlerFunc(func(params CreateDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("create-driver")
		log.Info("request", lager.Data{"driver-name": params.Driver.Name, "driver-type": params.Driver.DriverType})

		exist, err := configProvider.DriverTypeExists(*params.Driver.DriverType)
		if err != nil {
			return &CreateDriverInternalServerError{Payload: err.Error()}
		}
		if exist {
			return &CreateDriverConflict{}
		}

		var driver config.Driver

		driver.DriverType = *params.Driver.DriverType
		driver.DriverName = *params.Driver.Name

		driverID := uuid.NewV4().String()

		err = configProvider.SetDriver(driverID, driver)
		if err != nil {
			return &CreateDriverInternalServerError{Payload: err.Error()}
		}

		params.Driver.ID = driverID

		return &CreateDriverCreated{Payload: params.Driver}
	})

	api.GetDriverInstancesHandler = GetDriverInstancesHandlerFunc(func(params GetDriverInstancesParams, principal interface{}) middleware.Responder {
		log := log.Session("get-driver-instances")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		var driverInstances = make([]*genmodel.DriverInstance, 0)

		driver, err := configProvider.GetDriver(params.DriverID)
		if err != nil {
			return &GetDriverInstanceInternalServerError{Payload: err.Error()}
		}
		if driver == nil {
			return &GetDriverNotFound{}
		}

		for diID, di := range driver.DriverInstances {
			var dials = make([]string, 0)
			for dialID, _ := range di.Dials {
				dials = append(dials, dialID)
			}

			var driverInstanceID string
			driverInstanceID = diID

			var serviceID string
			serviceID = di.Service.ID

			driverInstance := &genmodel.DriverInstance{
				Configuration: di.Configuration,
				TargetURL:     di.TargetURL,
				Dials:         dials,
				DriverID:      &params.DriverID,
				ID:            driverInstanceID,
				Name:          &di.Name,
				Service:       serviceID,
			}

			driverInstances = append(driverInstances, driverInstance)
		}

		log.Debug("", lager.Data{"driver-instances-found": len(driverInstances)})

		return &GetDriverInstancesOK{Payload: driverInstances}
	})

	api.PingDriverInstanceHandler = PingDriverInstanceHandlerFunc(func(params PingDriverInstanceParams, principal interface{}) middleware.Responder {
		/*	log := log.Session("ping-driver-instance")
			log.Info("request", lager.Data{"driver-instance-id": params.DriverInstanceID})

			configuration, err := configProvider.LoadConfiguration()
			if err != nil {
				return &PingDriverInstanceInternalServerError{Payload: err.Error()}
			}

			for _, driver := range configuration.Drivers {
				for instanceID, instance := range driver.DriverInstances {
					if instanceID == params.DriverInstanceID {
						result, err := lib.Ping(instance.Configuration, configuration.DriversPath, driver.DriverType)
						if err != nil {
							return &PingDriverInstanceInternalServerError{Payload: err.Error()}
						}
						if result == true {
							return &PingDriverInstanceOK{}
						} else {
							return &PingDriverInstanceNotFound{}
						}
					}
				}
			}
			return &PingDriverInstanceNotFound{}*/
		return &PingDriverInstanceOK{}
	})
	api.DeleteDriverHandler = DeleteDriverHandlerFunc(func(params DeleteDriverParams, principal interface{}) middleware.Responder {
		log := log.Session("delete-driver")
		log.Info("request", lager.Data{"driver-id": params.DriverID})

		driver, err := configProvider.GetDriver(params.DriverID)

		if len(driver.DriverInstances) > 0 {
			return &DeleteDriverInstanceInternalServerError{Payload: fmt.Sprintf("Cannot delete driver '%s' while instances still exist", driver.DriverName)}
		}
		if err != nil {
			return &DeleteDriverInternalServerError{Payload: err.Error()}
		}
		if driver == nil {
			return &DeleteDriverNotFound{}
		}

		driverPath, err := configProvider.GetDriversPath()
		if err != nil {
			return &DeleteDriverInternalServerError{Payload: err.Error()}
		}

		driverPath = filepath.Join(driverPath, driver.DriverType)
		if runtime.GOOS == "windows" {
			driverPath = driverPath + ".exe"
		}

		if _, err := os.Stat(driverPath); err == nil {
			err = os.Remove(driverPath)
			if err != nil {
				return &DeleteDriverInternalServerError{Payload: err.Error()}
			}
		}

		err = configProvider.DeleteDriver(params.DriverID)
		if err != nil {
			return &DeleteDriverInternalServerError{Payload: err.Error()}
		}

		return &DeleteDriverNoContent{}
	})

}
