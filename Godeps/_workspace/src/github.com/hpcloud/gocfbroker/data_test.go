package gocfbroker

import (
	"encoding/json"
	"testing"
)

func TestProvisionRequest_Validate(t *testing.T) {
	t.Parallel()

	happy := testProvisionRequest
	happy.ServiceID = testConfig.Services[0].ID
	happy.PlanID = testConfig.Services[0].Plans[0].ID
	errs := happy.validate(testConfig.Catalog)
	if errs != nil {
		t.Error("There was a problem with the provision request:", validationErrors(errs))
	}

	unhappy := happy
	unhappy.ServiceID = "junk"
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "service_id does not reference a service from the catalog" {
		t.Error("Wrong message:", errs.Errors[0])
	}

	unhappy = happy
	unhappy.PlanID = "junk"
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "plan_id does not reference a services plan from the catalog" {
		t.Error("Wrong message:", errs.Errors[0])
	}

	unhappy = happy
	unhappy.OrganizationGUID = ""
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "organization_guid must not be blank" {
		t.Error("Wrong message:", errs.Errors[0])
	}

	unhappy = happy
	unhappy.SpaceGUID = ""
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "space_guid must not be blank" {
		t.Error("Wrong message:", errs.Errors[0])
	}
}

func TestProvisionRequest_Equal(t *testing.T) {
	t.Parallel()

	var req1, req2 ProvisionRequest

	if !req1.Equal(req2) {
		t.Errorf("Should be equal: \n%#v\n%#v", req1, req2)
	}

	req1.ServiceID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.ServiceID = ""
	req1.PlanID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.PlanID = ""
	req1.OrganizationGUID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.OrganizationGUID = ""
	req1.SpaceGUID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.SpaceGUID = ""
	req1.Parameters = MakeJSONRawMessage("a")
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}
}

func TestUpdateProvisionRequest_Validate(t *testing.T) {
	t.Parallel()

	happy := testUpdateProvisionRequest
	happy.ServiceID = testConfig.Services[0].ID
	happy.PlanID = testConfig.Services[0].Plans[0].ID
	errs := happy.validate(testConfig.Catalog)
	if errs != nil {
		t.Error("Should validate successfully")
	}

	unhappy := happy
	unhappy.PlanID = "junk"
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "plan_id does not reference a services plan from the catalog" {
		t.Error("Wrong message:", errs.Errors[0])
	}

	unhappy = happy
	unhappy.ServiceID = "junk"
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "service_id does not reference a service from the catalog" {
		t.Error("Wrong message:", errs.Errors[0])
	}

	unhappy = happy
	unhappy.ServiceID = testConfig.Services[1].ID // Non-updateable
	unhappy.PlanID = testConfig.Services[1].Plans[0].ID
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "service plan is not updatable" {
		t.Error("Wrong message:", errs.Errors[0])
	}
}

func TestUpdateProvisionRequest_Equal(t *testing.T) {
	t.Parallel()

	var req1, req2 UpdateProvisionRequest

	if !req1.Equal(req2) {
		t.Errorf("Should be equal: \n%#v\n%#v", req1, req2)
	}

	req1.PlanID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.PlanID = ""
	req1.Parameters = MakeJSONRawMessage("a")
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}
}

func TestBindingRequest_Validate(t *testing.T) {
	t.Parallel()

	happy := testBindingRequest
	happy.ServiceID = testConfig.Services[0].ID
	happy.PlanID = testConfig.Services[0].Plans[0].ID
	errs := happy.validate(testConfig.Catalog)
	if errs != nil {
		t.Error("There was a problem with the update request:", validationErrors(errs))
	}

	unhappy := happy
	unhappy.ServiceID = "junk"
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "service_id does not reference a service from the catalog" {
		t.Error("Wrong message:", errs.Errors[0])
	}

	unhappy = happy
	unhappy.PlanID = "junk"
	errs = unhappy.validate(testConfig.Catalog)
	if errs.Errors[0] != "plan_id does not reference a services plan from the catalog" {
		t.Error("Wrong message:", errs.Errors[0])
	}
}

func TestBindingRequest_Equal(t *testing.T) {
	t.Parallel()

	var req1, req2 BindingRequest

	if !req1.Equal(req2) {
		t.Errorf("Should be equal: \n%#v\n%#v", req1, req2)
	}

	req1.ServiceID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.ServiceID = ""
	req1.PlanID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.PlanID = ""
	req1.AppGUID = "a"
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}

	req1.AppGUID = ""
	req1.Parameters = MakeJSONRawMessage("a")
	if req1.Equal(req2) {
		t.Errorf("Should not be equal: \n%#v\n%#v", req1, req2)
	}
}

func TestJsonRawEquals(t *testing.T) {
	tests := []struct {
		A     json.RawMessage
		B     json.RawMessage
		Equal bool
	}{
		{json.RawMessage(`a`), json.RawMessage(`a`), true},
		{nil, nil, true},

		{json.RawMessage(`a`), json.RawMessage(`b`), false},
		{nil, json.RawMessage(`b`), false},
		{json.RawMessage(`a`), nil, false},
	}

	for i, test := range tests {
		if eq := jsonRawEqual(&test.A, &test.B); eq != test.Equal {
			t.Errorf("%d) %s and %s should eq %v, got: %v", i, test.A, test.B, test.Equal, eq)
		}
	}
}

func testMustEncode(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	return string(b)
}

var testProvisionRequest = ProvisionRequest{
	ServiceID:        testConfig.Services[0].ID,
	PlanID:           testConfig.Services[0].Plans[0].ID,
	OrganizationGUID: "organizationguid",
	SpaceGUID:        "spaceguid",
}
var testProvisionRequestJSON = testMustEncode(testProvisionRequest)

var testProvisionResponse = ProvisionResponse{
	DashboardURL: "url",
}
var testProvisionResponseJSON = testMustEncode(testProvisionResponse)

var testBindingRequest = BindingRequest{
	ServiceID: testConfig.Services[0].ID,
	PlanID:    testConfig.Services[0].Plans[0].ID,
	AppGUID:   "appguid",
}
var testBindingRequestJSON = testMustEncode(testBindingRequest)

var testBindingResponse = BindingResponse{
	Credentials:    MakeJSONRawMessage(`{"anything":["more","things"]}`),
	SyslogDrainURL: "syslogdrainurl",
}

var testUpdateProvisionRequest = UpdateProvisionRequest{
	ServiceID: "serviceid",
	PlanID:    "planid",
	PreviousValues: &PreviousValues{
		ServiceID:      "serviceid",
		PlanID:         "planid",
		OrganizationID: "organizationid",
		SpaceID:        "spaceid",
	},
}
var testUpdateProvisionRequestJSON = testMustEncode(testUpdateProvisionRequest)
