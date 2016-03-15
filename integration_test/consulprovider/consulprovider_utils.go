package consulprovider

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/golang/protobuf/proto" //workaround for godep + gomega
	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/pivotal-golang/localip"
	"github.com/satori/go.uuid"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var LoggerSB *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")
var DefaultConsulPath string = "consul"
var TempConsulPath string
var DefaultBuildDir string = "../../../build"
var TempDriversPath string
var ConsulProcess ifrit.Process
var BrokerApiPort uint16
var ManagementApiPort uint16

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

var UaaPublicKey = `-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHFr+KICms+tuT1OXJwhCUmR2d\nKVy7psa8xzElSyzqx7oJyfJ1JZyOzToj9T5SfTIq396agbHJWVfYphNahvZ/7uMX\nqHxf+ZH9BL1gk9Y6kCnbM5R60gfwjyW1/dQPjOzn9N394zd2FJoFHwdq9Qs0wBug\nspULZVNRxq7veq/fzwIDAQAB\n-----END PUBLIC KEY-----`

var ConsulConfig = struct {
	ConsulAddress    string
	ConsulDatacenter string
	ConsulUser       string
	ConsulPassword   string
	ConsulSchema     string
	ConsulToken      string
}{}

// list of drivers to be tested

var Drivers = []struct {
	DriverType                     string
	EnvVarsExistFunc               func() bool
	SetDriverInstanceValuesFunc    func(driverName, driverId string) []byte
	AssertDriverSchemaContainsFunc func(schemaContent string)
}{
	{"dummy-async", dummyEnvVarsExist, setDummyDriverInstanceValues, assertDummySchemaContains},
	{"postgres", postgresEnvVarsExist, setPostgresDriverInstanceValues, assertPostgresSchemaContains},
	{"mongo", mongoEnvVarsExist, setMongoDriverInstanceValues, assertMongoSchemaContains},
	{"mysql", mysqlEnvVarsExist, setMysqlDriverInstanceValues, assertMysqlSchemaContains},
}

// models http responses

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

type PlanResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Id          string `json:"id"`
	DialId      string `json:"dial_id"`
	Free        bool   `json:"free"`
}

type DialResponse struct {
	Configuration    interface{} `json:"configuration,omitempty"`
	Id               string      `json:"id"`
	DriverInstanceId string      `json:"driver_instance_id"`
	Plan             string      `json:"plan"`
}

type ServiceResponse struct {
	Bindable         bool        `json:"bindable"`
	Description      string      `json:"description"`
	DriverInstanceID string      `json:"driver_instance_id"`
	Id               string      `json:"id"`
	Metadata         interface{} `json:"metadata"`
	Name             string      `json:"name"`
	Tags             []string    `json:"tags"`
}

// init test functions

func CheckSolutionIsBuild() (string, bool, error) {
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	if err != nil {
		return "", false, err
	}

	binpath := path.Join(dir, DefaultBuildDir, architecture, "usb")

	_, err = os.Stat(binpath)
	return binpath, os.IsNotExist(err), err
}

