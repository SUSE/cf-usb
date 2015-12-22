package fileprovider_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/cf-usb/lib/config"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/pivotal-golang/localip"
)

func getBinPath() string {
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	os.Setenv("USB_DRIVER_PATH", path.Join(dir, "../../build", architecture))

	return path.Join(dir, "../../build", architecture, "usb")
}

func initializeRunner() (UsbRunner, *config.Config) {
	freePort, err := localip.LocalPort()
	Expect(err).NotTo(HaveOccurred())

	tempDir, err := ioutil.TempDir("", "cf-usb-test")
	Expect(err).NotTo(HaveOccurred())

	usbRunner := UsbRunner{
		UsbBrokerPort: freePort,
		Path:          getBinPath(),
		TempDir:       tempDir,
	}

	usbRunner.Start()

	provider := config.NewFileConfig(usbRunner.ConfigFile)
	configInfo, err := provider.LoadConfiguration()
	Expect(err).NotTo(HaveOccurred())

	return usbRunner, configInfo
}

func isValidJson(s []byte) bool {
	var m map[string]interface{}
	return json.Unmarshal(s, &m) == nil
}

func Test_BrokerWithFileConfigProviderCatalog(t *testing.T) {
	RegisterTestingT(t)

	binpath := getBinPath()
	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}

	usb, configInfo := initializeRunner()
	defer usb.Stop()

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password

	resp, err := http.Get(fmt.Sprintf("http://%s:%s@%s/v2/catalog", user, pass, usb.BrokerAddress()))
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())

	Expect(content).To(ContainSubstring(configInfo.Drivers["00000000-0000-0000-0000-000000000001"].DriverInstances["A0000000-0000-0000-0000-000000000002"].Service.ID))
}

func Test_BrokerWithFileConfigProviderProvision(t *testing.T) {
	RegisterTestingT(t)

	binpath := getBinPath()
	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}

	usb, configInfo := initializeRunner()
	defer usb.Stop()

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password

	var jsonStr = []byte(`
	{
		"service_id":"de8464a4-1d05-4f25-8a74-9790448d13cd",
		"plan_id": "1a7cc5ee-4a46-4af4-9af5-b6f2b6050ee9",
		"organization_guid": "832160fb-2b79-4565-b919-6dcbb7d60a9f",
		"space_guid": "117c31f4-2831-40e3-aad2-c63056ce7f74"
	}'`)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://%s:%s@%s/v2/service_instances/instance1", user, pass, usb.BrokerAddress()),
		bytes.NewBuffer(jsonStr),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(isValidJson(content)).To(BeTrue())

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))
	Expect(usb.Runner.Buffer()).Should(gbytes.Say("usb.usb-broker.provision"))
}

func Test_BrokerWithFileConfigProviderDeprovision(t *testing.T) {
	RegisterTestingT(t)

	binpath := getBinPath()
	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}

	usb, configInfo := initializeRunner()
	defer usb.Stop()

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password

	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://%s:%s@%s/v2/service_instances/instanceID?service_id=de8464a4-1d05-4f25-8a74-9790448d13cd&plan_id=1a7cc5ee-4a46-4af4-9af5-b6f2b6050ee9", user, pass, usb.BrokerAddress()),
		nil,
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(isValidJson(content)).To(BeTrue())

	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(usb.Runner.Buffer()).Should(gbytes.Say("usb.usb-broker.deprovision"))
}

func Test_BrokerWithFileConfigProviderBind(t *testing.T) {
	RegisterTestingT(t)

	binpath := getBinPath()
	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}

	usb, configInfo := initializeRunner()
	defer usb.Stop()

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password

	var jsonStr = []byte(`
	{
		"service_id":"de8464a4-1d05-4f25-8a74-9790448d13cd",
		"plan_id": "1a7cc5ee-4a46-4af4-9af5-b6f2b6050ee9",
		"organization_guid": "832160fb-2b79-4565-b919-6dcbb7d60a9f",
		"space_guid": "117c31f4-2831-40e3-aad2-c63056ce7f74"
	}'`)
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("http://%s:%s@%s/v2/service_instances/instance1/service_bindings/binding1", user, pass, usb.BrokerAddress()),
		bytes.NewBuffer(jsonStr),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(isValidJson(content)).To(BeTrue())

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))
}

func Test_BrokerWithFileConfigProviderUnbind(t *testing.T) {
	RegisterTestingT(t)

	binpath := getBinPath()
	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}

	usb, configInfo := initializeRunner()
	defer usb.Stop()

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password

	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("http://%s:%s@%s/v2/service_instances/instance1/service_bindings/credentialsID?service_id=de8464a4-1d05-4f25-8a74-9790448d13cd&plan_id=1a7cc5ee-4a46-4af4-9af5-b6f2b6050ee9", user, pass, usb.BrokerAddress()),
		nil,
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	Expect(isValidJson(content)).To(BeTrue())

	Expect(resp.StatusCode).To(Equal(http.StatusOK))
}
