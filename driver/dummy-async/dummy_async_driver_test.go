package dummyasync

import (
	"encoding/json"
	"testing"

	usbDriver "github.com/hpcloud/cf-usb/driver"
	"github.com/hpcloud/cf-usb/driver/status"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("postgres-provisioner")
var instanceID string = "00000000-0000-0000-0000-000000000001"
var settings = []byte(`{"succeed_count": "3"}`)
var config = (*json.RawMessage)(&settings)

func TestAsyncProvision(t *testing.T) {

	assert := assert.New(t)

	asyncDriver := NewDummyAsyncDriver(logger)
	response := usbDriver.Instance{}

	request := usbDriver.ProvisionInstanceRequest{}
	request.InstanceID = instanceID

	request.Config = config

	err := asyncDriver.ProvisionInstance(request, &response)
	assert.Nil(err)

	assert.Equal(status.InProgress, response.Status)

	getInstanceRequest := usbDriver.GetInstanceRequest{}
	getInstanceRequest.Config = config
	getInstanceRequest.InstanceID = instanceID

	err = asyncDriver.GetInstance(getInstanceRequest, &response)
	assert.Nil(err)

	//Two calls should return status in progress
	for i := 0; i < 2; i++ {
		err := asyncDriver.GetInstance(getInstanceRequest, &response)
		assert.Nil(err)
		assert.Equal(status.InProgress, response.Status)

	}

	//the third call must return status == created
	err = asyncDriver.GetInstance(getInstanceRequest, &response)
	assert.Nil(err)
	assert.Equal(status.Created, response.Status)

	deprovisionInstanceRequest := usbDriver.DeprovisionInstanceRequest{}
	deprovisionInstanceRequest.InstanceID = instanceID
	deprovisionInstanceRequest.Config = config
	err = asyncDriver.DeprovisionInstance(deprovisionInstanceRequest, &response)
	assert.Nil(err)
	assert.Equal(status.Deleted, response.Status)
}
