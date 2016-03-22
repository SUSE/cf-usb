package mgmttest

import (
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

func TestMgmtApiConsulProviderGetInfoCreateDriverAndInstance(t *testing.T) {
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

	executeGetInfoTest(t, ManagementApiPort)

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

func TestMgmtApiConsulProviderUpdateCatalog(t *testing.T) {
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

	SetupCcHttpFakeResponsesUpdateCatalog(uaaFakeServer, ccFakeServer)
	executeUpdateCatalogTest(t, ManagementApiPort)
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
							SetupCcHttpFakeResponsesDeleteDriverInstance(d.DriverType, uaaFakeServer, ccFakeServer)

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

func executeGetInfoTest(t *testing.T, managementApiPort uint16) {
	getInfoResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/info", managementApiPort), nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer getInfoResp.Body.Close()

	Expect(getInfoResp.StatusCode).To((Equal(200)))

	getInfoContent, err := ioutil.ReadAll(getInfoResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get info response content: %s", string(getInfoContent))
	Expect(string(getInfoContent)).To(ContainSubstring("version"))
}

func executeUpdateCatalogTest(t *testing.T, managementApiPort uint16) {
	updateCatalogResp, err := ExecuteHttpCall("POST", fmt.Sprintf("http://localhost:%[1]v/update_catalog", managementApiPort), nil, true)
	if err != nil {
		t.Fatal(err)
	}
	defer updateCatalogResp.Body.Close()

	Expect(updateCatalogResp.StatusCode).To((Equal(200)))

	updateCatalogContent, err := ioutil.ReadAll(updateCatalogResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("update catalog response content: %s", string(updateCatalogContent))

	Expect(string(updateCatalogContent)).To(Equal(""))
}