func SetupConsulForFirstRun() error {
	var err error

	BrokerApiPort, err = localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())
	ManagementApiPort, err = localip.LocalPort()
	Expect(err).ToNot(HaveOccurred())

	tmpConsul := path.Join(os.TempDir(), "consul")

	if _, err := os.Stat(tmpConsul); err == nil {
		err := os.RemoveAll(tmpConsul)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(tmpConsul, 0755)
	if err != nil {
		return err
	}

	TempConsulPath, err = ioutil.TempDir(tmpConsul, "")
	if err != nil {
		return err
	}

	TempDriversPath = path.Join(os.TempDir(), "drivers")

	err = os.MkdirAll(TempDriversPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func RunConsulProcess(brokerApiPort, managementApiPort uint16, ccServerUrl, driversPath string) (consul.ConsulProvisionerInterface, error) {
	var consulClient consul.ConsulProvisionerInterface

	getConsulReq, _ := http.NewRequest("GET", "http://localhost:8500", nil)
	getConsulResp, _ := http.DefaultClient.Do(getConsulReq)
	consulIsRunning := false
	if getConsulResp != nil && getConsulResp.StatusCode == 200 {
		consulIsRunning = true
	}

	if (strings.Contains(ConsulConfig.ConsulAddress, "127.0.0.1") || strings.Contains(ConsulConfig.ConsulAddress, "localhost")) && !consulIsRunning {
		ConsulConfig.ConsulAddress = "127.0.0.1:8500"
		ConsulConfig.ConsulSchema = "http"

		var err error
		ConsulProcess, err = startConsulProcess()
		if err != nil {
			return nil, err
		}
	}

	consulClient, err := initConsulProvisioner(driversPath)
	if err != nil {
		return nil, err
	}

	var list api.KVPairs

	list = append(list, &api.KVPair{Key: "usb/api_version", Value: []byte("2.0")})
	list = append(list, &api.KVPair{Key: "usb/broker_api", Value: []byte(fmt.Sprintf("{\"listen\":\":%v\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}", brokerApiPort))})
	list = append(list, &api.KVPair{Key: "usb/management_api", Value: []byte(fmt.Sprintf("{\"listen\":\":%v\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"%v\"}},\"cloud_controller\":{\"api\":\"%s\",\"skip_tls_validation\":true}}", ManagementApiPort, UaaPublicKey, ccServerUrl))})

	err = consulClient.PutKVs(&list, nil)
	if err != nil {
		return consulClient, err
	}

	return consulClient, nil
}

func RunUsbProcess(binPath, consulAddress string) (ifrit.Process, error) {
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

func SetFakeServers() (*ghttp.Server, *ghttp.Server) {
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

// external packages functions

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

func GetFileSha(filePath string) (string, error) {
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

func ExecuteHttpCall(verb, path string, body io.Reader) (*http.Response, error) {
	token, err := GenerateUaaToken()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(verb, path, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// unmarshal http responses

func UnmarshalDriverResponse(body io.ReadCloser) (string, DriverResponse, error) {
	var driver DriverResponse

	driverContent, err := ioutil.ReadAll(body)
	if err != nil {
		return "", driver, err
	}

	err = json.Unmarshal(driverContent, &driver)
	if err != nil {
		return "", driver, err
	}

	return string(driverContent), driver, nil
}

func UnmarshalDriversResponse(body io.ReadCloser) (string, []DriverResponse, error) {
	var drivers []DriverResponse

	driversContent, err := ioutil.ReadAll(body)
	if err != nil {
		return "", drivers, err
	}

	err = json.Unmarshal(driversContent, &drivers)
	if err != nil {
		return "", drivers, err
	}

	return string(driversContent), drivers, nil
}

func UnmarshalDriverInstanceResponse(body io.ReadCloser) (string, DriverInstanceResponse, error) {
	var driverInstance DriverInstanceResponse

	driverInstContent, err := ioutil.ReadAll(body)
	if err != nil {
		return "", driverInstance, err
	}

	err = json.Unmarshal(driverInstContent, &driverInstance)
	if err != nil {
		return "", driverInstance, err
	}

	return string(driverInstContent), driverInstance, nil
}

func UnmarshalDriverInstancesResponse(body io.ReadCloser) (string, []DriverInstanceResponse, error) {
	var driverInstances []DriverInstanceResponse

	driverInstancesContent, err := ioutil.ReadAll(body)
	if err != nil {
		return "", driverInstances, err
	}

	err = json.Unmarshal(driverInstancesContent, &driverInstances)
	if err != nil {
		return "", driverInstances, err
	}

	return string(driverInstancesContent), driverInstances, nil
}

// operations shared between tests

func GetDrivers(managementApiPort uint16) (int, string, []DriverResponse, error) {
	var drivers []DriverResponse

	getDriversResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/drivers", managementApiPort), nil)
	if err != nil {
		return 0, "", drivers, err
	}
	defer getDriversResp.Body.Close()

	driversContent, drivers, err := UnmarshalDriversResponse(getDriversResp.Body)
	if err != nil {
		return 0, "", drivers, err
	}

	return getDriversResp.StatusCode, driversContent, drivers, nil
}

func CreateDriver(managementApiPort uint16, driverName string) (int, string, DriverResponse, error) {
	var driver DriverResponse

	newDriverResp, err := ExecuteHttpCall("POST", fmt.Sprintf("http://localhost:%[1]v/drivers", managementApiPort), strings.NewReader(fmt.Sprintf(`{"name":"%[1]s", "driver_type":"%[2]s"}`, driverName, driverName)))
	if err != nil {
		return 0, "", driver, err
	}
	defer newDriverResp.Body.Close()

	driverContent, driver, err := UnmarshalDriverResponse(newDriverResp.Body)
	if err != nil {
		return 0, "", driver, err
	}

	return newDriverResp.StatusCode, driverContent, driver, nil
}

func UploadDriver(managementApiPort uint16, driverType, driverId string) (int, string, error) {
	token, err := GenerateUaaToken()
	if err != nil {
		return 0, "", err
	}

	// upload driver bits test
	dir, err := os.Getwd()
	if err != nil {
		return 0, "", err
	}

	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	bitsPath := path.Join(dir, DefaultBuildDir, architecture, driverType)
	sha, err := GetFileSha(bitsPath)
	if err != nil {
		return 0, "", err
	}

	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	err = body_writer.WriteField("sha", sha)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("Error writing to buffer: %s", err.Error()))
	}

	// use the body_writer to write the Part headers to the buffer
	_, err = body_writer.CreateFormFile("file", bitsPath)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("Error writing to buffer: %s", err.Error()))
	}

	// the file data will be the second part of the body
	fh, err := os.Open(bitsPath)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("Error opening file: %s", err.Error()))
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	_ = fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	_, err = fh.Stat()
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("Error stating file: %s. %s", bitsPath, err.Error()))
	}

	uploadDriverBitsReq, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s/bits", managementApiPort, driverId), request_reader)
	uploadDriverBitsReq.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	uploadDriverBitsReq.Header.Add("Accept", "application/json")
	uploadDriverBitsReq.Header.Add("Authorization", token)

	uploadDriverBitsResp, err := http.DefaultClient.Do(uploadDriverBitsReq)
	if err != nil {
		return 0, "", err
	}
	defer uploadDriverBitsResp.Body.Close()

	uploadDriverBitsContent, err := ioutil.ReadAll(uploadDriverBitsResp.Body)
	if err != nil {
		return 0, "", err
	}

	return uploadDriverBitsResp.StatusCode, string(uploadDriverBitsContent), nil
}

