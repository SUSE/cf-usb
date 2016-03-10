package driver_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hpcloud/cf-usb/lib/config"

	. "github.com/hpcloud/cf-usb/integration_test/consulprovider"
	. "github.com/onsi/gomega"
)

var tmpDriverPath string

func init() {
	ConsulConfig.ConsulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.ConsulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.ConsulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.ConsulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.ConsulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.ConsulToken = os.Getenv("CONSUL_TOKEN")

	tmpDriverPath = path.Join(os.TempDir(), "drivers")
}

func Test_MgmtApi_ConsulProvider_CreateDriver(t *testing.T) {
	RegisterTestingT(t)

	binpath, buildNotExist, err := Check_solutionBuild()
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

	_, ccFakeServer := Set_fakeServers()

	err = Setup_firstConsulRun()
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(tmpDriverPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	consulClient, err := Run_consul(BrokerApiPort, ManagementApiPort, ccFakeServer.URL(), tmpDriverPath)
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

	usbProcess, err := Start_usbProcess(binpath, ConsulConfig.ConsulAddress)
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

	//	err = os.RemoveAll(tmpDriverPath)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
}

func Test_MgmtApi_ConsulProvider_UpdateDriver(t *testing.T) {
	RegisterTestingT(t)

	binpath, buildNotExist, err := Check_solutionBuild()
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

	_, ccFakeServer := Set_fakeServers()

	consulClient, err := Run_consul(BrokerApiPort, ManagementApiPort, ccFakeServer.URL(), tmpDriverPath)
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

	usbProcess, err := Start_usbProcess(binpath, ConsulConfig.ConsulAddress)
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

func Test_MgmtApi_ConsulProvider_DeleteDriver(t *testing.T) {
	RegisterTestingT(t)

	binpath, buildNotExist, err := Check_solutionBuild()
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

	_, ccFakeServer := Set_fakeServers()

	consulClient, err := Run_consul(BrokerApiPort, ManagementApiPort, ccFakeServer.URL(), tmpDriverPath)
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

	usbProcess, err := Start_usbProcess(binpath, ConsulConfig.ConsulAddress)
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
	newDriverResp, err := ExecuteHttpCall("POST", fmt.Sprintf("http://localhost:%[1]v/drivers", managementApiPort), strings.NewReader(fmt.Sprintf(`{"name":"%[1]s", "driver_type":"%[2]s"}`, driverName, driverName)))
	if err != nil {
		t.Fatal(err)
	}
	defer newDriverResp.Body.Close()

	if newDriverResp.StatusCode == 409 {
		t.Skipf("Skipping test as driver type %[1]s already exists", driverName)
	}

	Expect(newDriverResp.StatusCode).To((Equal(201)))

	driverContent, driver, err := UnmarshalDriverResponse(newDriverResp.Body)
	if err != nil {
		t.Fatal(err)
	}

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

	Expect(getDriverResp.StatusCode).To((Equal(200)))

	getDriverContent, err := ioutil.ReadAll(getDriverResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver response content: %s", string(getDriverContent))
	Expect(getDriverContent).To(ContainSubstring(driverId))
}

func executeUploadDriverTest(t *testing.T, managementApiPort uint16, driverType, driverId string) {
	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	// upload driver bits test
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	bitsPath := path.Join(dir, "../../../build", architecture, driverType)
	sha, err := GetFileSha(bitsPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("bits path:", bitsPath)

	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	err = body_writer.WriteField("sha", sha)
	if err != nil {
		t.Log("error writing to buffer")
		t.Fatal(err)
	}

	// use the body_writer to write the Part headers to the buffer
	_, err = body_writer.CreateFormFile("file", bitsPath)
	if err != nil {
		t.Log("error writing to buffer")
		t.Fatal(err)
	}

	// the file data will be the second part of the body
	fh, err := os.Open(bitsPath)
	if err != nil {
		t.Log("error opening file")
		t.Fatal(err)
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	_ = fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	_, err = fh.Stat()
	if err != nil {
		t.Log("Error Stating file: ", bitsPath)
		t.Fatal(err)
	}

	uploadDriverBitsReq, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s/bits", managementApiPort, driverId), request_reader)
	uploadDriverBitsReq.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	uploadDriverBitsReq.Header.Add("Accept", "application/json")
	uploadDriverBitsReq.Header.Add("Authorization", token)

	uploadDriverBitsResp, err := http.DefaultClient.Do(uploadDriverBitsReq)
	if err != nil {
		t.Fatal(err)
	}
	defer uploadDriverBitsResp.Body.Close()

	Expect(uploadDriverBitsResp.StatusCode).To((Equal(200)))

	uploadDriverBitsContent, err := ioutil.ReadAll(uploadDriverBitsResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	Expect(string(uploadDriverBitsContent)).To(Equal(""))

	//check if driver bits file exists
	_, err = os.Stat(path.Join(os.Getenv("USB_DRIVER_PATH"), driverType))
	Expect(err).To(BeNil())
}

func executeGetDriversTest(t *testing.T, managementApiPort uint16) (string, []DriverResponse) {
	getDriversResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/drivers", managementApiPort), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDriversResp.Body.Close()

	Expect(getDriversResp.StatusCode).To((Equal(200)))

	driversContent, drivers, err := UnmarshalDriversResponse(getDriversResp.Body)
	if err != nil {
		t.Fatal(err)
	}

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

	Expect(updateDriverResp.StatusCode).To((Equal(200)))

	updatedDriverContent, err := ioutil.ReadAll(updateDriverResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver response content: %s", string(updatedDriverContent))

	Expect(updatedDriverContent).To(ContainSubstring(updateDriverName))
}

func executeDeleteDriverTest(t *testing.T, managementApiPort uint16, driver DriverResponse) {
	deleteDriverResp, err := ExecuteHttpCall("DELETE", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driver.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDriverResp.Body.Close()

	Expect(deleteDriverResp.StatusCode).To((Equal(204)))

	deleteDriverContent, err := ioutil.ReadAll(deleteDriverResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	Expect(string(deleteDriverContent)).To(Equal(""))

	getDriverDeletedResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driver.Id), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer getDriverDeletedResp.Body.Close()

	Expect(getDriverDeletedResp.StatusCode).To((Equal(404)))
}
