package config

import (
	"encoding/json"

	_ "github.com/golang/protobuf/proto" //workaround for godep + gomega
	"github.com/SUSE/cf-usb/lib/brokermodel"

	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/SUSE/cf-usb/lib/config/consul"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var IntegrationConfig = struct {
	Provider         Provider
	consulAddress    string
	consulDatacenter string
	consulUser       string
	consulPassword   string
	consulSchema     string
	consulToken      string
}{}

var DefaultConsulPath = "consul"

func init() {
	IntegrationConfig.consulAddress = os.Getenv("CONSUL_ADDRESS")
	IntegrationConfig.consulDatacenter = os.Getenv("CONSUL_DATACENTER")
	IntegrationConfig.consulPassword = os.Getenv("CONSUL_PASSWORD")
	IntegrationConfig.consulUser = os.Getenv("CONSUL_USER")
	IntegrationConfig.consulSchema = os.Getenv("CONSUL_SCHEMA")
	IntegrationConfig.consulToken = os.Getenv("CONSUL_TOKEN")
}

func initProvider() (bool, ifrit.Process, error) {
	var consulConfig api.Config
	if IntegrationConfig.consulAddress == "" {
		return false, nil, nil
	}
	consulConfig.Address = IntegrationConfig.consulAddress

	if consulConfig.Address == "" {
		return false, nil, nil
	}

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
			return false, nil, err
		}
	}

	provisioner, err := consul.New(&consulConfig)
	if err != nil {
		return false, nil, err
	}

	IntegrationConfig.Provider = NewConsulConfig(provisioner)
	return true, process, nil
}