func DeleteDriver(managementApiPort uint16, driverId string) (int, string, error) {
	deleteDriverResp, err := ExecuteHttpCall("DELETE", fmt.Sprintf("http://localhost:%[1]v/drivers/%[2]s", managementApiPort, driverId), nil)
	if err != nil {
		return 0, "", err
	}
	defer deleteDriverResp.Body.Close()

	deleteDriverContent, err := ioutil.ReadAll(deleteDriverResp.Body)
	if err != nil {
		return 0, "", err
	}

	return deleteDriverResp.StatusCode, string(deleteDriverContent), nil
}

func GetDriverInstances(managementApiPort uint16, driverId string) (int, string, []DriverInstanceResponse, error) {
	var driverInstances []DriverInstanceResponse

	getDriverInstancesResp, err := ExecuteHttpCall("GET", fmt.Sprintf("http://localhost:%[1]v/driver_instances?driver_id=%[2]s", managementApiPort, driverId), nil)
	if err != nil {
		return 0, "", driverInstances, err
	}
	defer getDriverInstancesResp.Body.Close()

	driverInstancesContent, driverInstances, err := UnmarshalDriverInstancesResponse(getDriverInstancesResp.Body)
	if err != nil {
		return 0, "", driverInstances, err
	}

	return getDriverInstancesResp.StatusCode, driverInstancesContent, driverInstances, nil
}

func CreateDriverInstance(managementApiPort uint16, driver DriverResponse, driverInstanceValues func(driverName, driverId string) []byte) (int, string, DriverInstanceResponse, error) {
	var driverInstance DriverInstanceResponse

	instanceValues := driverInstanceValues(driver.Name, driver.Id)

	newDriverInstResp, err := ExecuteHttpCall("POST", fmt.Sprintf("http://localhost:%[1]v/driver_instances", managementApiPort), bytes.NewBuffer(instanceValues))
	if err != nil {
		return 0, "", driverInstance, err
	}
	defer newDriverInstResp.Body.Close()

	driverInstContent, driverInstance, err := UnmarshalDriverInstanceResponse(newDriverInstResp.Body)
	if err != nil {
		return 0, "", driverInstance, err
	}

	return newDriverInstResp.StatusCode, driverInstContent, driverInstance, nil
}

