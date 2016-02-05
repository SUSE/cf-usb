package integration_test

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hpcloud/cf-usb/lib/config"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/golang/protobuf/proto" //workaround for godep + gomega
	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/pivotal-golang/localip"
	uuid "github.com/satori/go.uuid"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var loggerSB *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")
var consulProcess ifrit.Process

var uaaSigningKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDHFr+KICms+tuT1OXJwhCUmR2dKVy7psa8xzElSyzqx7oJyfJ1
JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMXqHxf+ZH9BL1gk9Y6kCnbM5R6
0gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBugspULZVNRxq7veq/fzwIDAQAB
AoGBAJ8dRTQFhIllbHx4GLbpTQsWXJ6w4hZvskJKCLM/o8R4n+0W45pQ1xEiYKdA
Z/DRcnjltylRImBD8XuLL8iYOQSZXNMb1h3g5/UGbUXLmCgQLOUUlnYt34QOQm+0
KvUqfMSFBbKMsYBAoQmNdTHBaz3dZa8ON9hh/f5TT8u0OWNRAkEA5opzsIXv+52J
duc1VGyX3SwlxiE2dStW8wZqGiuLH142n6MKnkLU4ctNLiclw6BZePXFZYIK+AkE
xQ+k16je5QJBAN0TIKMPWIbbHVr5rkdUqOyezlFFWYOwnMmw/BKa1d3zp54VP/P8
+5aQ2d4sMoKEOfdWH7UqMe3FszfYFvSu5KMCQFMYeFaaEEP7Jn8rGzfQ5HQd44ek
lQJqmq6CE2BXbY/i34FuvPcKU70HEEygY6Y9d8J3o6zQ0K9SYNu+pcXt4lkCQA3h
jJQQe5uEGJTExqed7jllQ0khFJzLMx0K6tj0NeeIzAaGCQz13oo2sCdeGRHO4aDh
HH6Qlq/6UOV5wP8+GAcCQFgRCcB+hrje8hfEEefHcFpyKH+5g1Eu1k0mLrxK2zd+
4SlotYRHgPCEubokb2S1zfZDWIXW3HmggnGgM949TlY=
-----END RSA PRIVATE KEY-----`

var uaaPublicKey = `-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2d\nKVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMX\nqHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBug\nspULZVNRxq7veq/fzwIDAQAB\n-----END PUBLIC KEY-----`

var drivers = []struct {
	driverType                     string
	envVarsExistFunc               func() bool
	setDriverInstanceValuesFunc    func(driverName, driverId string) []byte
	assertDriverSchemaContainsFunc func(schemaContent string)
}{
	//{"dummy-async", dummyEnvVarsExist, setDummyDriverInstanceValues, assertDummySchemaContains},
	{"postgres", postgresEnvVarsExist, setPostgresDriverInstanceValues, assertPostgresSchemaContains},
	{"mongo", mongoEnvVarsExist, setMongoDriverInstanceValues, assertMongoSchemaContains},
	{"mysql", mysqlEnvVarsExist, setMysqlDriverInstanceValues, assertMysqlSchemaContains},
}

var ConsulConfig = struct {
	consulAddress    string
	consulDatacenter string
	consulUser       string
	consulPassword   string
	consulSchema     string
	consulToken      string
}{}

type DriverResponse struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	DriverType      string   `json:"driver_type"`
	DriverInstances []string `json:"driver_instances,omitempty"`
}

type DriverInstanceResponse struct {
	Configuration interface{} `json:"configuration,omitempty"`
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	Service       string      `json:"service"`
	Dials         []string    `json:"dials,omitempty"`
}

func init() {
	ConsulConfig.consulAddress = os.Getenv("CONSUL_ADDRESS")
	ConsulConfig.consulDatacenter = os.Getenv("CONSUL_DATACENTER")
	ConsulConfig.consulPassword = os.Getenv("CONSUL_PASSWORD")
	ConsulConfig.consulUser = os.Getenv("CONSUL_USER")
	ConsulConfig.consulSchema = os.Getenv("CONSUL_SCHEMA")
	ConsulConfig.consulToken = os.Getenv("CONSUL_TOKEN")
}

func init_consulProvisioner() (consul.ConsulProvisionerInterface, error) {
	var consulConfig api.Config
	consulConfig.Address = ConsulConfig.consulAddress
	consulConfig.Datacenter = ConsulConfig.consulPassword

	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	buildDir := filepath.Join(workDir, "../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	os.Setenv("USB_DRIVER_PATH", buildDir)

	var auth api.HttpBasicAuth
	auth.Username = ConsulConfig.consulUser
	auth.Password = ConsulConfig.consulPassword

	consulConfig.HttpAuth = &auth
	consulConfig.Scheme = ConsulConfig.consulSchema

	consulConfig.Token = ConsulConfig.consulToken

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		return nil, err
	}
	return provisioner, nil
}

func start_usbProcess(binPath, consulAddress string) (ifrit.Process, error) {
	usbRunner := ginkgomon.New(ginkgomon.Config{
		Name:              "cf-usb",
		Command:           exec.Command(binPath, "--loglevel", "debug", "consulConfigProvider", "-a", consulAddress),
		StartCheck:        "usb.start-listening-brokerapi",
		StartCheckTimeout: 5 * time.Second,
		Cleanup:           func() {},
	})

	usbProcess := ginkgomon.Invoke(usbRunner)

	// wait for the processes to start before returning
	<-usbProcess.Ready()

	return usbProcess, nil
}

func start_consulProcess() (ifrit.Process, error) {

	defaultConsulPath := "consul"

	tmpConsul := path.Join(os.TempDir(), "consul")

	if _, err := os.Stat(tmpConsul); err == nil {
		err := os.RemoveAll(tmpConsul)
		if err != nil {
			return nil, err
		}
	}

	err := os.MkdirAll(tmpConsul, 0755)
	if err != nil {
		return nil, err
	}

	tempPath, err := ioutil.TempDir(tmpConsul, "")
	if err != nil {
		return nil, err
	}

	consulRunner := ginkgomon.New(ginkgomon.Config{
		Name:              "consul",
		Command:           exec.Command(defaultConsulPath, "agent", "-server", "-bootstrap-expect", "1", "-data-dir", tempPath, "-advertise", "127.0.0.1"),
		AnsiColorCode:     "",
		StartCheck:        "New leader elected",
		StartCheckTimeout: 5 * time.Second,
		Cleanup:           func() {},
	})

	consulProcess := ginkgomon.Invoke(consulRunner)

	// wait for the processes to start before returning
	<-consulProcess.Ready()

	return consulProcess, nil
}

func run_consul(brokerApiPort, managementApiPort uint16, ccServerUrl string) (consul.ConsulProvisionerInterface, error) {
	var consulClient consul.ConsulProvisionerInterface

	getConsulReq, _ := http.NewRequest("GET", "http://localhost:8500", nil)
	getConsulResp, _ := http.DefaultClient.Do(getConsulReq)
	consulIsRunning := false
	if getConsulResp != nil && getConsulResp.StatusCode == 200 {
		consulIsRunning = true
	}

	if (strings.Contains(ConsulConfig.consulAddress, "127.0.0.1") || strings.Contains(ConsulConfig.consulAddress, "localhost")) && !consulIsRunning {
		ConsulConfig.consulAddress = "127.0.0.1:8500"
		ConsulConfig.consulSchema = "http"

		consulProcess, err := start_consulProcess()
		if err != nil {
			return nil, err
		}

		defer func() {
			consulProcess.Signal(os.Kill)
			<-consulProcess.Wait()
		}()
	}

	consulClient, err := init_consulProvisioner()
	if err != nil {
		return nil, err
	}

	var list api.KVPairs

	list = append(list, &api.KVPair{Key: "usb/api_version", Value: []byte("2.0")})
	list = append(list, &api.KVPair{Key: "usb/broker_api", Value: []byte(fmt.Sprintf("{\"listen\":\":%v\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", brokerApiPort))})
	list = append(list, &api.KVPair{Key: "usb/management_api", Value: []byte(fmt.Sprintf("{\"listen\":\":%v\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"%v\"}},\"cloud_controller\":{\"api\":\"%s\",\"skip_tsl_validation\":true}}", managementApiPort, uaaPublicKey, ccServerUrl))})

	err = consulClient.PutKVs(&list, nil)
	if err != nil {
		return consulClient, err
	}

	return consulClient, nil
}

func set_fakeServers() (*ghttp.Server, *ghttp.Server) {
	var uaaFakeServer *ghttp.Server
	uaaFakeServer = ghttp.NewServer()

	var ccFakeServer *ghttp.Server
	ccFakeServer = ghttp.NewServer()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/info"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(
					`{"name": "vcap","authorization_endpoint": "%[1]s","token_endpoint":"%[1]s","api_version":"2.44.0"}`,
					uaaFakeServer.URL()),
			),
		),
	)

	return uaaFakeServer, ccFakeServer
}

func GenerateUaaToken() (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	token.Header = map[string]interface{}{
		"alg": "RS256",
	}

	token.Claims = map[string]interface{}{
		"exp":   3404281214,
		"scope": []string{"usb.management.admin"},
	}

	signedKey, err := token.SignedString([]byte(uaaSigningKey))
	if err != nil {
		return "", err
	}

	return "bearer " + signedKey, nil
}

func Test_BrokerWithConsulConfigProviderCatalog(t *testing.T) {
	RegisterTestingT(t)
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binpath := path.Join(dir, "../build", architecture, "usb")

	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}

	if ConsulConfig.consulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}
	if err != nil {
		t.Fatal(err)
	}

	_, ccFakeServer := set_fakeServers()

	brokerApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())
	managementApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())

	consulClient, err := run_consul(brokerApiPort, managementApiPort, ccFakeServer.URL())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("consul started")

	usbProcess, err := start_usbProcess(binpath, ConsulConfig.consulAddress)
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

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password

	resp, err := http.Get(fmt.Sprintf("http://%s:%s@%s/v2/catalog", user, pass, "localhost:"+strconv.Itoa(int(brokerApiPort))))

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Test catalog:")
	t.Log(string(content))
}

// create driver, create driver instance, ping tests
func Test_BrokerWithConsulConfigProviderCreateDriverInstance(t *testing.T) {
	RegisterTestingT(t)
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binpath := path.Join(dir, "../build", architecture, "usb")

	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	if ConsulConfig.consulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}

	uaaFakeServer, ccFakeServer := set_fakeServers()

	brokerApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())
	managementApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())

	consulClient, err := run_consul(brokerApiPort, managementApiPort, ccFakeServer.URL())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("consul started")

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	usbProcess, err := start_usbProcess(binpath, ConsulConfig.consulAddress)
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

	for _, driver := range drivers {
		if driver.envVarsExistFunc() {
			setupCcHttpFakeResponsesCreateDriverInstance(uaaFakeServer, ccFakeServer)
			executeCreateDriverInstanceTest(t, managementApiPort, driver.driverType, driver.setDriverInstanceValuesFunc)
		}
	}
}

// update plan and update driver instance tests
func Test_BrokerWithConsulConfigProviderUpdateDriverInstance(t *testing.T) {
	RegisterTestingT(t)
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binpath := path.Join(dir, "../build", architecture, "usb")

	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	if ConsulConfig.consulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}

	uaaFakeServer, ccFakeServer := set_fakeServers()

	brokerApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())
	managementApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())

	consulClient, err := run_consul(brokerApiPort, managementApiPort, ccFakeServer.URL())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("consul started")

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	usbProcess, err := start_usbProcess(binpath, ConsulConfig.consulAddress)
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

	driversContent, driversResp := executeGetDriversTest(t, managementApiPort)
	brokerGuid := uuid.NewV4().String()

	for _, driver := range drivers {
		if driver.envVarsExistFunc() {

			Expect(driversContent).To(ContainSubstring(driver.driverType))

			for _, d := range driversResp {
				if d.DriverType == driver.driverType {
					setupCcHttpFakeResponsesUpdatePlan(brokerGuid, uaaFakeServer, ccFakeServer)
					executeTestUpdateDriverInstance(t, managementApiPort, d, driver.setDriverInstanceValuesFunc)
				}
			}
		}
	}
}

// update service and update driver tests
func Test_BrokerWithConsulConfigProviderUpdateService(t *testing.T) {
	RegisterTestingT(t)
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binpath := path.Join(dir, "../build", architecture, "usb")

	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	if ConsulConfig.consulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}

	uaaFakeServer, ccFakeServer := set_fakeServers()

	brokerApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())
	managementApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())

	consulClient, err := run_consul(brokerApiPort, managementApiPort, ccFakeServer.URL())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("consul started")

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	usbProcess, err := start_usbProcess(binpath, ConsulConfig.consulAddress)
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

	driversContent, driversResp := executeGetDriversTest(t, managementApiPort)
	brokerGuid := uuid.NewV4().String()

	for _, driver := range drivers {
		if driver.envVarsExistFunc() {

			Expect(driversContent).To(ContainSubstring(driver.driverType))

			for _, d := range driversResp {
				if d.DriverType == driver.driverType {
					setupCcHttpFakeResponsesUpdateService(brokerGuid, uaaFakeServer, ccFakeServer)
					executeTestUpdateService(t, managementApiPort, d, driver.setDriverInstanceValuesFunc)

					executeTestUpdateDriver(t, managementApiPort, d, driver.assertDriverSchemaContainsFunc)
				}
			}
		}
	}
}

// delete plan, dial, driver instance, driver tests
func Test_BrokerWithConsulConfigProviderDeleteDriver(t *testing.T) {
	RegisterTestingT(t)
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	binpath := path.Join(dir, "../build", architecture, "usb")

	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	if ConsulConfig.consulAddress == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS")
	}

	uaaFakeServer, ccFakeServer := set_fakeServers()

	brokerApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())
	managementApiPort, err := localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())

	consulClient, err := run_consul(brokerApiPort, managementApiPort, ccFakeServer.URL())
	if err != nil {
		t.Fatal(err)
	}
	t.Log("consul started")

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	usbProcess, err := start_usbProcess(binpath, ConsulConfig.consulAddress)
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

	driversContent, driversResp := executeGetDriversTest(t, managementApiPort)
	brokerGuid := uuid.NewV4().String()

	for _, driver := range drivers {
		if driver.envVarsExistFunc() {

			Expect(driversContent).To(ContainSubstring(driver.driverType))

			for _, d := range driversResp {
				if d.DriverType == driver.driverType {
					setupCcHttpFakeResponsesDeleteDriver(brokerGuid, uaaFakeServer, ccFakeServer)
					executeTestDeleteDriver(t, managementApiPort, d)
				}
			}
		}
	}
}

func setupCcHttpFakeResponsesCreateDriverInstance(uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				`{"resources":[]}`),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{"resources":[]}`),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/v2/service_brokers"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(201, `{}`),
		),
	)

	serviceGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, serviceGuid)),
		),
	)

	servicePlanGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_plans"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"},"entity":{"name":"default","free":true,"description":"default plan","public":false,"service_guid":"%[2]s"}}]}`, servicePlanGuid, serviceGuid)),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("/v2/service_plans/%[1]s", servicePlanGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(201, `{}`),
		),
	)
}

func executeCreateDriverInstanceTest(t *testing.T, managementApiPort uint16, driverName string, driverInstanceValues func(driverName, driverId string) []byte) {
	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	newDriverReq, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%[1]v/drivers", managementApiPort), strings.NewReader(fmt.Sprintf(`{"name":"%[1]s", "driver_type":"%[2]s"}`, driverName, driverName)))
	newDriverReq.Header.Add("Content-Type", "application/json")
	newDriverReq.Header.Add("Accept", "application/json")
	newDriverReq.Header.Add("Authorization", token)

	newDriverResp, err := http.DefaultClient.Do(newDriverReq)

	if err != nil {
		t.Fatal(err)
	}
	defer newDriverResp.Body.Close()

	if newDriverResp.StatusCode == 409 {
		t.Logf("Skipping test as driver type %[1]s already exists", driverName)
		return
	}

	driverContent, err := ioutil.ReadAll(newDriverResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var driver DriverResponse

	err = json.Unmarshal(driverContent, &driver)
	if err != nil {
		fmt.Println("error:", err)
	}
	t.Logf("create driver response content: %s", string(driverContent))

	Expect(driver.Id).ToNot(BeNil())
	Expect(driver.Name).To(Equal(driverName))
	driverId := driver.Id

	getDriverReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driverId), nil)
	getDriverReq.Header.Add("Content-Type", "application/json")
	getDriverReq.Header.Add("Accept", "application/json")
	getDriverReq.Header.Add("Authorization", token)

	getDriverResp, err := http.DefaultClient.Do(getDriverReq)

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

	// upload driver bits test
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	bitsPath := path.Join(dir, "../cmd/driver", driver.DriverType, driver.DriverType)
	sha, err := getFileSha(bitsPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("bits path:", bitsPath)
	t.Log("sha:", sha)

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

	instanceValues := driverInstanceValues(driverName, driverId)
	newDriverInstReq, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%[1]v/driver_instances", managementApiPort), bytes.NewBuffer(instanceValues))
	newDriverInstReq.Header.Add("Content-Type", "application/json")
	newDriverInstReq.Header.Add("Accept", "application/json")
	newDriverInstReq.Header.Add("Authorization", token)

	newDriverInstResp, err := http.DefaultClient.Do(newDriverInstReq)

	if err != nil {
		t.Fatal(err)
	}
	defer newDriverInstResp.Body.Close()

	Expect(newDriverInstResp.StatusCode).To((Equal(201)))

	driverInstContent, err := ioutil.ReadAll(newDriverInstResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("driver instance: %s", string(driverInstContent))

	var driverInstance DriverInstanceResponse

	err = json.Unmarshal(driverInstContent, &driverInstance)
	if err != nil {
		fmt.Println("error:", err)
	}
	Expect(driverInstContent).To(ContainSubstring(driver.Name))

	pingDriverInstReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s/ping", managementApiPort, driverInstance.Id), nil)
	pingDriverInstReq.Header.Add("Content-Type", "application/json")
	pingDriverInstReq.Header.Add("Accept", "application/json")
	pingDriverInstReq.Header.Add("Authorization", token)

	pingDriverInstResp, err := http.DefaultClient.Do(pingDriverInstReq)

	if err != nil {
		t.Fatal(err)
	}
	defer pingDriverInstResp.Body.Close()

	Expect(pingDriverInstResp.StatusCode).To((Equal(200)))

	driverInstPingContent, err := ioutil.ReadAll(pingDriverInstResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	Expect(string(driverInstPingContent)).To(Equal(""))

	getPlanReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/plans?driver_instance_id=%[2]s", managementApiPort, driverInstance.Id), nil)
	getPlanReq.Header.Add("Content-Type", "application/json")
	getPlanReq.Header.Add("Accept", "application/json")
	getPlanReq.Header.Add("Authorization", token)

	getPlanResp, err := http.DefaultClient.Do(getPlanReq)

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
	Expect(getPlanContent).To(ContainSubstring("default"))

	getDialReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/dials?driver_instance_id=%[2]s", managementApiPort, driverInstance.Id), nil)
	getDialReq.Header.Add("Content-Type", "application/json")
	getDialReq.Header.Add("Accept", "application/json")
	getDialReq.Header.Add("Authorization", token)

	getDialResp, err := http.DefaultClient.Do(getDialReq)
	if err != nil {
		t.Fatal(err)
	}
	defer getDialResp.Body.Close()

	Expect(getDialResp.StatusCode).To((Equal(200)))

	getDialContent, err := ioutil.ReadAll(getDialResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get dial response content: %s", string(getDialContent))
	Expect(getDialContent).To(ContainSubstring("plan"))

	getServiceReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/services?driver_instance_id=%[2]s", managementApiPort, driverInstance.Id), nil)
	getServiceReq.Header.Add("Content-Type", "application/json")
	getServiceReq.Header.Add("Accept", "application/json")
	getServiceReq.Header.Add("Authorization", token)

	getServiceResp, err := http.DefaultClient.Do(getServiceReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getServiceResp.Body.Close()

	Expect(getServiceResp.StatusCode).To((Equal(200)))

	getServiceContent, err := ioutil.ReadAll(getServiceResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get service response content: %s", string(getServiceContent))
	Expect(getServiceContent).To(ContainSubstring(driverInstance.Id))
}

func executeGetDriversTest(t *testing.T, managementApiPort uint16) (string, []DriverResponse) {
	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	getDriversReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/drivers", managementApiPort), nil)
	getDriversReq.Header.Add("Content-Type", "application/json")
	getDriversReq.Header.Add("Accept", "application/json")
	getDriversReq.Header.Add("Authorization", token)

	getDriversResp, err := http.DefaultClient.Do(getDriversReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getDriversResp.Body.Close()

	Expect(getDriversResp.StatusCode).To((Equal(200)))

	driversContent, err := ioutil.ReadAll(getDriversResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get drivers response content: %s", string(driversContent))

	var drivers []DriverResponse

	err = json.Unmarshal(driversContent, &drivers)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(driversContent), drivers
}

func setupCcHttpFakeResponsesUpdatePlan(brokerGuid string, uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	ccFakeServer.RouteToHandler("GET", "/v2/service_brokers",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			ghttp.RespondWith(200, fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, brokerGuid)),
		),
	)

	ccFakeServer.RouteToHandler("PUT", "/v2/service_brokers/"+brokerGuid,
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/v2/service_brokers/"+brokerGuid),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{}`),
		),
	)

	serviceGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, serviceGuid)),
		),
	)

	servicePlanGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_plans"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"},"entity":{"name":"default","free":true,"description":"default plan","public":false,"service_guid":"%[2]s"}}]}`, servicePlanGuid, serviceGuid)),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("/v2/service_plans/%[1]s", servicePlanGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(201, `{}`),
		),
	)
}

func executeTestUpdateDriverInstance(t *testing.T, managementApiPort uint16, driver DriverResponse, driverInstanceValues func(driverName, driverId string) []byte) {
	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	getDriverInstancesReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances?driver_id=%[2]s", managementApiPort, driver.Id), nil)
	getDriverInstancesReq.Header.Add("Content-Type", "application/json")
	getDriverInstancesReq.Header.Add("Accept", "application/json")
	getDriverInstancesReq.Header.Add("Authorization", token)

	getDriverInstancesResp, err := http.DefaultClient.Do(getDriverInstancesReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getDriverInstancesResp.Body.Close()

	Expect(getDriverInstancesResp.StatusCode).To((Equal(200)))

	driverInstancesContent, err := ioutil.ReadAll(getDriverInstancesResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver instances response content: %s", string(driverInstancesContent))

	var driverInstances []DriverInstanceResponse

	err = json.Unmarshal(driverInstancesContent, &driverInstances)
	if err != nil {
		fmt.Println("error:", err)
	}
	t.Logf("driver instances count: %v", len(driverInstances))

	firstDriverInstance := driverInstances[0]

	getDialReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, firstDriverInstance.Dials[0]), nil)
	getDialReq.Header.Add("Content-Type", "application/json")
	getDialReq.Header.Add("Accept", "application/json")
	getDialReq.Header.Add("Authorization", token)

	getDialResp, err := http.DefaultClient.Do(getDialReq)
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

	type DialResponse struct {
		Configuration    interface{} `json:"configuration,omitempty"`
		Id               string      `json:"id"`
		DriverInstanceId string      `json:"driver_instance_id"`
		Plan             string      `json:"plan"`
	}

	var dial DialResponse

	err = json.Unmarshal(getDialContent, &dial)
	if err != nil {
		fmt.Println("error:", err)
	}
	Expect(dial.DriverInstanceId).To(Equal(firstDriverInstance.Id))

	dialValues := []byte(fmt.Sprintf(`{"configuration":{"min_dbsize_mb":1},"driver_instance_id":"%[1]s","id":"%[2]s"}`,
		dial.DriverInstanceId,
		dial.Id))

	updateDialReq, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, dial.Id), bytes.NewBuffer(dialValues))
	updateDialReq.Header.Add("Content-Type", "application/json")
	updateDialReq.Header.Add("Accept", "application/json")
	updateDialReq.Header.Add("Authorization", token)

	updateDialResp, err := http.DefaultClient.Do(updateDialReq)
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

	getPlanReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/plans/%[2]s", managementApiPort, dial.Plan), nil)
	getPlanReq.Header.Add("Content-Type", "application/json")
	getPlanReq.Header.Add("Accept", "application/json")
	getPlanReq.Header.Add("Authorization", token)

	getPlanResp, err := http.DefaultClient.Do(getPlanReq)

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

	type PlanResponse struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Id          string `json:"id"`
		DialId      string `json:"dial_id"`
		Free        bool   `json:"free"`
	}

	var plan PlanResponse

	err = json.Unmarshal(getPlanContent, &plan)
	if err != nil {
		fmt.Println("error:", err)
	}
	Expect(plan.Name).To(ContainSubstring("default"))
	Expect(plan.Free).To(Equal(true))

	updatePlanName := plan.Name + "updp"
	updatePlanDesc := plan.Description + "updp"

	planValues := []byte(fmt.Sprintf(`{"description":"%[1]s","dial_id":"%[2]s","id":"%[3]s","free":true,"name":"%[4]s"}`,
		updatePlanDesc,
		dial.Id,
		plan.Id,
		updatePlanName))

	updatePlanReq, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:%[1]v/plans/%[2]s", managementApiPort, dial.Plan), bytes.NewBuffer(planValues))
	updatePlanReq.Header.Add("Content-Type", "application/json")
	updatePlanReq.Header.Add("Accept", "application/json")
	updatePlanReq.Header.Add("Authorization", token)

	updatePlanResp, err := http.DefaultClient.Do(updatePlanReq)

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

	// negative test update driver instance wrong configuration
	newDriverInstanceNameNeg := firstDriverInstance.Name + "updierr"

	instanceValuesNeg := []byte(fmt.Sprintf(`{"name":"%[1]s", "driver_id":"%[2]s", "configuration": {"aKey": "aValue"}}`,
		newDriverInstanceNameNeg,
		driver.Id))

	updateDriverInstNegReq, err := http.NewRequest("PUT",
		fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, firstDriverInstance.Id),
		bytes.NewBuffer(instanceValuesNeg))
	updateDriverInstNegReq.Header.Add("Content-Type", "application/json")
	updateDriverInstNegReq.Header.Add("Accept", "application/json")
	updateDriverInstNegReq.Header.Add("Authorization", token)

	updateDriverInstNegResp, err := http.DefaultClient.Do(updateDriverInstNegReq)

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
	Expect(string(updateDriverInstNegContent)).To(ContainSubstring("Invalid configuration schema"))

	newDriverInstanceName := firstDriverInstance.Name + "updi"

	instanceValues := driverInstanceValues(newDriverInstanceName, driver.Id)

	updateDriverInstReq, err := http.NewRequest("PUT",
		fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, firstDriverInstance.Id),
		bytes.NewBuffer(instanceValues))
	updateDriverInstReq.Header.Add("Content-Type", "application/json")
	updateDriverInstReq.Header.Add("Accept", "application/json")
	updateDriverInstReq.Header.Add("Authorization", token)

	updateDriverInstResp, err := http.DefaultClient.Do(updateDriverInstReq)

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

func setupCcHttpFakeResponsesUpdateService(brokerGuid string, uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				`{"resources":[]}`),
		),
	)

	ccFakeServer.RouteToHandler("GET", "/v2/service_brokers",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			ghttp.RespondWith(200, fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, brokerGuid)),
		),
	)

	ccFakeServer.RouteToHandler("PUT", "/v2/service_brokers/"+brokerGuid,
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/v2/service_brokers/"+brokerGuid),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{}`),
		),
	)

	serviceGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, serviceGuid)),
		),
	)

	servicePlanGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_plans"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"},"entity":{"name":"default","free":true,"description":"default plan","public":false,"service_guid":"%[2]s"}}]}`, servicePlanGuid, serviceGuid)),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("/v2/service_plans/%[1]s", servicePlanGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(201, `{}`),
		),
	)
}

func executeTestUpdateService(t *testing.T, managementApiPort uint16, driver DriverResponse, driverInstanceValues func(driverName, driverId string) []byte) {
	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	getDriverInstancesReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances?driver_id=%[2]s", managementApiPort, driver.Id), nil)
	getDriverInstancesReq.Header.Add("Content-Type", "application/json")
	getDriverInstancesReq.Header.Add("Accept", "application/json")
	getDriverInstancesReq.Header.Add("Authorization", token)

	getDriverInstancesResp, err := http.DefaultClient.Do(getDriverInstancesReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getDriverInstancesResp.Body.Close()

	Expect(getDriverInstancesResp.StatusCode).To((Equal(200)))

	driverInstancesContent, err := ioutil.ReadAll(getDriverInstancesResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver instances response content: %s", string(driverInstancesContent))

	var driverInstances []DriverInstanceResponse

	err = json.Unmarshal(driverInstancesContent, &driverInstances)
	if err != nil {
		fmt.Println("error:", err)
	}
	t.Logf("driver instances count: %v", len(driverInstances))

	firstDriverInstance := driverInstances[0]

	getServiceReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/services/%[2]s", managementApiPort, firstDriverInstance.Service), nil)
	getServiceReq.Header.Add("Content-Type", "application/json")
	getServiceReq.Header.Add("Accept", "application/json")
	getServiceReq.Header.Add("Authorization", token)

	getServiceResp, err := http.DefaultClient.Do(getServiceReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getServiceResp.Body.Close()

	Expect(getServiceResp.StatusCode).To((Equal(200)))

	getServiceContent, err := ioutil.ReadAll(getServiceResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get service response content: %s", string(getServiceContent))

	type ServiceResponse struct {
		Bindable         bool        `json:"bindable"`
		Description      string      `json:"description"`
		DriverInstanceID string      `json:"driver_instance_id"`
		ID               string      `json:"id"`
		Metadata         interface{} `json:"metadata"`
		Name             string      `json:"name"`
		Tags             []string    `json:"tags"`
	}

	var service ServiceResponse

	err = json.Unmarshal(getServiceContent, &service)
	if err != nil {
		fmt.Println("error:", err)
	}
	Expect(service.Description).To(ContainSubstring("Default"))
	Expect(service.Bindable).To(Equal(true))

	updateServiceName := service.Name + "upds"
	updateServiceDesc := service.Description + "upds"

	serviceValues := []byte(fmt.Sprintf(`{"bindable":%[1]v,"description":"%[2]s","driver_instance_id":"%[3]s","id":"%[4]s","name":"%[5]s"}`,
		service.Bindable,
		updateServiceDesc,
		service.DriverInstanceID,
		service.ID,
		updateServiceName))

	updateServiceReq, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:%[1]v/services/%[2]s", managementApiPort, firstDriverInstance.Service), bytes.NewBuffer(serviceValues))
	updateServiceReq.Header.Add("Content-Type", "application/json")
	updateServiceReq.Header.Add("Accept", "application/json")
	updateServiceReq.Header.Add("Authorization", token)

	updateServiceResp, err := http.DefaultClient.Do(updateServiceReq)

	if err != nil {
		t.Fatal(err)
	}
	defer updateServiceResp.Body.Close()

	Expect(updateServiceResp.StatusCode).To(Equal(200))

	updateServiceContent, err := ioutil.ReadAll(updateServiceResp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("updated service content: %s", string(updateServiceContent))

	Expect(updateServiceContent).To(ContainSubstring(updateServiceName))
	Expect(updateServiceContent).To(ContainSubstring(updateServiceDesc))
}

func executeTestUpdateDriver(t *testing.T, managementApiPort uint16, driver DriverResponse, assertDriverSchemaContains func(schemaContent string)) {
	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	getDialSchemaReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s/dial_schema", managementApiPort, driver.Id), nil)
	getDialSchemaReq.Header.Add("Content-Type", "application/json")
	getDialSchemaReq.Header.Add("Accept", "application/json")
	getDialSchemaReq.Header.Add("Authorization", token)

	getDialSchemaResp, err := http.DefaultClient.Do(getDialSchemaReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getDialSchemaResp.Body.Close()

	Expect(getDialSchemaResp.StatusCode).To((Equal(200)))

	getDialSchemaContent, err := ioutil.ReadAll(getDialSchemaResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver dial schema content: %s", string(getDialSchemaContent))
	Expect(getDialSchemaContent).To(ContainSubstring("{}"))

	getConfigSchemaReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s/config_schema", managementApiPort, driver.Id), nil)
	getConfigSchemaReq.Header.Add("Content-Type", "application/json")
	getConfigSchemaReq.Header.Add("Accept", "application/json")
	getConfigSchemaReq.Header.Add("Authorization", token)

	getConfigSchemaResp, err := http.DefaultClient.Do(getConfigSchemaReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getConfigSchemaResp.Body.Close()

	Expect(getConfigSchemaResp.StatusCode).To((Equal(200)))

	getConfigSchemaContent, err := ioutil.ReadAll(getConfigSchemaResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver config schema content: %s", string(getConfigSchemaContent))
	assertDriverSchemaContains(string(getConfigSchemaContent))

	updateDriverName := driver.Name + "updd"

	driverValues := []byte(fmt.Sprintf(`{"driver_type":"%[1]s","id":"%[2]s","name":"%[3]s"}`,
		driver.DriverType,
		driver.Id,
		updateDriverName))

	updateDriverReq, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driver.Id), bytes.NewBuffer(driverValues))
	updateDriverReq.Header.Add("Content-Type", "application/json")
	updateDriverReq.Header.Add("Accept", "application/json")
	updateDriverReq.Header.Add("Authorization", token)

	updateDriverResp, err := http.DefaultClient.Do(updateDriverReq)

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

func setupCcHttpFakeResponsesDeleteDriver(brokerGuid string, uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	ccFakeServer.RouteToHandler("GET", "/v2/service_brokers",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			ghttp.RespondWith(200, fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, brokerGuid)),
		),
	)

	ccFakeServer.RouteToHandler("PUT", "/v2/service_brokers/"+brokerGuid,
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", "/v2/service_brokers/"+brokerGuid),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{}`),
		),
	)

	ccFakeServer.RouteToHandler("DELETE", "/v2/service_brokers/"+brokerGuid,
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("DELETE", "/v2/service_brokers/"+brokerGuid),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(204, `{}`),
		),
	)

	serviceGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, serviceGuid)),
		),
	)

	servicePlanGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_plans"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"},"entity":{"name":"default","free":true,"description":"default plan","public":false,"service_guid":"%[2]s"}}]}`, servicePlanGuid, serviceGuid)),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("/v2/service_plans/%[1]s", servicePlanGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(201, `{}`),
		),
	)
}

func executeTestDeleteDriver(t *testing.T, managementApiPort uint16, driver DriverResponse) {
	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	getDriverInstancesReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances?driver_id=%[2]s", managementApiPort, driver.Id), nil)
	getDriverInstancesReq.Header.Add("Content-Type", "application/json")
	getDriverInstancesReq.Header.Add("Accept", "application/json")
	getDriverInstancesReq.Header.Add("Authorization", token)

	getDriverInstancesResp, err := http.DefaultClient.Do(getDriverInstancesReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getDriverInstancesResp.Body.Close()

	Expect(getDriverInstancesResp.StatusCode).To((Equal(200)))

	driverInstancesContent, err := ioutil.ReadAll(getDriverInstancesResp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("get driver instances response content: %s", string(driverInstancesContent))

	var driverInstances []DriverInstanceResponse

	err = json.Unmarshal(driverInstancesContent, &driverInstances)
	if err != nil {
		fmt.Println("error:", err)
	}
	t.Logf("driver instances count: %v", len(driverInstances))

	if len(driverInstances) > 0 {
		firstDriverInstance := driverInstances[0]

		getDialReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, firstDriverInstance.Dials[0]), nil)
		getDialReq.Header.Add("Content-Type", "application/json")
		getDialReq.Header.Add("Accept", "application/json")
		getDialReq.Header.Add("Authorization", token)

		getDialResp, err := http.DefaultClient.Do(getDialReq)
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

		type DialResponse struct {
			Configuration    interface{} `json:"configuration,omitempty"`
			Id               string      `json:"id"`
			DriverInstanceId string      `json:"driver_instance_id"`
			Plan             string      `json:"plan"`
		}

		var dial DialResponse

		err = json.Unmarshal(getDialContent, &dial)
		if err != nil {
			fmt.Println("error:", err)
		}
		Expect(dial.DriverInstanceId).To(Equal(firstDriverInstance.Id))

		// delete plan
		deletePlanReq, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:%[1]v/plans/%[2]s", managementApiPort, dial.Plan), nil)
		deletePlanReq.Header.Add("Content-Type", "application/json")
		deletePlanReq.Header.Add("Accept", "application/json")
		deletePlanReq.Header.Add("Authorization", token)

		deletePlanResp, err := http.DefaultClient.Do(deletePlanReq)

		if err != nil {
			t.Fatal(err)
		}
		defer deletePlanResp.Body.Close()

		Expect(deletePlanResp.StatusCode).To((Equal(204)))

		deletePlanContent, err := ioutil.ReadAll(deletePlanResp.Body)
		if err != nil {
			t.Fatal(err)
		}
		Expect(string(deletePlanContent)).To(Equal(""))

		getPlanReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/plans/%[2]s", managementApiPort, dial.Plan), nil)
		getPlanReq.Header.Add("Content-Type", "application/json")
		getPlanReq.Header.Add("Accept", "application/json")
		getPlanReq.Header.Add("Authorization", token)

		getPlanResp, err := http.DefaultClient.Do(getPlanReq)

		if err != nil {
			t.Fatal(err)
		}
		defer getPlanResp.Body.Close()

		Expect(getPlanResp.StatusCode).To((Equal(404)))

		//check dial is deleted too
		deleteDialReq, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:%[1]v/dials/%[2]s", managementApiPort, dial.Id), nil)
		deleteDialReq.Header.Add("Content-Type", "application/json")
		deleteDialReq.Header.Add("Accept", "application/json")
		deleteDialReq.Header.Add("Authorization", token)

		deleteDialResp, err := http.DefaultClient.Do(deleteDialReq)
		if err != nil {
			t.Fatal(err)
		}
		defer deleteDialResp.Body.Close()

		Expect(deleteDialResp.StatusCode).To(Equal(404))

		deleteDialContent, err := ioutil.ReadAll(deleteDialResp.Body)
		if err != nil {
			t.Fatal(err)
		}
		Expect(string(deleteDialContent)).To(Equal(""))

		// delete driver instance
		deleteDriverInstReq, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, firstDriverInstance.Id), nil)
		deleteDriverInstReq.Header.Add("Content-Type", "application/json")
		deleteDriverInstReq.Header.Add("Accept", "application/json")
		deleteDriverInstReq.Header.Add("Authorization", token)

		deleteDriverInstResp, err := http.DefaultClient.Do(deleteDriverInstReq)

		if err != nil {
			t.Fatal(err)
		}
		defer deleteDriverInstResp.Body.Close()

		Expect(deleteDriverInstResp.StatusCode).To(Equal(204))

		deleteDriverInstContent, err := ioutil.ReadAll(deleteDriverInstResp.Body)
		if err != nil {
			t.Fatal(err)
		}
		Expect(string(deleteDriverInstContent)).To(Equal(""))

		getDriverDeletedInstancesReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, firstDriverInstance.Id), nil)
		getDriverDeletedInstancesReq.Header.Add("Content-Type", "application/json")
		getDriverDeletedInstancesReq.Header.Add("Accept", "application/json")
		getDriverDeletedInstancesReq.Header.Add("Authorization", token)

		getDriverDeletedInstancesResp, err := http.DefaultClient.Do(getDriverDeletedInstancesReq)

		if err != nil {
			t.Fatal(err)
		}
		defer getDriverDeletedInstancesResp.Body.Close()

		Expect(getDriverDeletedInstancesResp.StatusCode).To((Equal(404)))
	}

	//delete driver
	deleteDriverReq, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driver.Id), nil)
	deleteDriverReq.Header.Add("Content-Type", "application/json")
	deleteDriverReq.Header.Add("Accept", "application/json")
	deleteDriverReq.Header.Add("Authorization", token)

	deleteDriverResp, err := http.DefaultClient.Do(deleteDriverReq)

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

	getDriverDeletedReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driver.Id), nil)
	getDriverDeletedReq.Header.Add("Content-Type", "application/json")
	getDriverDeletedReq.Header.Add("Accept", "application/json")
	getDriverDeletedReq.Header.Add("Authorization", token)

	getDriverDeletedResp, err := http.DefaultClient.Do(getDriverDeletedReq)

	if err != nil {
		t.Fatal(err)
	}
	defer getDriverDeletedResp.Body.Close()

	Expect(getDriverDeletedResp.StatusCode).To((Equal(404)))
}

func postgresEnvVarsExist() bool {
	return os.Getenv("POSTGRES_USER") != "" && os.Getenv("POSTGRES_PASSWORD") != "" && os.Getenv("POSTGRES_HOST") != "" &&
		os.Getenv("POSTGRES_PORT") != "" && os.Getenv("POSTGRES_DBNAME") != "" && os.Getenv("POSTGRES_SSLMODE") != ""
}

func setPostgresDriverInstanceValues(driverName, driverId string) []byte {
	values := []byte(fmt.Sprintf(`{"name":"%[1]s", "driver_id":"%[2]s", "configuration": {"host":"%[3]s","port":"%[4]s","user":"%[5]s","password":"%[6]s","dbname":"%[7]s","sslmode":"%[8]s"}}`,
		driverName,
		driverId,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_SSLMODE")))

	return values
}

func assertPostgresSchemaContains(schemaContent string) {
	Expect(schemaContent).To(ContainSubstring(`\"host\"`))
	Expect(schemaContent).To(ContainSubstring(`\"port\"`))
	Expect(schemaContent).To(ContainSubstring(`\"user\"`))
	Expect(schemaContent).To(ContainSubstring(`\"password\"`))
}

