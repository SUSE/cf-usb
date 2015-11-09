package ccapi

import (
	"encoding/json"
	"testing"

	"github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/mocks"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var loggerSP *lagertest.TestLogger = lagertest.NewTestLogger("cc-api")

func TestUpdateServicePlanVisibility(t *testing.T) {
	tokenGenerator := new(mocks.GetTokenInterface)
	tokenGenerator.On("GetToken").Return("bearer atoken", nil)

	var pra []PlanResource
	pr := PlanResource{Values: PlanMetadata{}, Entity: PlanValues{}}
	pra = append(pra, pr)

	getPlanMocked := PlanResources{Resources: pra}
	values, err := json.Marshal(getPlanMocked)
	if err != nil {
		t.Errorf("Error marshall get plan: %v", err)
	}

	client := new(mocks.HttpClient)
	client.Mock.On("Request", mock.Anything).Return(values, nil)
	client.Mock.On("Request", mock.Anything).Return(nil, nil)

	sp := NewServicePlan(client, tokenGenerator, "ccApi", loggerSP)

	err = sp.Update("a-service-guid")
	if err != nil {
		t.Errorf("Error enable service access: %v", err)
	}

	assert.NoError(t, err)
}
