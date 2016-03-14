package drivertest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func TestMgmtApiConsulProviderCreateDriver(t *testing.T) {
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
			driverId := executeCreateDriverTest(t, ManagementApiPort, driver.DriverType)
			executeGetDriverTest(t, ManagementApiPort, driverId)
			executeUploadDriverTest(t, ManagementApiPort, driver.DriverType, driverId)
		}
	}
}

func TestMgmtApiConsulProviderUpdateDriver(t *testing.T) {
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

	driversContent, driversResp := executeGetDriversTest(t, ManagementApiPort)

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {

			Expect(driversContent).To(ContainSubstring(driver.DriverType))

			for _, d := range driversResp {
				if d.DriverType == driver.DriverType {
					executeUpdateDriverTest(t, ManagementApiPort, d, driver.AssertDriverSchemaContainsFunc)
				}
			}
		}
	}
}

func TestMgmtApiConsulProviderDeleteDriver(t *testing.T) {
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

	driversContent, driversResp := executeGetDriversTest(t, ManagementApiPort)

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {

			Expect(driversContent).To(ContainSubstring(driver.DriverType))

			for _, d := range driversResp {
				if d.DriverType == driver.DriverType {
					executeDeleteDriverTest(t, ManagementApiPort, d)
				}
			}
		}
	}
}

func executeCreateDriverTest(t *testing.T, managementApiPort uint16, driverName string) string {
	createDriverStatusCode, driverContent, driver, err := CreateDriver(managementApiPort, driverName)
	if err != nil {
		t.Fatal(err)
	}

	if createDriverStatusCode == 409 {
		t.Skipf("Skipping test as driver type %[1]s already exists", driverName)
	}

	Expect(createDriverStatusCode).To(Equal(201))

	t.Logf("create driver response content: %s", driverContent)

	Expect(driver.Id).ToNot(BeNil())
	Expect(driver.Name).To(Equal(driverName))

	return driver.Id
}

func executeGetDriverTest(t *testing.T, managementApiPort uint16, driverId string) {
	getDriverResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driverId), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDriverResp.Body.Close()

	Expect(getDriverResp.StatusCode).To(Equal(200))

	getDriverContent, err := ioutil.ReadAll(getDriverResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver response content: %s", string(getDriverContent))
	Expect(getDriverContent).To(ContainSubstring(driverId))
}

func executeUploadDriverTest(t *testing.T, managementApiPort uint16, driverType, driverId string) {
	statusCode, uploadDriverBitsContent, err := UploadDriver(managementApiPort, driverType, driverId)
	if err != nil {
		t.Fatal(err)
	}

	Expect(statusCode).To(Equal(200))
	Expect(string(uploadDriverBitsContent)).To(Equal(""))

	//check if driver bits file exists
	_, err = os.Stat(path.Join(os.Getenv("USB_DRIVER_PATH"), driverType))
	Expect(err).To(BeNil())
}

func executeGetDriversTest(t *testing.T, managementApiPort uint16) (string, []DriverResponse) {
	getDriversStatusCode, driversContent, drivers, err := GetDrivers(managementApiPort)
	if err != nil {
		t.Fatal(err)
	}

	Expect(getDriversStatusCode).To(Equal(200))

	t.Logf("get drivers response content: %s", driversContent)

	return driversContent, drivers
}

func executeUpdateDriverTest(t *testing.T, managementApiPort uint16, driver DriverResponse, assertDriverSchemaContains func(schemaContent string)) {
	updateDriverName := driver.Name + "updd"

	driverValues := []byte(fmt.Sprintf(`{"driver_type":"%[1]s","name":"%[2]s"}`,
		driver.DriverType,
		updateDriverName))

	updateDriverResp, err := ExecuteHttpCall("PUT", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driver.Id), bytes.NewBuffer(driverValues))
	if err != nil {
		t.Fatal(err)
	}
	defer updateDriverResp.Body.Close()

	Expect(updateDriverResp.StatusCode).To(Equal(200))

	updatedDriverContent, err := ioutil.ReadAll(updateDriverResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver response content: %s", string(updatedDriverContent))

	Expect(updatedDriverContent).To(ContainSubstring(updateDriverName))
}

func executeDeleteDriverTest(t *testing.T, managementApiPort uint16, driver DriverResponse) {
	deleteDriverStatusCode, deleteDriverContent, err := DeleteDriver(managementApiPort, driver.Id)
	if err != nil {
		t.Fatal(err)
	}

	Expect(deleteDriverStatusCode).To(Equal(204))

	Expect(string(deleteDriverContent)).To(Equal(""))

	getDriverDeletedResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driver.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDriverDeletedResp.Body.Close()

	Expect(getDriverDeletedResp.StatusCode).To(Equal(404))
}
