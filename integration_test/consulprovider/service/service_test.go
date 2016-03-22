package servicetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hpcloud/cf-usb/lib/config"

	. "github.com/hpcloud/cf-usb/integration_test/consulprovider"
	. "github.com/onsi/gomega"
)

func init() {
	ConsulConfig.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.ConsulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.ConsulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.ConsulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.ConsulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.ConsulToken = os.Getenv("CONSUL_TOKEN")
}

func TestMgmtApiConsulProviderCreateDriverAndInstance(t *testing.T) {
	RegisterTestingT(t)

	binpath, buildNotExist, err := CheckSolutionIsBuild()
	if buildNotExist {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}

	uaaFakeServer, ccFakeServer := SetFakeServers()

	err = SetupConsulForFirstRun()
	if err != nil {
		t.Fatal(err)
	}

	consulClient, err := RunConsulProcess(BrokerApiPort, ManagementApiPort, ccFakeServer.URL(), TempDriversPath)
	if err != nil {
		t.Fatal(err)
	}

	if ConsulProcess != nil {
		defer func() {
			ConsulProcess.Signal(os.Kill)
			<-ConsulProcess.Wait()
		}()
	}

	t.Log("consul started")

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	usbProcess, err := RunUsbProcess(binpath, ConsulConfig.ConsulAddress)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		usbProcess.Signal(os.Kill)
		<-usbProcess.Wait()
	}()

	t.Log("usb started")

	//wait for process to start
	time.Sleep(5 * time.Second)

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			_, _, driverResp, err := CreateDriver(ManagementApiPort, driver.DriverType)
			if err != nil {
				t.Fatal(err)
			}

			_, _, err = UploadDriver(ManagementApiPort, driverResp.DriverType, driverResp.Id)
			if err != nil {
				t.Fatal(err)
			}

			SetupCcHttpFakeResponsesCreateDriverInstance(uaaFakeServer, ccFakeServer)

			_, _, _, err = CreateDriverInstance(ManagementApiPort, driverResp, driver.SetDriverInstanceValuesFunc)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestMgmtApiConsulProviderUpdateService(t *testing.T) {
	RegisterTestingT(t)

	binpath, buildNotExist, err := CheckSolutionIsBuild()
	if buildNotExist {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}

	uaaFakeServer, ccFakeServer := SetFakeServers()

	err = SetupConsulForFirstRun()
	if err != nil {
		t.Fatal(err)
	}

	consulClient, err := RunConsulProcess(BrokerApiPort, ManagementApiPort, ccFakeServer.URL(), TempDriversPath)
	if err != nil {
		t.Fatal(err)
	}

	if ConsulProcess != nil {
		defer func() {
			ConsulProcess.Signal(os.Kill)
			<-ConsulProcess.Wait()
		}()
	}

	t.Log("consul started")

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	usbProcess, err := RunUsbProcess(binpath, ConsulConfig.ConsulAddress)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		usbProcess.Signal(os.Kill)
		<-usbProcess.Wait()
	}()

	t.Log("usb started")

	//wait for process to start
	time.Sleep(5 * time.Second)

	_, _, driversResp, err := GetDrivers(ManagementApiPort)
	if err != nil {
		t.Fatal(err)
	}

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			for _, d := range driversResp {
				if d.DriverType == driver.DriverType {
					_, _, driverInstancesResponse, err := GetDriverInstances(ManagementApiPort, d.Id)
					if err != nil {
						t.Fatal(err)
					}

					if len(driverInstancesResponse) > 0 {
						firstDriverInstance := driverInstancesResponse[0]

						SetupCcHttpFakeResponsesUpdateService(fmt.Sprintf("%supds", d.DriverType), uaaFakeServer, ccFakeServer)

						serviceByInstanceId := executeGetServiceByInstanceIdTest(t, ManagementApiPort, firstDriverInstance)
						service := executeGetServiceTest(t, ManagementApiPort, firstDriverInstance)

						Expect(serviceByInstanceId.Description).To(Equal(service.Description))
						Expect(serviceByInstanceId.Bindable).To(Equal(service.Bindable))
						Expect(serviceByInstanceId.DriverInstanceID).To(Equal(service.DriverInstanceID))
						Expect(serviceByInstanceId.Id).To(Equal(service.Id))
						Expect(serviceByInstanceId.Name).To(Equal(service.Name))
						Expect(serviceByInstanceId.Tags).To(Equal(service.Tags))

						executeUpdateServiceTest(t, ManagementApiPort, service)
					}
				}
			}
		}
	}
}

