package mgmt

import (
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/hpcloud/cf-usb/lib/data"
	"github.com/hpcloud/cf-usb/lib/mgmt/authentication/uaa"
	"github.com/hpcloud/cf-usb/lib/operations"
	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	assert := assert.New(t)

	swaggerJSON, err := data.Asset("swagger-spec/api.json")
	if err != nil {
		t.Errorf("Error loading swagger data: %v", err)
	}

	swaggerSpec, err := spec.New(swaggerJSON, "")
	if err != nil {
		t.Errorf("Error loading swagger: %v", err)
	}
	mgmtAPI := operations.NewUsbMgmtAPI(swaggerSpec)

	auth, err := uaa.NewUaaAuth("", "", true)
	if err != nil {
		t.Errorf("Error instantiating uaa auth: %v", err)
	}

	ConfigureAPI(mgmtAPI, auth)

	params := operations.GetInfoParams{""}

	info, err := mgmtAPI.GetInfoHandler.Handle(params)
	if err != nil {
		t.Errorf("Error get info: %v", err)
	}

	assert.Equal("0.0.1", info.Version)
}
