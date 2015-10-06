package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
	"github.com/stretchr/testify/assert"
)

var testService = brokerapi.Service{
	ID:          "GIUD",
	Name:        "testService",
	Description: "",
	Bindable:    true,
	Plans:       []brokerapi.ServicePlan{},
	Metadata:    brokerapi.ServiceMetadata{},
	Tags:        []string{},
}

var testDriverConfig = config.DriverConfig{
	DriverType: "dummy",
	ServiceIDs: []string{"GUID"},
}

var testConfig = config.Config{
	Crednetials:    brokerapi.BrokerCredentials{},
	ServiceCatalog: []brokerapi.Service{testService},
	DriverConfigs:  []config.DriverConfig{testDriverConfig},
	Listen:         ":5580",
	APIVersion:     "2.6",
	LogLevel:       "debug",
}

func setupEnv() (*UsbBroker, error) {
	buildDir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "hpcloud", "cf-usb", "build",
		fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	os.Setenv("USB_DRIVER_PATH", buildDir)

	configProperties := `{  
            "property_one":"one",
            "property_two":"two"
         }`

	data, err := json.Marshal(configProperties)
	if err != nil {
		return nil, err
	}

	*testDriverConfig.Configuration = json.RawMessage(data)

	driverProvider, err := NewDriverProvider("dummy", config.DriverProperties{})
	if err != nil {
		return nil, err
	}

	broker := NewUsbBroker([]*DriverProvider{driverProvider}, &testConfig, lager.NewLogger("brokerTests"))
	return broker, nil
}

func TestGetCatalog(t *testing.T) {
	assert := assert.New(t)
	broker, err := setupEnv()
	assert.True(true)
}
