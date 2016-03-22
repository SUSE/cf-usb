package test_utils

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

type Configurator func(conf *config.Config)

type UsbRunner struct {
	Path    string
	TempDir string

	UsbBrokerPort      uint16
	JsonConfigDefaults string
	Configurator       Configurator

	Runner  *ginkgomon.Runner
	Process ifrit.Process

	ConfigFile string
}

func (r *UsbRunner) configure() *ginkgomon.Runner {
	cfgFile, err := ioutil.TempFile(r.TempDir, "usb-config.json")
	Expect(err).NotTo(HaveOccurred())

	var conf *config.Config

	err = json.Unmarshal([]byte(r.JsonConfigDefaults), &conf)
	Expect(err).NotTo(HaveOccurred())

	conf.BrokerAPI.Listen = r.BrokerAddress()

	r.Configurator(conf)

	modifiedConfigJson, err := json.Marshal(conf)
	Expect(err).NotTo(HaveOccurred())

	config := string(modifiedConfigJson)

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
	runner := r.configure()
	r.Process = ginkgomon.Invoke(runner)

	<-r.Process.Ready()

	return r.Process
}

func (r *UsbRunner) Stop() {
	ginkgomon.Interrupt(r.Process, 2)
}

func (r *UsbRunner) BrokerAddress() string {
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(int(r.UsbBrokerPort)))
}
