package config

import (
	"encoding/json"
	"github.com/frodenas/brokerapi"
	_ "github.com/golang/protobuf/proto" //workaround for godep + gomega

	"github.com/hashicorp/consul/api"
	"github.com/hpcloud/cf-usb/lib/config/consul"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

var IntegrationConfig = struct {
	Provider         ConfigProvider
	consulAddress    string
	consulDatacenter string
	consulUser       string
	consulPassword   string
	consulSchema     string
	consulToken      string
}{}

var DefaultConsulPath string = "consul"

func init() {
	IntegrationConfig.consulAddress = os.Getenv("CONSUL_ADDRESS")
	IntegrationConfig.consulDatacenter = os.Getenv("CONSUL_DATACENTER")
	IntegrationConfig.consulPassword = os.Getenv("CONSUL_PASSWORD")
	IntegrationConfig.consulUser = os.Getenv("CONSUL_USER")
	IntegrationConfig.consulSchema = os.Getenv("CONSUL_SCHEMA")
	IntegrationConfig.consulToken = os.Getenv("CONSUL_TOKEN")
}

func initProvider() (bool, error, ifrit.Process) {
	var consulConfig api.Config
	if IntegrationConfig.consulAddress == "" {
		return false, nil, nil
	}
	consulConfig.Address = IntegrationConfig.consulAddress
	consulConfig.Datacenter = IntegrationConfig.consulPassword

	var auth api.HttpBasicAuth
	auth.Username = IntegrationConfig.consulUser
	auth.Password = IntegrationConfig.consulPassword

	consulConfig.HttpAuth = &auth
	consulConfig.Scheme = IntegrationConfig.consulSchema

	consulConfig.Token = IntegrationConfig.consulToken
	getConsulReq, _ := http.NewRequest("GET", "http://localhost:8500", nil)
	getConsulResp, _ := http.DefaultClient.Do(getConsulReq)
	consulIsRunning := false
	if getConsulResp != nil && getConsulResp.StatusCode == 200 {
		consulIsRunning = true
	}

	var process ifrit.Process
	var err error
	if consulIsRunning == false {
		process, err = startConsulProcess()
		if err != nil {
			return false, err, nil
		}
	}

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		return false, err, nil
	}

	IntegrationConfig.Provider = NewConsulConfig(provisioner)
	return true, nil, process
}

func Test_IntDriverInstance(t *testing.T) {
	RegisterTestingT(t)

	initialized, err, process := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	var instance Instance
	instance.Name = "testInstance"
	raw := json.RawMessage("{\"a1\":\"b1\"}")
	instance.Configuration = &raw
	err = IntegrationConfig.Provider.SetInstance("testInstanceID", instance)
	assert.NoError(err)

	instanceInfo, _, err := IntegrationConfig.Provider.GetInstance("testInstanceID")

	assert.Equal("testInstance", instanceInfo.Name)
	assert.NoError(err)

	exist, err := IntegrationConfig.Provider.InstanceNameExists("testInstance")
	if err != nil {
		assert.Error(err, "Unable to check driver instance name existance")
	}
	assert.NoError(err)
	assert.True(exist)

	instanceDetails, err := IntegrationConfig.Provider.LoadDriverInstance("testInstanceID")
	t.Log("Load driver instance results:")
	t.Log(instanceDetails.Configuration)
	t.Log(instanceDetails.Dials)
	t.Log(instanceDetails.Service)
	assert.Equal("testInstance", instanceDetails.Name)
	assert.NoError(err)

	if process != nil {
		process.Signal(os.Kill)
		<-process.Wait()
	}
}

func Test_IntDial(t *testing.T) {
	RegisterTestingT(t)

	initialized, err, process := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	var instance Instance
	instance.Name = "testInstance"
	rawInstance := json.RawMessage("{\"a1\":\"b1\"}")
	instance.Configuration = &rawInstance
	err = IntegrationConfig.Provider.SetInstance("testInstanceID", instance)
	assert.NoError(err)

	var dialInfo Dial

	var plan brokerapi.ServicePlan
	plan.Description = "testPlan desc"
	plan.ID = "testPlanID"
	plan.Name = "free"
	plan.Metadata = &brokerapi.ServicePlanMetadata{DisplayName: "test plan"}

	raw := json.RawMessage("{\"a1\":\"b1\"}")

	dialInfo.Configuration = &raw
	dialInfo.Plan = plan

	err = IntegrationConfig.Provider.SetDial("testInstanceID", "testdialID", dialInfo)
	assert.NoError(err)

	dialDetails, instanceID, err := IntegrationConfig.Provider.GetDial("testdialID")
	t.Log(dialDetails)
	t.Log(instanceID)
	assert.NoError(err)

	if process != nil {
		process.Signal(os.Kill)
		<-process.Wait()
	}

}

func startConsulProcess() (ifrit.Process, error) {

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

	TempConsulPath, err := ioutil.TempDir(tmpConsul, "")
	if err != nil {
		return nil, err
	}

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