func TestMgmtApiConsulProviderDeleteDriverAndInstance(t *testing.T) {
	RegisterTestingT(t)

	binpath, buildNotExist, err := CheckSolutionIsBuild()
	if buildNotExist {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	if ConsulConfig.ConsulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}

	uaaFakeServer, ccFakeServer := SetFakeServers()

	err = SetupConsulForFirstRun()
	if err != nil {
		t.Fatal(err)
	}

	consulClient, err := RunConsulProcess(BrokerApiPort, ManagementApiPort, ccFakeServer.URL(), TempDriversPath)
	if err != nil {
		t.Fatal(err)
	}

	if ConsulProcess != nil {
		defer func() {
			ConsulProcess.Signal(os.Kill)
			<-ConsulProcess.Wait()
		}()
	}

	t.Log("consul started")

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	usbProcess, err := RunUsbProcess(binpath, ConsulConfig.ConsulAddress)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		usbProcess.Signal(os.Kill)
		<-usbProcess.Wait()
	}()

	t.Log("usb started")

	//wait for process to start
	time.Sleep(5 * time.Second)

	_, _, driversResp, err := GetDrivers(ManagementApiPort)
	if err != nil {
		t.Fatal(err)
	}

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			for _, d := range driversResp {
				if d.DriverType == driver.DriverType {
					_, _, driverInstancesResponse, err := GetDriverInstances(ManagementApiPort, d.Id)
					if err != nil {
						t.Fatal(err)
					}

					if len(driverInstancesResponse) > 0 {
						for _, di := range driverInstancesResponse {
							SetupCcHttpFakeResponsesDeleteDriverInstance(fmt.Sprintf("%supds", d.DriverType), uaaFakeServer, ccFakeServer)

							_, _, err = DeleteDriverInstance(ManagementApiPort, di.Id)
							if err != nil {
								t.Fatal(err)
							}
						}
					}
					_, _, err = DeleteDriver(ManagementApiPort, d.Id)
					if err != nil {
						t.Fatal(err)
					}
				}
			}
		}
	}
}

func executeGetServiceByInstanceIdTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse) ServiceResponse {
	getServiceResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/services?driver_instance_id=%[2]s", managementApiPort, firstDriverInstance.Id), nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer getServiceResp.Body.Close()

	Expect(200).To((Equal(getServiceResp.StatusCode)))

	getServiceContent, err := ioutil.ReadAll(getServiceResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get service by driver instance id response content: %s", string(getServiceContent))

	var service ServiceResponse

	err = json.Unmarshal(getServiceContent, &service)
	if err != nil {
		t.Fatal(err)
	}
	Expect(service.Description).To(ContainSubstring("Default"))
	Expect(service.Bindable).To(Equal(true))

	return service
}

func executeGetServiceTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse) ServiceResponse {
	getServiceResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/services/%[2]s", managementApiPort, firstDriverInstance.Service), nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer getServiceResp.Body.Close()

	Expect(200).To((Equal(getServiceResp.StatusCode)))

	getServiceContent, err := ioutil.ReadAll(getServiceResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get service response content: %s", string(getServiceContent))

	var service ServiceResponse

	err = json.Unmarshal(getServiceContent, &service)
	if err != nil {
		t.Fatal(err)
	}
	Expect(service.Description).To(ContainSubstring("Default"))
	Expect(service.Bindable).To(Equal(true))

	return service
}

func executeUpdateServiceTest(t *testing.T, managementApiPort uint16, service ServiceResponse) {
	updateServiceName := service.Name + "upds"
	updateServiceDesc := service.Description + "upds"

	serviceValues := []byte(fmt.Sprintf(`{"bindable":%[1]v,"description":"%[2]s","driver_instance_id":"%[3]s","name":"%[4]s"}`,
		service.Bindable,
		updateServiceDesc,
		service.DriverInstanceID,
		updateServiceName))

	updateServiceResp, err := ExecuteHttpCall("PUT", fmt.Sprintf("http://localhost:%[1]v/services/%[2]s", managementApiPort, service.Id), bytes.NewBuffer(serviceValues), true)
	if err != nil {
		t.Fatal(err)
	}
	defer updateServiceResp.Body.Close()

	Expect(200).To(Equal(updateServiceResp.StatusCode))

	updateServiceContent, err := ioutil.ReadAll(updateServiceResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("updated service content: %s", string(updateServiceContent))

	Expect(updateServiceContent).To(ContainSubstring(updateServiceName))
	Expect(updateServiceContent).To(ContainSubstring(updateServiceDesc))
}
