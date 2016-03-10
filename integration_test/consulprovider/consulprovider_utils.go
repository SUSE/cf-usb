package consulprovider

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var LoggerSB *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")
var DefaultConsulPath string = "consul"
var TempConsulPath string
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

var ConsulConfig = struct {
	ConsulAddress    string
	ConsulDatacenter string
	ConsulUser       string
	ConsulPassword   string
	ConsulSchema     string
	ConsulToken      string
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

func init_consulProvisioner(driversPath string) (consul.ConsulProvisionerInterface, error) {
	var consulConfig api.Config
	consulConfig.Address = ConsulConfig.ConsulAddress
	consulConfig.Datacenter = ConsulConfig.ConsulPassword

	if driversPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		buildDir := filepath.Join(workDir, "../../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
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

func Start_usbProcess(binPath, consulAddress string) (ifrit.Process, error) {
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

func Check_solutionBuild() (string, bool, error) {
	architecture := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	dir, err := os.Getwd()
	if err != nil {
		return "", false, err
	}

	binpath := path.Join(dir, "../../../build", architecture, "usb")

	_, err = os.Stat(binpath)
	return binpath, os.IsNotExist(err), err
}

func Run_consul(brokerApiPort, managementApiPort uint16, ccServerUrl, driversPath string) (consul.ConsulProvisionerInterface, error) {
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
		ConsulProcess, err = start_consulProcess()
		if err != nil {
			return nil, err
		}
	}

	consulClient, err := init_consulProvisioner(driversPath)
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

func Set_fakeServers() (*ghttp.Server, *ghttp.Server) {
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

func Setup_firstConsulRun() error {
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

	return nil
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
