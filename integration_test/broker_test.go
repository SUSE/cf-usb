package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/cf-usb/lib/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/golang/protobuf/proto" //workaround for godep + gomega
	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-golang/lager/lagertest"
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

var ConsulConfig = struct {
	consulAddress    string
	consulDatacenter string
	consulUser       string
	consulPassword   string
	consulSchema     string
	consulToken      string
}{}

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

func start_consulProcess() (ifrit.Process, error) {

	defaultConsulPath := "/usr/local/bin/consul"

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
	cmd := exec.Command(binpath, "consulConfigProvider", "-a", ConsulConfig.consulAddress)

	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	defer cmd.Process.Kill()

	consulClient, err := init_consulProvisioner()
	if err != nil {
		t.Fatal(err)
	}

	provider := config.NewConsulConfig(consulClient)

	configInfo, err := provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password
	//wait for process to start
	time.Sleep(2 * time.Second)

	resp, err := http.Get(fmt.Sprintf("http://%s:%s@%s/v2/catalog", user, pass, ConsulConfig.consulAddress))

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

	var uaaFakeServer *ghttp.Server
	uaaFakeServer = ghttp.NewServer()

	var ccFakeServer *ghttp.Server
	ccFakeServer = ghttp.NewServer()

	uaaFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281214}`),
		),
	)

	uaaFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281216}`),
		),
	)

	uaaFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/oauth/token"),
			ghttp.RespondWith(200, `{"access_token":"replace-me", "expires_in": 3404281218}`),
		),
	)

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

	ccFakeServer.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/v2/service_brokers"),
			func(http.ResponseWriter, *http.Request) {
				time.Sleep(0 * time.Second)
			},
			ghttp.RespondWith(200, `{"resources":[{"metadata":{"guid":""}}]}`),
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

	serviceGuid := "83E94C97-C755-46A5-8653-461517EB442M"

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

	servicePlanGuid := "98E94C97-C755-46A5-8653-461517EB442M"

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

	var consulClient consul.ConsulProvisionerInterface

	if ConsulConfig.consulAddress == "" {
		ConsulConfig.consulAddress = "127.0.0.1:8500"
		ConsulConfig.consulSchema = "http"

		consulProcess, err := start_consulProcess()
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			consulProcess.Signal(os.Kill)
			<-consulProcess.Wait()
		}()

		t.Log("consul started")
	}

	consulClient, err = init_consulProvisioner()
	if err != nil {
		t.Fatal(err)
	}

	var list api.KVPairs

	list = append(list, &api.KVPair{Key: "usb/loglevel", Value: []byte("debug")})
	list = append(list, &api.KVPair{Key: "usb/broker_api", Value: []byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}")})
	list = append(list, &api.KVPair{Key: "usb/management_api", Value: []byte(fmt.Sprintf("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"%[1]s\"}},\"cloud_controller\":{\"api\":\"%[2]s\",\"skip_tsl_validation\":true}}", uaaPublicKey, ccFakeServer.URL()))})

	err = consulClient.PutKVs(&list, nil)
	if err != nil {
		t.Fatal(err)
	}

	provider := config.NewConsulConfig(consulClient)

	_, err = provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binpath, "consulConfigProvider", "-a", ConsulConfig.consulAddress)

	stde, err := cmd.StderrPipe()
	stdo, err := cmd.StdoutPipe()
	go io.Copy(os.Stdout, stde)
	go io.Copy(os.Stdout, stdo)

	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	defer cmd.Process.Kill()
	//wait for usb to start
	time.Sleep(2 * time.Second)

	token, err := GenerateUaaToken()
	if err != nil {
		t.Fatal(err)
	}

	if postgresEnvVarsExist() {
		driverName := "postgres"
		newDriverReq, err := http.NewRequest("POST", "http://localhost:54053/drivers", strings.NewReader(fmt.Sprintf(`{"name":"%[1]s", "driver_type":"%[2]s"}`, driverName, driverName)))
		newDriverReq.Header.Add("Content-Type", "application/json")
		newDriverReq.Header.Add("Accept", "application/json")
		newDriverReq.Header.Add("Authorization", token)

		newDriverResp, err := http.DefaultClient.Do(newDriverReq)

		if err != nil {
			t.Fatal(err)
		}
		defer newDriverResp.Body.Close()

		driverContent, err := ioutil.ReadAll(newDriverResp.Body)
		if err != nil {
			t.Fatal(err)
		}

		type DriverResponse struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		}

		var driver DriverResponse

		err = json.Unmarshal(driverContent, &driver)
		if err != nil {
			fmt.Println("error:", err)
		}

		Expect(driver.Id).ToNot(BeNil())
		Expect(driver.Name).To(Equal(driverName))

		getDriverReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:54053/drivers/%[1]s", driver.Id), nil)
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
		Expect(getDriverContent).To(ContainSubstring(driver.Id))

		config := []byte(fmt.Sprintf(`{"name":"%[1]s", "driver_id":"%[2]s", "configuration": {"host":"%[3]s","port":"%[4]s","user":"%[5]s","password":"%[6]s","dbname":"%[7]s","sslmode":"%[8]s"}}`,
			driverName,
			driver.Id,
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DBNAME"),
			os.Getenv("POSTGRES_SSLMODE")))
		newDriverInstReq, err := http.NewRequest("POST", "http://localhost:54053/driver_instances", bytes.NewBuffer(config))
		newDriverInstReq.Header.Add("Content-Type", "application/json")
		newDriverInstReq.Header.Add("Accept", "application/json")
		newDriverInstReq.Header.Add("Authorization", token)

		newDriverInstResp, err := http.DefaultClient.Do(newDriverInstReq)

		if err != nil {
			t.Fatal(err)
		}
		defer newDriverInstResp.Body.Close()

		driverInstContent, err := ioutil.ReadAll(newDriverInstResp.Body)
		if err != nil {
			t.Fatal(err)
		}

		type DriverInstanceResponse struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Service string `json:"service"`
		}

		var driverInstance DriverInstanceResponse

		err = json.Unmarshal(driverInstContent, &driverInstance)
		if err != nil {
			fmt.Println("error:", err)
		}
		t.Logf("driver instance: %s", string(driverInstContent))
		Expect(driverInstContent).To(ContainSubstring(driver.Name))

		getPlanReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:54053/plans?q=driverInstanceId:%[1]s", driver.Id), nil)
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

		getServiceReq, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:54053/services/%[1]s", driverInstance.Service), nil)
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
}

func postgresEnvVarsExist() bool {
	return os.Getenv("POSTGRES_USER") != "" && os.Getenv("POSTGRES_PASSWORD") != "" && os.Getenv("POSTGRES_HOST") != "" &&
		os.Getenv("POSTGRES_PORT") != "" && os.Getenv("POSTGRES_DBNAME") != "" && os.Getenv("POSTGRES_SSLMODE") != ""
}
