package dialtest

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

func TestMgmtApiConsulProviderCreateDial(t *testing.T) {
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

			_, _, driverInstanceResp, err := CreateDriverInstance(ManagementApiPort, driverResp, driver.SetDriverInstanceValuesFunc)
			if err != nil {
				t.Fatal(err)
			}

			SetupCcHttpFakeResponsesCreateDial(driverInstanceResp.Name, uaaFakeServer, ccFakeServer)

			executeCreateDialTest(t, ManagementApiPort, driverInstanceResp)
			executeGetDialsTest(t, ManagementApiPort, driverInstanceResp)
			executeGetPlansTest(t, ManagementApiPort, driverInstanceResp)
		}
	}
}

func TestMgmtApiConsulProviderUpdateDial(t *testing.T) {
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

					_, _, driverInstancesResponse, err := GetDriverInstances(ManagementApiPort, d.Id)
					if err != nil {
						t.Fatal(err)
					}

					if len(driverInstancesResponse) > 0 {
						firstDriverInstance := driverInstancesResponse[0]

						if len(firstDriverInstance.Dials) > 0 {
							firstDial := executeGetDialTest(t, ManagementApiPort, firstDriverInstance)

							executeUpdateDialTest(t, ManagementApiPort, firstDial)

							plan := executeGetPlanTest(t, ManagementApiPort, firstDial)

							SetupCcHttpFakeResponsesUpdatePlan(uaaFakeServer, ccFakeServer)

							executeUpdatePlanTest(t, ManagementApiPort, plan, firstDial)
						}
					}
				}
			}
		}
	}
}

