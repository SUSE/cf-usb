package instancetest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

func TestMgmtApiConsulProviderCreateDriverInstance(t *testing.T) {
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

			driverInstanceId := executeCreateDriverInstanceTest(t, ManagementApiPort, driverResp, driver.SetDriverInstanceValuesFunc)
			executePingDriverInstanceTest(t, ManagementApiPort, driverInstanceId)
			executeCheckServiceCreatedTest(t, ManagementApiPort, driverInstanceId)
		}
	}
}

func TestMgmtApiConsulProviderUpdateDriverInstance(t *testing.T) {
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

	_, ccFakeServer := SetFakeServers()

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
					driverInstances := executeGetDriverInstancesTest(t, ManagementApiPort, d.Id)

					if len(driverInstances) > 0 {
						firstDriverInstance := driverInstances[0]

						executeGetDriverDialSchemaTest(t, ManagementApiPort, d.Id)

						executeGetDriverConfigSchemaTest(t, ManagementApiPort, d.Id, driver.AssertDriverSchemaContainsFunc)

						executeNegativeUpdateDriverInstanceTest(t, ManagementApiPort, firstDriverInstance, d.Id)

						executeUpdateDriverInstanceTest(t, ManagementApiPort, firstDriverInstance, d.Id, driver.SetDriverInstanceValuesFunc)
					}
				}
			}
		}
	}
}

func TestMgmtApiConsulProviderDeleteDriverInstance(t *testing.T) {
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
					driverInstances := executeGetDriverInstancesTest(t, ManagementApiPort, d.Id)

					if len(driverInstances) > 0 {
						for _, di := range driverInstances {
							SetupCcHttpFakeResponsesDeleteDriverInstance(strings.Replace(di.Name, "updi", "", -1), uaaFakeServer, ccFakeServer)

							executeDeleteDriverInstanceTest(t, ManagementApiPort, di.Id)
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

func executeCreateDriverInstanceTest(t *testing.T, managementApiPort uint16, driver DriverResponse, driverInstanceValues func(driverName, driverId string) []byte) string {
	createDriverStatusCode, driverInstContent, driverInstance, err := CreateDriverInstance(managementApiPort, driver, driverInstanceValues)
	if err != nil {
		t.Fatal(err)
	}

	Expect(createDriverStatusCode).To(Equal(201))

	t.Logf("create driver instance content: %s", driverInstContent)

	Expect(driverInstance.Name).To(ContainSubstring(driver.Name))

	return driverInstance.Id
}

func executePingDriverInstanceTest(t *testing.T, managementApiPort uint16, driverInstanceId string) {
	pingDriverInstResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s/ping", managementApiPort, driverInstanceId), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer pingDriverInstResp.Body.Close()

	Expect(pingDriverInstResp.StatusCode).To(Equal(200))

	driverInstPingContent, err := ioutil.ReadAll(pingDriverInstResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	Expect(string(driverInstPingContent)).To(Equal(""))

	t.Logf("ping driver instance succedeed")
}

func executeCheckServiceCreatedTest(t *testing.T, managementApiPort uint16, driverInstanceId string) {
	getServiceResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/services?driver_instance_id=%[2]s", managementApiPort, driverInstanceId), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getServiceResp.Body.Close()

	Expect(getServiceResp.StatusCode).To(Equal(200))

	getServiceContent, err := ioutil.ReadAll(getServiceResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("get service response content: %s", string(getServiceContent))

	Expect(getServiceContent).To(ContainSubstring(driverInstanceId))
}

func executeGetDriverInstancesTest(t *testing.T, managementApiPort uint16, driverId string) []DriverInstanceResponse {
	getDriverInstancesStatusCode, driverInstancesContent, driverInstances, err := GetDriverInstances(managementApiPort, driverId)
	if err != nil {
		t.Fatal(err)
	}

	Expect(getDriverInstancesStatusCode).To(Equal(200))

	t.Logf("get driver instances response content: %s", string(driverInstancesContent))
	t.Logf("driver instances count: %v", len(driverInstances))

	return driverInstances
}

func executeGetDriverDialSchemaTest(t *testing.T, managementApiPort uint16, driverId string) {
	getDialSchemaResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s/dial_schema", managementApiPort, driverId), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDialSchemaResp.Body.Close()

	Expect(getDialSchemaResp.StatusCode).To(Equal(200))

	getDialSchemaContent, err := ioutil.ReadAll(getDialSchemaResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver dial schema content: %s", string(getDialSchemaContent))
	Expect(getDialSchemaContent).To(ContainSubstring("{}"))
}

func executeGetDriverConfigSchemaTest(t *testing.T, managementApiPort uint16, driverId string, assertDriverSchemaContains func(schemaContent string)) {
	getConfigSchemaResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s/config_schema", managementApiPort, driverId), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getConfigSchemaResp.Body.Close()

	Expect(getConfigSchemaResp.StatusCode).To(Equal(200))

	getConfigSchemaContent, err := ioutil.ReadAll(getConfigSchemaResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver config schema content: %s", string(getConfigSchemaContent))
	assertDriverSchemaContains(string(getConfigSchemaContent))
}

func executeNegativeUpdateDriverInstanceTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse, driverId string) {
	// negative test update driver instance wrong configuration
	newDriverInstanceNameNeg := firstDriverInstance.Name + "updierr"

	instanceValuesNeg := []byte(fmt.Sprintf(`{"name":"%[1]s", "driver_id":"%[2]s", "configuration": {"aKey": "aValue"}}`,
		newDriverInstanceNameNeg,
		driverId))

	updateDriverInstNegResp, err := ExecuteHttpCall("PUT",
		fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, firstDriverInstance.Id),
		bytes.NewBuffer(instanceValuesNeg))
	if err != nil {
		t.Fatal(err)
	}
	defer updateDriverInstNegResp.Body.Close()

	Expect(updateDriverInstNegResp.StatusCode).To(Equal(500))
	updateDriverInstNegContent, err := ioutil.ReadAll(updateDriverInstNegResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("negative test update driver instance response content: %s", string(updateDriverInstNegContent))
	Expect(string(updateDriverInstNegContent)).To(ContainSubstring("is required"))
}

func executeUpdateDriverInstanceTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse, driverId string, driverInstanceValues func(driverName, driverId string) []byte) {
	newDriverInstanceName := firstDriverInstance.Name + "updi"

	instanceValues := driverInstanceValues(newDriverInstanceName, driverId)

	updateDriverInstResp, err := ExecuteHttpCall("PUT",
		fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, firstDriverInstance.Id),
		bytes.NewBuffer(instanceValues))

	if err != nil {
		t.Fatal(err)
	}
	defer updateDriverInstResp.Body.Close()

	updateDriverInstContent, err := ioutil.ReadAll(updateDriverInstResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("update driver instance response content: %s", string(updateDriverInstContent))

	if updateDriverInstResp.StatusCode == 200 {
		Expect(updateDriverInstContent).To(ContainSubstring(newDriverInstanceName))
	}
}

func executeDeleteDriverInstanceTest(t *testing.T, managementApiPort uint16, driverInstanceId string) {
	deleteDriverInstStatusCode, deleteDriverInstContent, err := DeleteDriverInstance(managementApiPort, driverInstanceId)
	if err != nil {
		t.Fatal(err)
	}

	Expect(deleteDriverInstStatusCode).To(Equal(204))

	Expect(string(deleteDriverInstContent)).To(Equal(""))

	getDriverDeletedInstancesResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, driverInstanceId), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDriverDeletedInstancesResp.Body.Close()

	Expect(getDriverDeletedInstancesResp.StatusCode).To((Equal(404)))
}