func Test_IntDriverInstance(t *testing.T) {
	RegisterTestingT(t)

	initialized, process, err := initProvider()
	if initialized == false || err != nil {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), seting CONSUL_ADDRESS to 127.0.0.1:8500 will start a service on the local machine")
		t.Log(err)
	}

	assert := assert.New(t)

	var instance Instance
	instance.Name = "testInstance"
	instance.CaCert = "-----BEGIN CERTIFICATE-----\nMIIFCzCCAvOgAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwppbnRl\ncm5hbENBMB4XDTE2MDcxODIxMjYxOVoXDTI2MDcxODIxMjYyMlowFTETMBEGA1UE\nAxMKaW50ZXJuYWxDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAMPy\n1fazELUA0mi9ffXpCGamjZsvtv/agOastPN51LTJcdX7uPQWw0XgCVjJ1qdUfuy1\n6DG56YS7CTRgVZ66Rf3TQN0wRa7iMfz/Ngkfwy8qYi8VgiOms55aZwn+RfJAioac\n5GsLMUsXq7So6R5CW2RvEo47xTOYAoM3FhmniMsu7KJq/Iq9aldgtqrVkbicqzgq\n4Cqs4f5TlpzqvlT1pzgV55FbX5arqj9vldRXS32s0xDWvgAt7tdgV16tPYu6qyWa\nsFYBGQE1Fjl0edmR9WNCopZaZUHdTjJ6+hVQsEOyQ28JypO0dwtuzv/d/DwduhiU\nS1KApX4M03ChKKusJXw/Wc8GDmaWt82gd2LhptpEgL//4fzt39yCTGdfOZL/irjD\nl28jBnNYSdz019OsWqE3uPOKgEdimQNZnPGIGNu1R/yyb4BbvX8ZiT0fj4unjGd/\n77F7oVf4z8T4lUAgGArcsbIc3SgKBBi/uQwzJEOo7cg+p9PcE7e6PBRIBRQYeYgG\nPsBGMTJO+4BJrzcUgz3r7Fginn3cMhWofB0XYOP73u0h3hyJDe+GPedWiO0SFeZ2\nNZREF131Fods0vV5m1VcOTYDZxaAj9b1DvBMtJvwgR5J98aVx8OrihoNOcEGGHON\nHqsx/ik7cDShrHf2ZWKtdY1egoKwfFKKarhi552fAgMBAAGjZjBkMA4GA1UdDwEB\n/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBS4BaKkSHho/32t\nlFUCDKQlM/CPfDAfBgNVHSMEGDAWgBS4BaKkSHho/32tlFUCDKQlM/CPfDANBgkq\nhkiG9w0BAQsFAAOCAgEAi+OcyYCIA5cYoituIUXmQfAhh7qjOrcZNsPZq6KjrE5W\nhMwyEHrfWB0Rv2+WSIa4GvvMgkweGeY4s+3Y/3ml/AciYzolqlUUzAWudqi3gW9G\nMKY89boFgzcEAFeX0MJq9SSQuSlNmpHy4glPEKV1abevMAuiD3CRq7XoDdN68K4L\nULMfdrE78ae+Kd4xdCtJyccNGZ1PAPEIMfGIY+I1VbBre9M0z0YJi+O9EuBilxYb\ns4tuN5LVEaLk8AuC+Ik+nH+zWTKfLiSAe1r5Gvx5ERweYkvylYV5u7ERb/B6PXrN\nRwpmTkLGMr4bMNWGEpw/Z9mxHFxiJgNgZLg4IPkEYr2dtpO8s3FalT2DvOClsVIJ\njsSNZONumhdn2GVhOJoctTwh2kXSXDGsuNBhsAg6Q3F2ejNJCqKfqyW4PLnxi8/2\nrtYBKDRpjjSjX1MOkBluAuoxR6DjybldK2+cag9dPvAVYidxCH8j7y7eawBEBsFs\nIsuu+bom/Wv8zNDQpXXEUBEuJPRqsrLZdHCJH2e7hoKS03oXRXNkERRjekQLaUnc\n51Zrq5GowsmMSircZJPfufIltTurBK5uwX+iaxalb/ynqAzD2s8BxI5dzrRi2Jkv\nTk0STTyj0QVppy4kwJfI3vuAPTg45gNkhNDGAQ3pRnh3g4Jq9jgptcgahpkGoGs=\n-----END CERTIFICATE-----"

	err = IntegrationConfig.Provider.SetInstance("testInstanceID", instance)
	assert.NoError(err)

	instanceInfo, _, err := IntegrationConfig.Provider.GetInstance("testInstanceID")

	assert.Equal("testInstance", instanceInfo.Name)
	//Test if default is false
	assert.Equal(false, instance.SkipSsl)
	assert.Equal("-----BEGIN CERTIFICATE-----\nMIIFCzCCAvOgAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwppbnRl\ncm5hbENBMB4XDTE2MDcxODIxMjYxOVoXDTI2MDcxODIxMjYyMlowFTETMBEGA1UE\nAxMKaW50ZXJuYWxDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAMPy\n1fazELUA0mi9ffXpCGamjZsvtv/agOastPN51LTJcdX7uPQWw0XgCVjJ1qdUfuy1\n6DG56YS7CTRgVZ66Rf3TQN0wRa7iMfz/Ngkfwy8qYi8VgiOms55aZwn+RfJAioac\n5GsLMUsXq7So6R5CW2RvEo47xTOYAoM3FhmniMsu7KJq/Iq9aldgtqrVkbicqzgq\n4Cqs4f5TlpzqvlT1pzgV55FbX5arqj9vldRXS32s0xDWvgAt7tdgV16tPYu6qyWa\nsFYBGQE1Fjl0edmR9WNCopZaZUHdTjJ6+hVQsEOyQ28JypO0dwtuzv/d/DwduhiU\nS1KApX4M03ChKKusJXw/Wc8GDmaWt82gd2LhptpEgL//4fzt39yCTGdfOZL/irjD\nl28jBnNYSdz019OsWqE3uPOKgEdimQNZnPGIGNu1R/yyb4BbvX8ZiT0fj4unjGd/\n77F7oVf4z8T4lUAgGArcsbIc3SgKBBi/uQwzJEOo7cg+p9PcE7e6PBRIBRQYeYgG\nPsBGMTJO+4BJrzcUgz3r7Fginn3cMhWofB0XYOP73u0h3hyJDe+GPedWiO0SFeZ2\nNZREF131Fods0vV5m1VcOTYDZxaAj9b1DvBMtJvwgR5J98aVx8OrihoNOcEGGHON\nHqsx/ik7cDShrHf2ZWKtdY1egoKwfFKKarhi552fAgMBAAGjZjBkMA4GA1UdDwEB\n/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBS4BaKkSHho/32t\nlFUCDKQlM/CPfDAfBgNVHSMEGDAWgBS4BaKkSHho/32tlFUCDKQlM/CPfDANBgkq\nhkiG9w0BAQsFAAOCAgEAi+OcyYCIA5cYoituIUXmQfAhh7qjOrcZNsPZq6KjrE5W\nhMwyEHrfWB0Rv2+WSIa4GvvMgkweGeY4s+3Y/3ml/AciYzolqlUUzAWudqi3gW9G\nMKY89boFgzcEAFeX0MJq9SSQuSlNmpHy4glPEKV1abevMAuiD3CRq7XoDdN68K4L\nULMfdrE78ae+Kd4xdCtJyccNGZ1PAPEIMfGIY+I1VbBre9M0z0YJi+O9EuBilxYb\ns4tuN5LVEaLk8AuC+Ik+nH+zWTKfLiSAe1r5Gvx5ERweYkvylYV5u7ERb/B6PXrN\nRwpmTkLGMr4bMNWGEpw/Z9mxHFxiJgNgZLg4IPkEYr2dtpO8s3FalT2DvOClsVIJ\njsSNZONumhdn2GVhOJoctTwh2kXSXDGsuNBhsAg6Q3F2ejNJCqKfqyW4PLnxi8/2\nrtYBKDRpjjSjX1MOkBluAuoxR6DjybldK2+cag9dPvAVYidxCH8j7y7eawBEBsFs\nIsuu+bom/Wv8zNDQpXXEUBEuJPRqsrLZdHCJH2e7hoKS03oXRXNkERRjekQLaUnc\n51Zrq5GowsmMSircZJPfufIltTurBK5uwX+iaxalb/ynqAzD2s8BxI5dzrRi2Jkv\nTk0STTyj0QVppy4kwJfI3vuAPTg45gNkhNDGAQ3pRnh3g4Jq9jgptcgahpkGoGs=\n-----END CERTIFICATE-----", instanceInfo.CaCert)

	assert.NoError(err)

	exist, err := IntegrationConfig.Provider.InstanceNameExists("testInstance")
	if err != nil {
		assert.Error(err, "Unable to check driver instance name existance")
	}
	assert.NoError(err)
	assert.True(exist)

	instanceDetails, err := IntegrationConfig.Provider.LoadDriverInstance("testInstanceID")
	t.Log("Load driver instance results:")
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

	initialized, process, err := initProvider()
	if initialized == false {
		t.Skip("Skipping Consul Set Driver test, environment variables not set: CONSUL_ADDRESS(host:port), CONSUL_DATACENTER, CONSUL_TOKEN / CONSUL_USER + CONSUL_PASSWORD, CONSUL_SCHEMA")
		t.Log(err)
	}

	assert := assert.New(t)

	var instance Instance
	instance.Name = "testInstance"
	err = IntegrationConfig.Provider.SetInstance("testInstanceID", instance)
	assert.NoError(err)

	var dialInfo Dial

	var plan brokermodel.Plan
	plan.Description = "testPlan desc"
	plan.ID = "testPlanID"
	plan.Name = "free"
	plan.Metadata = &brokermodel.PlanMetadata{Metadata: struct{ DisplayName string }{"test plan"}}

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