func TestMgmtApiConsulProviderDeleteDial(t *testing.T) {
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

					_, _, driverInstancesResponse, err := GetDriverInstances(ManagementApiPort, d.Id)
					if err != nil {
						t.Fatal(err)
					}

					if len(driverInstancesResponse) > 0 {
						firstDriverInstance := driverInstancesResponse[0]

						if len(firstDriverInstance.Dials) > 0 {
							firstDial := executeGetDialTest(t, ManagementApiPort, firstDriverInstance)

							executeDeleteDialTest(t, ManagementApiPort, firstDial)

							executeCheckPlanIsDeletedTest(t, ManagementApiPort, firstDial)
						}

						for _, di := range driverInstancesResponse {
							SetupCcHttpFakeResponsesDeleteDriverInstance(di.Name, uaaFakeServer, ccFakeServer)

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

func executeCreateDialTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse) {
	dialValues := []byte(fmt.Sprintf(`{"configuration":{"max_dbsize_mb":3},"driver_instance_id":"%[1]s"}`, firstDriverInstance.Id))

	createDialResp, err := ExecuteHttpCall("POST", fmt.Sprintf("http://localhost:%[1]v/dials", managementApiPort), bytes.NewBuffer(dialValues))
	if err != nil {
		t.Fatal(err)
	}
	defer createDialResp.Body.Close()

	Expect(createDialResp.StatusCode).To(Equal(201))

	createDialContent, err := ioutil.ReadAll(createDialResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("create dial response content: %s", string(createDialContent))
	Expect(string(createDialContent)).To(ContainSubstring(`"max_dbsize_mb":3`))

	var createdDial DialResponse

	err = json.Unmarshal(createDialContent, &createdDial)
	if err != nil {
		t.Fatal(err)
	}
	Expect(createdDial.Plan).NotTo(Equal(""))
}

func executeGetDialsTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse) {
	getDialsResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/dials?driver_instance_id=%[2]s", managementApiPort, firstDriverInstance.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDialsResp.Body.Close()

	Expect(getDialsResp.StatusCode).To((Equal(200)))

	getDialsContent, err := ioutil.ReadAll(getDialsResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get dials response content: %s", string(getDialsContent))
	Expect(getDialsContent).To(ContainSubstring("plan"))

	var dials []DialResponse

	err = json.Unmarshal(getDialsContent, &dials)
	if err != nil {
		t.Fatal(err)
	}
	Expect(len(dials)).To(Equal(2))
	Expect(dials[0].Id).NotTo(Equal(dials[1].Id))
	Expect(dials[0].Plan).NotTo(Equal(dials[1].Plan))
}

func executeGetPlansTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse) {
	getPlansResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/plans?driver_instance_id=%[2]s", managementApiPort, firstDriverInstance.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getPlansResp.Body.Close()

	Expect(getPlansResp.StatusCode).To((Equal(200)))

	getPlansContent, err := ioutil.ReadAll(getPlansResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get plans response content: %s", string(getPlansContent))
	Expect(getPlansContent).To(ContainSubstring(`"name":"default"`))
	Expect(getPlansContent).To(ContainSubstring(`"name":"plan-`))
	Expect(getPlansContent).To(ContainSubstring(`"description":"N/A"`))
	Expect(getPlansContent).To(ContainSubstring(`"free":true`))
	Expect(getPlansContent).To(ContainSubstring(`"free":false`))

	var plans []PlanResponse

	err = json.Unmarshal(getPlansContent, &plans)
	if err != nil {
		t.Fatal(err)
	}
	Expect(len(plans)).To(Equal(2))
	Expect(plans[0].Id).NotTo(Equal(plans[1].Id))
	Expect(plans[0].Description).NotTo(Equal(plans[1].Description))
	Expect(plans[0].Name).NotTo(Equal(plans[1].Name))
	Expect(plans[0].Free).NotTo(Equal(plans[1].Free))
	Expect(plans[0].DialId).NotTo(Equal(plans[1].DialId))
}

func executeGetDialTest(t *testing.T, managementApiPort uint16, firstDriverInstance DriverInstanceResponse) DialResponse {
	getDialResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, firstDriverInstance.Dials[0]), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDialResp.Body.Close()

	Expect(getDialResp.StatusCode).To(Equal(200))

	getDialContent, err := ioutil.ReadAll(getDialResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get dial response content: %s", string(getDialContent))

	var dial DialResponse

	err = json.Unmarshal(getDialContent, &dial)
	if err != nil {
		t.Fatal(err)
	}
	Expect(dial.DriverInstanceId).To(Equal(firstDriverInstance.Id))

	return dial
}

func executeUpdateDialTest(t *testing.T, managementApiPort uint16, dial DialResponse) {
	dialValues := []byte(fmt.Sprintf(`{"configuration":{"min_dbsize_mb":1},"driver_instance_id":"%[1]s"}`, dial.DriverInstanceId))

	updateDialResp, err := ExecuteHttpCall("PUT", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, dial.Id), bytes.NewBuffer(dialValues))
	if err != nil {
		t.Fatal(err)
	}
	defer updateDialResp.Body.Close()

	Expect(updateDialResp.StatusCode).To(Equal(200))

	updateDialContent, err := ioutil.ReadAll(updateDialResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	Expect(string(updateDialContent)).To(ContainSubstring(`"min_dbsize_mb":1`))
}

func executeGetPlanTest(t *testing.T, managementApiPort uint16, dial DialResponse) PlanResponse {
	getPlanResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/plans/%[2]s", managementApiPort, dial.Plan), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getPlanResp.Body.Close()

	Expect(getPlanResp.StatusCode).To((Equal(200)))

	getPlanContent, err := ioutil.ReadAll(getPlanResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get plan response content: %s", string(getPlanContent))

	var plan PlanResponse

	err = json.Unmarshal(getPlanContent, &plan)
	if err != nil {
		t.Fatal(err)
	}
	Expect(plan.Name).To(Or(ContainSubstring("default"), ContainSubstring("plan-")))
	Expect(plan.Description).To(Or(ContainSubstring("N/A"), ContainSubstring("default plan")))

	return plan
}

func executeUpdatePlanTest(t *testing.T, managementApiPort uint16, plan PlanResponse, dial DialResponse) {
	updatePlanName := plan.Name + "updp"
	updatePlanDesc := plan.Description + "updp"

	planValues := []byte(fmt.Sprintf(`{"description":"%[1]s","dial_id":"%[2]s","free":true,"name":"%[3]s"}`,
		updatePlanDesc,
		dial.Id,
		updatePlanName))

	updatePlanResp, err := ExecuteHttpCall("PUT", fmt.Sprintf("http://localhost:%[1]v/plans/%[2]s", managementApiPort, dial.Plan), bytes.NewBuffer(planValues))
	if err != nil {
		t.Fatal(err)
	}
	defer updatePlanResp.Body.Close()

	Expect(updatePlanResp.StatusCode).To(Equal(200))

	updatePlanContent, err := ioutil.ReadAll(updatePlanResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	Expect(updatePlanContent).To(ContainSubstring(updatePlanName))
	Expect(updatePlanContent).To(ContainSubstring(updatePlanDesc))
}

func executeDeleteDialTest(t *testing.T, managementApiPort uint16, dial DialResponse) {
	deleteDialResp, err := ExecuteHttpCall("DELETE", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, dial.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDialResp.Body.Close()

	Expect(deleteDialResp.StatusCode).To(Equal(204))

	deleteDialContent, err := ioutil.ReadAll(deleteDialResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	Expect(string(deleteDialContent)).To(Equal(""))

	getDialDeletedResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, dial.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDialDeletedResp.Body.Close()

	Expect(getDialDeletedResp.StatusCode).To((Equal(404)))
}

func executeCheckPlanIsDeletedTest(t *testing.T, managementApiPort uint16, dial DialResponse) {
	getPlanResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/plans/%[2]s", managementApiPort, dial.Plan), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getPlanResp.Body.Close()

	Expect(getPlanResp.StatusCode).To((Equal(404)))
}
