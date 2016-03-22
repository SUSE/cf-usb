package consul

import (
	"github.com/hashicorp/consul/api"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

var consulConfig = struct {
	Address         string
	Datacenter      string
	User            string
	Pass            string
	Scheme          string
	Token           string
	TestProvisioner ConsulProvisionerInterface
}{}

var DefaultConsulPath string = "consul"

func init() {
	consulConfig.Address = os.Getenv("CONSUL_ADDRESS")
	consulConfig.Datacenter = os.Getenv("CONSUL_DATACENTER")
	consulConfig.User = os.Getenv("CONSUL_USERNAME")
	consulConfig.Pass = os.Getenv("CONSUL_PASSWORD")
	consulConfig.Scheme = os.Getenv("CONSUL_SCHEME")
	consulConfig.Token = os.Getenv("CONSUL_TOKEN")
}

func initProvider() (error, ifrit.Process) {
	var config api.Config
	var err error

	config.Address = consulConfig.Address
	config.Datacenter = consulConfig.Datacenter
	config.Scheme = consulConfig.Scheme
	config.Token = consulConfig.Token

	if consulConfig.User != "" {
		var auth api.HttpBasicAuth
		auth.Username = consulConfig.User
		auth.Password = consulConfig.Pass
		config.HttpAuth = &auth
	}

	consulConfig.TestProvisioner, err = New(&config)

	getConsulReq, _ := http.NewRequest("GET", "http://localhost:8500", nil)
	getConsulResp, _ := http.DefaultClient.Do(getConsulReq)
	consulIsRunning := false
	if getConsulResp != nil && getConsulResp.StatusCode == 200 {
		consulIsRunning = true
	}

	var process ifrit.Process
	if consulIsRunning == false {
		process, err = startConsulProcess()
		if err != nil {
			return err, nil
		}
	}

	return nil, process
}

func Test_KeyValue(t *testing.T) {
	RegisterTestingT(t)
	assert := assert.New(t)

	if consulConfig.Address == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS, CONSUL_DATACENTER, (CONSUL_USERNAME, CONSUL_PASSWORD) / CONSUL_TOKEN")
	}

	err, process := initProvider()
	assert.NoError(err)

	log.Println("Testing add key-value")
	err = consulConfig.TestProvisioner.AddKV("testKey", []byte("testValue"), nil)

	assert.NoError(err)

	log.Println("Testing get key-value")
	val, err := consulConfig.TestProvisioner.GetValue("testKey")

	assert.NoError(err)

	log.Println("Testing delete key-value")
	var options api.WriteOptions
	err = consulConfig.TestProvisioner.DeleteKV("testKey", &options)
	assert.NoError(err)

	if val != nil {
		log.Println("Test successful, retrieved value: ", string(val))
	}

	if process != nil {
		process.Signal(os.Kill)
		<-process.Wait()
	}

}

func Test_KVList(t *testing.T) {
	RegisterTestingT(t)
	assert := assert.New(t)

	var list api.KVPairs

	if consulConfig.Address == "" {
		t.Skip("Skipping test as Consul env vars are not set: CONSUL_ADDRESS, CONSUL_DATACENTER, (CONSUL_USERNAME, CONSUL_PASSWORD) / CONSUL_TOKEN")
	}

	err, process := initProvider()
	assert.NoError(err)

	log.Println("Testing put key-value list")

	list = append(list, &api.KVPair{Key: "usb/api_version", Value: []byte("2.0")})

	list = append(list, &api.KVPair{Key: "usb/broker_api", Value: []byte("{\"listen\":\":54054\",\"credentials\":{\"username\":\"demouser\",\"password\":\"demopassword\"}}")})

	list = append(list, &api.KVPair{Key: "usb/management_api", Value: []byte("{\"listen\":\":54053\",\"uaa_secret\":\"myuaasecret\",\"uaa_client\":\"myuaaclient\",\"authentication\":{\"uaa\":{\"adminscope\":\"usb.management.admin\",\"public_key\":\"\"}}}")})

	list = append(list, &api.KVPair{Key: "usb/drivers/mysql2", Value: []byte("mysql2")})

	list = append(list, &api.KVPair{Key: "usb/drivers/mysql2/instances/00000000-0000-0000-0000-0000000000M1/Name", Value: []byte("mysql-local")})

	list = append(list, &api.KVPair{Key: "usb/drivers/mysql2/instances/00000000-0000-0000-0000-0000000000M1/Configuration", Value: []byte("{\"server\":\"127.0.0.1\",\"port\":\"3306\",\"userid\":\"root\",\"password\":\"password\"}")})

	list = append(list, &api.KVPair{Key: "usb/drivers/mysql2/instances/00000000-0000-0000-0000-0000000000M1/dials/B0000000-0000-0000-0000-000000000LM1", Value: []byte("{\"id\":\"B0000000-0000-0000-0000-000000000LM1\",\"configuration\":{},\"plan\":{\"name\":\"free\",\"id\":\"53425178-F731-49E7-9E53-5CF4BE9D807L\",\"description\":\"This is the first plan\",\"free\":true}}")})

	list = append(list, &api.KVPair{Key: "usb/drivers/mysql2/instances/00000000-0000-0000-0000-0000000000M1/service", Value: []byte("{\"id\":\"83E94C97-C755-46A5-8653-461517EB442L\",\"bindable\":true,\"name\":\"MysqlLocalService\",\"description\":\"Mysql Local Service\",\"tags\":[\"mysql\"],\"metadata\":{\"providerDisplayName\":\"Mysql Local Service\"}}")})

	err = consulConfig.TestProvisioner.PutKVs(&list, nil)
	assert.NoError(err)

	log.Println("Testing get key-value list")
	val, err := consulConfig.TestProvisioner.GetAllKVs("usb", nil)

	assert.NoError(err)

	if val != nil {
		for _, kv := range val {
			log.Println("Retrieved value: ", kv.Key, string(kv.Value))
		}
	}

	log.Println("Testing get keys")
	result, err := consulConfig.TestProvisioner.GetAllKeys("usb/", "", nil)
	log.Println("Retrieved keys:", result)
	assert.NoError(err)

	log.Println("Testing delete key-value list with prefix")
	err = consulConfig.TestProvisioner.DeleteKVs("usb/drivers/", nil)
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
