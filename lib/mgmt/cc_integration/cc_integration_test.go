package ccintegration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hpcloud/cf-usb/lib/config"
	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("cc-integration")

func TestCreateServiceBroker(t *testing.T) {
	brokerName := "testbroker"
	assert := assert.New(t)

	config, err := loadConfig()
	assert.NoError(err)

	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("", nil)

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return(nil, nil)

	cci := NewCCIntegration(config, tokenGenerator, client, logger)
	assert.NotNil(cci)

	err = cci.CreateServiceBroker(brokerName)
	assert.NoError(err)
}

func loadConfig() (*config.Config, error) {
	workDir, err := os.Getwd()
	configFile := filepath.Join(workDir, "../../../test-assets/file-config/config.json")

	fileConfig := config.NewFileConfig(configFile)

	conf, err := fileConfig.LoadConfiguration()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
