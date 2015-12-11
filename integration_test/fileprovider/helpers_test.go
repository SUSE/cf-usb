package fileprovider_test

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os/exec"
	"strconv"

	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	"github.com/hpcloud/cf-usb/lib/config"
	. "github.com/onsi/gomega"
)

var fileProviderConfig = `
{
    "api_version": "2.6",
    "logLevel": "debug",
    "broker_api": {
		"external_url": "http://127.0.0.1:54054",
        "listen": ":54054",		
        "credentials": {
            "username": "username",
            "password": "password"
        }
    },
    "drivers": {
          "f533244b-e270-4ee0-aa9e-4733d52f097a":{
            "driver_type": "dummy",
            "driver_instances": {
                "dd4fcc63-b28f-4795-b398-36d9fc75efe7":{
                    "name": "dummy1",
                    "configuration": {
                        "property_one": "one",
                        "property_two": "two"
                    },
                    "dials": {
                          "881d876b-c933-4d9e-87c1-6d4b238abc0b": {
                            "configuration": {
                                "max_dbsize_mb": 2
                            },
                            "plan": {
                                "name": "free",
                                "id": "1a7cc5ee-4a46-4af4-9af5-b6f2b6050ee9",
                                "description": "free plan",
                                "free": true
                            }
                        }
                    },
                    "service": {
                        "id": "de8464a4-1d05-4f25-8a74-9790448d13cd",
                        "bindable": true,
                        "name": "dummy-test",
                        "description": "Dummy test service",
                        "tags": [
                            "dummy"
                        ],
                        "metadata": {
                            "providerDisplayName": "Dummy"
                        }
                    }
                }
            }
        }
    }
}
`

func setupConfigFile(sourceConfig string, brokerListen string) (string, error) {
	var conf *config.Config

	err := json.Unmarshal([]byte(sourceConfig), &conf)
	if err != nil {
		return "", err
	}

	conf.BrokerAPI.Listen = brokerListen

	modifiedConfigJson, err := json.Marshal(conf)
	if err != nil {
		return "", err
	}

	return string(modifiedConfigJson), nil
}

type UsbRunner struct {
	Path    string
	TempDir string

	UsbBrokerPort uint16

	Runner  *ginkgomon.Runner
	Process ifrit.Process

	ConfigFile string
}

func (r *UsbRunner) Configure() *ginkgomon.Runner {
	cfgFile, err := ioutil.TempFile(r.TempDir, "usb-config.json")
	Expect(err).NotTo(HaveOccurred())

	config, err := setupConfigFile(fileProviderConfig, r.BrokerAddress())
	Expect(err).NotTo(HaveOccurred())

	_, err = cfgFile.WriteString(config)
	Expect(err).NotTo(HaveOccurred())
	cfgFile.Close()

	r.ConfigFile = cfgFile.Name()

	r.Runner = ginkgomon.New(ginkgomon.Config{
		Name:       "cf-usb",
		StartCheck: "usb.start-listening-brokerapi",
		Command: exec.Command(
			r.Path,
			"fileConfigProvider", "-p", cfgFile.Name(),
		),
	})

	return r.Runner
}

func (r *UsbRunner) Start() ifrit.Process {
	runner := r.Configure()
	r.Process = ginkgomon.Invoke(runner)
	return r.Process
}

func (r *UsbRunner) Stop() {
	ginkgomon.Interrupt(r.Process, 2)
}

func (r *UsbRunner) BrokerAddress() string {
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(int(r.UsbBrokerPort)))
}
