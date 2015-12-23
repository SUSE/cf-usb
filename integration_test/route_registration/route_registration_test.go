package route_registration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"

	. "github.com/hpcloud/cf-usb/integration_test/test_utils"
	"github.com/hpcloud/cf-usb/lib/config"

	"github.com/apcera/nats"
	"github.com/cloudfoundry/gunk/natsrunner"

	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/localip"
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
	"routes_register": {
        "nats_members": ["nats://nats:nats@127.0.0.1:4222"],
        "broker_api_host": "usb-broker.bosh-lite.com",
        "management_api_host": "usb-mgm.bosh-lite.com"
    },
    "drivers": {
    }
}
`

func getBinPath() string {
	dir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	return path.Join(dir, "../../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH), "usb")
}

func setDriverPathEnv() {
	dir, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	os.Setenv("USB_DRIVER_PATH", path.Join(dir, "../../build", fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)))
}

func initializeRunner() (*UsbRunner, *natsrunner.NATSRunner, *config.Config) {
	freePort, err := localip.LocalPort()
	Expect(err).NotTo(HaveOccurred())

	tempDir, err := ioutil.TempDir("", "cf-usb-test")
	Expect(err).NotTo(HaveOccurred())

	setDriverPathEnv()

	natsFreePort, err := localip.LocalPort()
	Expect(err).NotTo(HaveOccurred())
	natsRunner := natsrunner.NewNATSRunner(int(natsFreePort))
	natsRunner.Start()

	usbRunner := &UsbRunner{
		UsbBrokerPort:      freePort,
		Path:               getBinPath(),
		TempDir:            tempDir,
		JsonConfigDefaults: fileProviderConfig,
		Configurator: func(conf *config.Config) {
			conf.RoutesRegister.NatsMembers = []string{fmt.Sprintf("nats://127.0.0.1:%d", natsFreePort)}
		},
	}

	usbRunner.Start()

	provider := config.NewFileConfig(usbRunner.ConfigFile)
	configInfo, err := provider.LoadConfiguration()
	Expect(err).NotTo(HaveOccurred())

	return usbRunner, natsRunner, configInfo
}

func Test_BrokerApiPortIsRegisterd(t *testing.T) {
	RegisterTestingT(t)

	binpath := getBinPath()
	if _, err := os.Stat(binpath); os.IsNotExist(err) {
		t.Skip("Please build the solution before testing ", binpath)
		return
	}

	usb, natsRun, configInfo := initializeRunner()
	defer usb.Stop()
	defer natsRun.KillWithFire()

	Expect(configInfo.RoutesRegister).ToNot(BeNil())
	Expect(natsRun.MessageBus.Ping()).To(Equal(true))

	natsRun.MessageBus.Publish("router.start", []byte(`{"id":"some-router-id","minimumRegisterIntervalInSeconds":1}`))

	routerRegisterChannel := make(chan []byte)

	natsRun.MessageBus.Subscribe("router.register", func(msg *nats.Msg) {
		routerRegisterChannel <- msg.Data
	})

	localip, err := localip.LocalIP()
	Expect(err).NotTo(HaveOccurred())

	registerMsg := []byte{}
	Eventually(routerRegisterChannel, 3).Should(Receive(&registerMsg))
	Expect(registerMsg).To(ContainSubstring(configInfo.RoutesRegister.BrokerAPIHost))
	Expect(registerMsg).To(ContainSubstring(localip))
	Expect(registerMsg).To(ContainSubstring(strconv.Itoa(int(usb.UsbBrokerPort))))
}
