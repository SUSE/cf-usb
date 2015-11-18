package integration_test

import (
	"fmt"
	"github.com/hpcloud/cf-usb/lib/config"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/hpcloud/cf-usb/lib/config/consul"

	"github.com/hashicorp/consul/api"
)

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

func init_consulProvisioner() (*consul.ConsulProvisionerInterface, error) {
	var consulConfig api.Config
	consulConfig.Address = ConsulConfig.consulAddress
	consulConfig.Datacenter = ConsulConfig.consulPassword

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
	return &provisioner, nil
}

func Test_BrokerWithFileConfigProviderCatalog(t *testing.T) {

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

	var testFile string

	testFile = os.Getenv("TEST_CONFIG_FILE")
	if len(testFile) == 0 {
		testFile = path.Join(dir, "../test-assets", "file-config", "config.json")

		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Skip("Test configuration not found ", testFile)
			return
		}
	}

	cmd := exec.Command(binpath, "fileConfigProvider", "-p", testFile)

	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	defer cmd.Process.Kill()

	provider := config.NewFileConfig(testFile)

	configInfo, err := provider.LoadConfiguration()
	if err != nil {
		t.Fatal(err)
	}

	user := configInfo.BrokerAPI.Credentials.Username
	pass := configInfo.BrokerAPI.Credentials.Password
	//wait for process to start
	time.Sleep(2 * time.Second)

	resp, err := http.Get(fmt.Sprintf("http://%s:%s@localhost:54054/v2/catalog", user, pass))
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

	provider := config.NewConsulConfig(*consulClient)

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