func mongoEnvVarsExist() bool {
	return os.Getenv("MONGO_USER") != "" && os.Getenv("MONGO_PASS") != "" && os.Getenv("MONGO_HOST") != "" && os.Getenv("MONGO_PORT") != ""
}

func setMongoDriverInstanceValues(driverName, driverId string) []byte {
	values := []byte(fmt.Sprintf(`{"name":"%[1]s", "driver_id":"%[2]s", "configuration": {"server":"%[3]s","port":"%[4]s","userid":"%[5]s","password":"%[6]s"}}`,
		driverName,
		driverId,
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASS")))

	return values
}

func assertMongoSchemaContains(schemaContent string) {
	Expect(schemaContent).To(ContainSubstring(`\"server\"`))
	Expect(schemaContent).To(ContainSubstring(`\"port\"`))
}

func mysqlEnvVarsExist() bool {
	return os.Getenv("MYSQL_USER") != "" && os.Getenv("MYSQL_PASS") != "" && os.Getenv("MYSQL_HOST") != "" && os.Getenv("MYSQL_PORT") != ""
}

func setMysqlDriverInstanceValues(driverName, driverId string) []byte {
	values := []byte(fmt.Sprintf(`{"name":"%[1]s", "driver_id":"%[2]s", "configuration": {"server":"%[3]s","port":"%[4]s","userid":"%[5]s","password":"%[6]s"}}`,
		driverName,
		driverId,
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASS")))

	return values
}

func assertMysqlSchemaContains(schemaContent string) {
	Expect(schemaContent).To(ContainSubstring(`\"server\"`))
	Expect(schemaContent).To(ContainSubstring(`\"port\"`))
	Expect(schemaContent).To(ContainSubstring(`\"userid\"`))
	Expect(schemaContent).To(ContainSubstring(`\"password\"`))
}

func dummyEnvVarsExist() bool {
	return true
}

func setDummyDriverInstanceValues(driverName, driverId string) []byte {
	values := []byte(fmt.Sprintf(`{"name":"%[1]s", "driver_id":"%[2]s", "configuration": {"succeed_count": "3"}}`,
		driverName,
		driverId))

	return values
}

func assertDummySchemaContains(schemaContent string) {
	Expect(schemaContent).To(ContainSubstring(`\"succeed_count\"`))
}

func getFileSha(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sha1 := sha1.New()
	_, err = io.Copy(sha1, f)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sha1.Sum(nil)), nil
}
