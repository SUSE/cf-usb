package brokertest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/frodenas/brokerapi"
	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/satori/go.uuid"

	. "github.com/hpcloud/cf-usb/integration_test/consulprovider"
	. "github.com/onsi/gomega"
)

var orgGuid string = uuid.NewV4().String()
var spaceGuid string = uuid.NewV4().String()
var serviceGuid string = uuid.NewV4().String()
var serviceGuidAsync string = fmt.Sprintf("%[1]s-async", uuid.NewV4().String())
var serviceBindingGuid string = uuid.NewV4().String()
var driverInstances []DriverInstanceResponse

func init() {
	ConsulConfig.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.ConsulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.ConsulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.ConsulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.ConsulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.ConsulToken = os.Getenv("CONSUL_TOKEN")
}

func TestBrokerApiConsulProviderCreateDriverAndInstance(t *testing.T) {
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

func TestBrokerApiConsulProviderCatalog(t *testing.T) {
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

	provider := config.NewConsulConfig(consulClient)

	configInfo, err := provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	services := executeCatalogTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials)

	driversCounter := 0
	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			driversCounter++

			for _, service := range services {
				if service.Name == driver.DriverType {

					Expect(service.ID).NotTo(BeNil())
					Expect(service.Name).To(Equal(driver.DriverType))
					Expect(service.Description).To(ContainSubstring("Default"))
					Expect(service.Bindable).To(Equal(true))
					Expect(service.Tags).NotTo(BeNil())
					Expect(service.PlanUpdateable).To(Equal(false))
					Expect(len(service.Plans)).To(Equal(1))
					Expect(service.Plans[0].ID).NotTo(BeNil())
					Expect(service.Plans[0].Name).To(Equal("default"))
					Expect(service.Plans[0].Description).To(Equal("default plan"))
					Expect(service.Plans[0].Metadata).NotTo(BeNil())
					Expect(service.Plans[0].Free).To(Equal(true))

				}
			}
		}
	}

	Expect(len(services)).To(Equal(driversCounter))
}

func TestBrokerApiConsulProviderProvision(t *testing.T) {
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

	provider := config.NewConsulConfig(consulClient)

	configInfo, err := provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	services := executeCatalogTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials)

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			for _, service := range services {
				if service.Name == driver.DriverType {

					if strings.Contains(driver.DriverType, "-async") {
						executeProvisionAsyncTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials, service)
					} else {
						executeProvisionTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials, service)
					}

				}
			}
		}
	}
}

func TestBrokerApiConsulProviderBind(t *testing.T) {
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

	provider := config.NewConsulConfig(consulClient)

	configInfo, err := provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	services := executeCatalogTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials)

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			for _, service := range services {
				if service.Name == driver.DriverType {
					executeBindTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials, service)
				}
			}
		}
	}
}

func TestBrokerApiConsulProviderUnbind(t *testing.T) {
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

	provider := config.NewConsulConfig(consulClient)

	configInfo, err := provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	services := executeCatalogTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials)

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			for _, service := range services {
				if service.Name == driver.DriverType {

					executeUnbindTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials, service, driver.DriverType)

				}
			}
		}
	}
}

func TestBrokerApiConsulProviderDeprovision(t *testing.T) {
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

	provider := config.NewConsulConfig(consulClient)

	configInfo, err := provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	services := executeCatalogTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials)

	for _, driver := range Drivers {
		if driver.EnvVarsExistFunc() {
			for _, service := range services {
				if service.Name == driver.DriverType {

					executeDeprovisionTest(t, BrokerApiPort, configInfo.BrokerAPI.Credentials, service, driver.DriverType)
				}
			}
		}
	}
}