func DeleteDriverInstance(managementApiPort uint16, driverInstanceId string) (int, string, error) {
	deleteDriverInstResp, err := ExecuteHttpCall("DELETE", fmt.Sprintf("http://localhost:%[1]v/driver_instances/%[2]s", managementApiPort, driverInstanceId), nil)
	if err != nil {
		return 0, "", err
	}
	defer deleteDriverInstResp.Body.Close()

	deleteDriverInstContent, err := ioutil.ReadAll(deleteDriverInstResp.Body)
	if err != nil {
		return 0, "", err
	}

	return deleteDriverInstResp.StatusCode, string(deleteDriverInstContent), nil
}

// cc fake responses functions

func SetupCcHttpFakeResponsesCreateDriverInstance(uaaFakeServer, ccFakeServer *ghttp.Server) {
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

func SetupCcHttpFakeResponsesDeleteDriverInstance(notExistLabel string, uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	serviceGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services", fmt.Sprintf("q=label:%s", notExistLabel)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				`{"resources":[]}`),
		),
	)

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
			ghttp.VerifyRequest("GET", fmt.Sprintf("/v2/service_plans/%[1]s/service_instances", servicePlanGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{}`),
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

	brokerGuid := uuid.NewV4().String()

	ccFakeServer.RouteToHandler("GET", "/v2/service_brokers",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			ghttp.RespondWith(200, fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, brokerGuid)),
		),
	)

	ccFakeServer.RouteToHandler("PUT", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{}`),
		),
	)

	ccFakeServer.RouteToHandler("DELETE", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("DELETE", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(204, `{}`),
		),
	)
}

func SetupCcHttpFakeResponsesCreateDial(notExistLabel string, uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	serviceGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services", fmt.Sprintf("q=label:%s", notExistLabel)),
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

func SetupCcHttpFakeResponsesUpdatePlan(uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	brokerGuid := uuid.NewV4().String()

	ccFakeServer.RouteToHandler("GET", "/v2/service_brokers",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			ghttp.RespondWith(200, fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, brokerGuid)),
		),
	)

	ccFakeServer.RouteToHandler("PUT", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid)),
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

func SetupCcHttpFakeResponsesUpdateService(notExistLabel string, uaaFakeServer, ccFakeServer *ghttp.Server) {
	uaaFakeServer.RouteToHandler("POST", "/oauth/token",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	serviceGuid := uuid.NewV4().String()

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services", fmt.Sprintf("q=label:%s", notExistLabel)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				`{"resources":[]}`),
		),
	)

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/services"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200,
				fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"},"entity":{"name":"%[2]s"}}]}`, serviceGuid, notExistLabel)),
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

	brokerGuid := uuid.NewV4().String()

	ccFakeServer.RouteToHandler("GET", "/v2/service_brokers",
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			ghttp.RespondWith(200, fmt.Sprintf(`{"resources":[{"metadata":{"guid":"%[1]s"}}]}`, brokerGuid)),
		),
	)

	ccFakeServer.RouteToHandler("PUT", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("PUT", fmt.Sprintf("/v2/service_brokers/%[1]s", brokerGuid)),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{}`),
		),
	)
}

// functions for running consul

func initConsulProvisioner(driversPath string) (consul.ConsulProvisionerInterface, error) {
	var consulConfig api.Config
	consulConfig.Address = ConsulConfig.ConsulAddress
	consulConfig.Datacenter = ConsulConfig.ConsulPassword

	if driversPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		buildDir := filepath.Join(workDir, DefaultBuildDir, fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
		os.Setenv("USB_DRIVER_PATH", buildDir)
	} else {
		os.Setenv("USB_DRIVER_PATH", driversPath)
	}

	var auth api.HttpBasicAuth
	auth.Username = ConsulConfig.ConsulUser
	auth.Password = ConsulConfig.ConsulPassword

	consulConfig.HttpAuth = &auth
	consulConfig.Scheme = ConsulConfig.ConsulSchema

	consulConfig.Token = ConsulConfig.ConsulToken

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		return nil, err
	}
	return provisioner, nil
}

func startConsulProcess() (ifrit.Process, error) {

	consulRunner := ginkgomon.New(ginkgomon.Config{
		Name:              "consul",
		Command:           exec.Command(DefaultConsulPath, "agent", "-server", "-bootstrap-expect", "1", "-data-dir", TempConsulPath, "-advertise", "127.0.0.1"),
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

// postgres driver specific functions

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

// mongo driver specific functions

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

// mysql driver specific functions

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

// dummy-async driver specific functions

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