func TestBrokerApiConsulProviderDeleteDriverAndInstance(t *testing.T) {
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

func executeCatalogTest(t *testing.T, brokerApiPort uint16, credentials brokerapi.BrokerCredentials) []brokerapi.Service {
	catalogResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://%[1]s:%[2]s@localhost:%[3]v/v2/catalog", credentials.Username, credentials.Password, brokerApiPort), nil, false)
	if err != nil {
		t.Fatal(err)
	}
	defer catalogResp.Body.Close()

	Expect(catalogResp.StatusCode).To(Equal(200))

	catalogContent, err := ioutil.ReadAll(catalogResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("catalog response content: %s", string(catalogContent))

	var catalog brokerapi.CatalogResponse

	err = json.Unmarshal(catalogContent, &catalog)
	if err != nil {
		t.Fatal(err)
	}

	return catalog.Services
}

func executeProvisionAsyncTest(t *testing.T, brokerApiPort uint16, credentials brokerapi.BrokerCredentials, service brokerapi.Service) {
	serviceValues := []byte(fmt.Sprintf(`{"service_id":"%[1]s", "plan_id":"%[2]s", "organization_guid": "%[3]s", "space_guid": "%[4]s"}`,
		service.ID,
		service.Plans[0].ID,
		orgGuid,
		spaceGuid))

	t.Logf("start provisioning service %[1]s, with service guid %[2]s", service.Name, serviceGuid)

	provisionResp, err := ExecuteHttpCall(
		"PUT",
		fmt.Sprintf("http://%s:%s@localhost:%[3]v/v2/service_instances/%[4]s?accepts_incomplete=true", credentials.Username, credentials.Password, brokerApiPort, serviceGuid),
		bytes.NewBuffer(serviceValues),
		false)
	if err != nil {
		t.Fatal(err)
	}
	defer provisionResp.Body.Close()

	Expect(provisionResp.StatusCode).To(Equal(202))

	provisionContent, err := ioutil.ReadAll(provisionResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("provision response content: %s", string(provisionContent))
}

func executeProvisionTest(t *testing.T, brokerApiPort uint16, credentials brokerapi.BrokerCredentials, service brokerapi.Service) {
	serviceValues := []byte(fmt.Sprintf(`{"service_id":"%[1]s", "plan_id":"%[2]s", "organization_guid": "%[3]s", "space_guid": "%[4]s"}`,
		service.ID,
		service.Plans[0].ID,
		orgGuid,
		spaceGuid))

	t.Logf("start provisioning service %[1]s, with service guid %[2]s", service.Name, serviceGuid)

	provisionResp, err := ExecuteHttpCall(
		"PUT",
		fmt.Sprintf("http://%s:%s@localhost:%[3]v/v2/service_instances/%[4]s", credentials.Username, credentials.Password, brokerApiPort, serviceGuid),
		bytes.NewBuffer(serviceValues),
		false)
	if err != nil {
		t.Fatal(err)
	}
	defer provisionResp.Body.Close()

	Expect(provisionResp.StatusCode).To(Equal(201))

	provisionContent, err := ioutil.ReadAll(provisionResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("provision response content: %s", string(provisionContent))
}

func executeBindTest(t *testing.T, brokerApiPort uint16, credentials brokerapi.BrokerCredentials, service brokerapi.Service) {
	bindValues := []byte(fmt.Sprintf(`{"service_id":"%[1]s", "plan_id":"%[2]s"}`,
		service.ID,
		service.Plans[0].ID))

	t.Logf("start binding service %[1]s, with service binding guid %[2]s", service.Name, serviceBindingGuid)

	bindResp, err := ExecuteHttpCall(
		"PUT",
		fmt.Sprintf("http://%s:%s@localhost:%[3]v/v2/service_instances/%[4]s/service_bindings/%[5]s", credentials.Username, credentials.Password, brokerApiPort, serviceGuid, serviceBindingGuid),
		bytes.NewBuffer(bindValues),
		false)
	if err != nil {
		t.Fatal(err)
	}
	defer bindResp.Body.Close()

	Expect(bindResp.StatusCode).To(Equal(201))

	bindContent, err := ioutil.ReadAll(bindResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("bind response content: %s", string(bindContent))
}

func executeUnbindTest(t *testing.T, brokerApiPort uint16, credentials brokerapi.BrokerCredentials, service brokerapi.Service, driverType string) {
	servicePlanId := service.Plans[0].ID

	if strings.Contains(driverType, "dummy") {
		serviceBindingGuid = "credentialsID"
	}

	t.Logf("start unbinding service %[1]s, with service binding guid %[2]s", service.Name, serviceBindingGuid)

	unbindResp, err := ExecuteHttpCall(
		"DELETE",
		fmt.Sprintf("http://%[1]s:%[2]s@localhost:%[3]v/v2/service_instances/%[4]s/service_bindings/%[5]s?service_id=%[6]s&plan_id=%[7]s", credentials.Username, credentials.Password, brokerApiPort, serviceGuid, serviceBindingGuid, service.ID, servicePlanId),
		nil,
		false)
	if err != nil {
		t.Fatal(err)
	}
	defer unbindResp.Body.Close()

	Expect(unbindResp.StatusCode).To(Equal(200))

	unbindContent, err := ioutil.ReadAll(unbindResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("unbind response content: %s", string(unbindContent))
}

func executeDeprovisionTest(t *testing.T, brokerApiPort uint16, credentials brokerapi.BrokerCredentials, service brokerapi.Service, driverType string) {
	servicePlanId := service.Plans[0].ID

	if driverType == "dummy" {
		serviceGuid = "instanceID"
	}

	t.Logf("start deprovisioning service %[1]s, with service guid %[2]s", service.Name, serviceGuid)

	deprovisionResp, err := ExecuteHttpCall(
		"DELETE",
		fmt.Sprintf("http://%[1]s:%[2]s@localhost:%[3]v/v2/service_instances/%[4]s?service_id=%[5]s&plan_id=%[6]s", credentials.Username, credentials.Password, brokerApiPort, serviceGuid, service.ID, servicePlanId),
		nil,
		false)
	if err != nil {
		t.Fatal(err)
	}
	defer deprovisionResp.Body.Close()

	Expect(deprovisionResp.StatusCode).To(Equal(200))

	deprovisionContent, err := ioutil.ReadAll(deprovisionResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("deprovision response content: %s", string(deprovisionContent))
}
