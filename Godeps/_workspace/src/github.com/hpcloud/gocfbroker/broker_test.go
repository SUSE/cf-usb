package gocfbroker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func init() {
	logOutput = ioutil.Discard
}

func setupTestBroker() (*testService, testStore, *Broker) {
	service := &testService{}
	storer := newTestStore()

	broker, err := New(service, storer, *testConfig)
	if err != nil {
		panic("Could not create broker: " + err.Error())
	}

	return service, storer, broker
}

func setupHTTPReq(method, path string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	r, _ := http.NewRequest(method, path, body)
	r.SetBasicAuth(testConfig.AuthUser, testConfig.AuthPassword)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	return w, r
}

func verifyJobResult(b *Broker, instanceID string, errMsg string, status jobStatus) error {
	if res := b.jobResult(instanceID); res.ErrorMsg != errMsg {
		return fmt.Errorf("Want error: %q, got %q", errMsg, res.ErrorMsg)
	} else if res.Status != status {
		return fmt.Errorf("Want status: %q, got %q", status, res.Status)
	}
	return nil
}

func TestBroker_New(t *testing.T) {
	t.Parallel()

	broker, err := New(&testService{}, newTestStore(), *testConfig)
	if err != nil {
		t.Fatal(err)
	}

	if broker.provisioner == nil {
		t.Error("Need a provisioner")
	}
	if broker.db.Storer == nil {
		t.Error("Need db")
	}
	if len(broker.apiVersion.String()) == 0 {
		t.Error("Need version")
	}
}

func TestBroker_LoadConfig(t *testing.T) {
	var opts Options
	if err := LoadConfigReader(bytes.NewReader(testConfigJSON), &opts); err != nil {
		t.Error(err)
	}

	b, err := json.Marshal(opts)
	if err != nil {
		t.Error(err)
	}
	if strings.TrimSpace(string(b)) != strings.TrimSpace(string(testConfigJSON)) {
		t.Error("The JSON config did not serialize back to its original form:", string(b))
	}
}

func TestBroker_BuildAPI(t *testing.T) {
	t.Parallel()

	_, _, broker := setupTestBroker()

	router := broker.buildRoutes()
	api := broker.buildAPI(router)

	// Inject a route for testing into the router
	called := false
	httpr := router.(*httprouter.Router)
	httpr.GET("/testpath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		called = true
	})

	r, _ := http.NewRequest("GET", "/testpath", nil)

	// Test that clients that don't accept JSON fail
	w := httptest.NewRecorder()
	r.Header.Set("Accept", "text/html")
	api.ServeHTTP(w, r)
	if w.Code != http.StatusNotAcceptable {
		t.Error("Want status not acceptable from json mw:", w.Code)
	} else if called {
		t.Error("Should have stopped the call")
	}

	// Check that they must be authed
	w = httptest.NewRecorder()
	r.Header.Set("Accept", "*/*")
	api.ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Error("Want unauthorized from basic auth mw:", w.Code)
	} else if called {
		t.Error("Should have stopped the call")
	}

	// Check that they must request a sane version
	w = httptest.NewRecorder()
	r.SetBasicAuth(testConfig.AuthUser, testConfig.AuthPassword)
	r.Header.Set(headerBroker, "5.6")
	api.ServeHTTP(w, r)
	if w.Code != http.StatusPreconditionFailed {
		t.Error("Want precodindition failed from version mw:", w.Code)
	} else if called {
		t.Error("Should have stopped the call")
	}

	// Check that it calls the handler in the end.
	w = httptest.NewRecorder()
	r.Header.Set(headerBroker, "2.4")
	api.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Error("Status OK", w.Code)
	} else if !called {
		t.Error("It should have called the handler")
	}
}

func TestBroker_BuildRoutes(t *testing.T) {
	t.Parallel()

	_, _, broker := setupTestBroker()

	router := broker.buildRoutes()

	w, r := setupHTTPReq("GET", "/v2/catalog", nil)
	router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Error("Wrong code:", w.Code)
	}

	w, r = setupHTTPReq("PUT", "/v2/service_instances/a?accepts_incomplete=true", strings.NewReader("a"))
	router.ServeHTTP(w, r)
	// Error out from JSON decode failure
	if w.Code != http.StatusBadRequest {
		t.Error("Wrong code:", w.Code)
	}

	w, r = setupHTTPReq("DELETE", "/v2/service_instances/a?accepts_incomplete=true", nil)
	go func() {
		<-broker.jobs
		broker.jobAccepts <- nil
	}()
	router.ServeHTTP(w, r)
	// No instance exists
	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	w, r = setupHTTPReq("PATCH", "/v2/service_instances/a?accepts_incomplete=true", strings.NewReader("a"))
	go func() {
		<-broker.jobs
		broker.jobAccepts <- nil
	}()
	router.ServeHTTP(w, r)
	// Error out from JSON decode failure
	if w.Code != http.StatusBadRequest {
		t.Error("Wrong code:", w.Code)
	}

	w, r = setupHTTPReq("PUT", "/v2/service_instances/a/service_bindings/b", strings.NewReader("a"))
	router.ServeHTTP(w, r)
	// Error out from JSON decode failure
	if w.Code != http.StatusBadRequest {
		t.Error("Wrong code:", w.Code)
	}

	w, r = setupHTTPReq("DELETE", "/v2/service_instances/a/service_bindings/b", nil)
	router.ServeHTTP(w, r)
	// Error out from missing binding
	if w.Code != http.StatusGone {
		t.Error("Wrong code:", w.Code)
	}

	w, r = setupHTTPReq("GET", "/v2/service_instances/a/last_operation", nil)
	go func() {
		req := <-broker.jobResultReqs
		req.resultCh <- jobResult{ErrorMsg: errMsgNoJobExists}
	}()
	router.ServeHTTP(w, r)
	// Error out from missing binding
	if w.Code != http.StatusNotFound {
		t.Error("Wrong code:", w.Code)
	}

	w, r = setupHTTPReq("GET", "/v2/not_a_route", nil)
	router.ServeHTTP(w, r)
	// Error out from missing route
	if w.Code != http.StatusNotFound {
		t.Error("Wrong code:", w.Code)
	}
	if s := w.Body.String(); strings.TrimSpace(s) != `{"description":"api endpoint does not exist"}` {
		t.Error("Wrong body:", s)
	}
}

func TestBroker_GetCatalog(t *testing.T) {
	t.Parallel()

	_, _, broker := setupTestBroker()

	w, r := setupHTTPReq("GET", "/", nil)
	if err := broker.fetchCatalog(w, r, nil); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusOK {
		t.Error("Wrong code:", w.Code)
	}

	var cat Catalog
	dec := json.NewDecoder(w.Body)
	if err := dec.Decode(&cat); err != nil {
		t.Error(err)
	}

	if len(cat.Services) != len(testConfig.Services) {
		t.Error("Wrong # of services:", len(cat.Services))
	}
}

func TestBroker_ProvisionInstance(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	instanceID := "instance-id"

	w, r := setupHTTPReq("GET", "/", strings.NewReader(testProvisionRequestJSON))
	p := httprouter.Params{{"instance_id", instanceID}}

	service, opener, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	err := verifyJobResult(broker, instanceID, "", jobStatusSucceeded)
	if err != nil {
		t.Error(err)
	}

	if _, ok := opener.keyValues["instance-id"]; !ok {
		t.Error("Instance should be in the DB.")
	}

	if !service.provisioned {
		t.Error("Should have called Provision")
	}
}

func TestBroker_ProvisionInstanceErrors(t *testing.T) {
	t.Parallel()

	w, r := setupHTTPReq("GET", "/", strings.NewReader(":D"))
	var p httprouter.Params

	_, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	// Fail on instance_id missing
	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusBadRequest {
		t.Error("Wrong code:", w.Code)
	}

	// Fail on bad json
	p = append(p, httprouter.Param{Key: "instance_id", Value: "instance-id"})
	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusBadRequest {
		t.Error("Wrong code:", w.Code)
	}
}

func TestBroker_ProvisionInstanceDupSame(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	instanceID := "instance-id"
	w, r := setupHTTPReq("GET", "/", strings.NewReader(testProvisionRequestJSON))
	p := httprouter.Params{{"instance_id", instanceID}}

	service, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	// Provision first time
	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	if !service.provisioned {
		t.Error("Should have called provision")
	}

	err := verifyJobResult(broker, instanceID, "", jobStatusSucceeded)
	if err != nil {
		t.Error(err)
	}

	// Provision second time, same value
	service.provisioned = false
	w, r = setupHTTPReq("GET", "/", strings.NewReader(testProvisionRequestJSON))
	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	err = verifyJobResult(broker, instanceID, errMsgNoDiff, jobStatusSucceeded)
	if err != nil {
		t.Error(err)
	}

	if service.provisioned {
		t.Error("It should not have called provision")
	}
}

func TestBroker_ProvisionInstanceDupDiff(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	instanceID := "instance-id"
	w, r := setupHTTPReq("GET", "/", strings.NewReader(testProvisionRequestJSON))
	p := httprouter.Params{{"instance_id", instanceID}}

	service, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	// Provision first time
	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	if !service.provisioned {
		t.Error("Should have called provision")
	}

	err := verifyJobResult(broker, instanceID, "", jobStatusSucceeded)
	if err != nil {
		t.Error(err)
	}

	// Provision second time, diff values
	newProvisionReq := testProvisionRequest
	newProvisionReq.PlanID = testConfig.Services[0].Plans[1].ID
	service.provisioned = false

	data, err := json.Marshal(newProvisionReq)
	if err != nil {
		t.Error(err)
	}

	body := ioutil.NopCloser(bytes.NewReader(data))
	w, r = setupHTTPReq("GET", "/", body)
	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	err = verifyJobResult(broker, instanceID, errMsgConflict, jobStatusFailed)
	if err != nil {
		t.Error(err)
	}

	if service.provisioned {
		t.Error("Should not have called provision")
	}
}

func TestBroker_ProvisionInstanceMissingFields(t *testing.T) {
	t.Parallel()

	service, _, broker := setupTestBroker()

	doTest := func(key string, testData ProvisionRequest, description string) {
		data, err := json.Marshal(testData)
		if err != nil {
			t.Errorf("%s) Failed json marshal: %v %#v", key, err, testData)
		}
		w, r := setupHTTPReq("GET", "/", bytes.NewReader(data))
		p := httprouter.Params{{"instance_id", "instance-id"}}
		if err = broker.provisionInstance(w, r, p); err != nil {
			t.Errorf("%s) %v", key, err)
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("%s) Wrong code: %v", key, w.Code)
		}
		if str := w.Body.String(); !strings.Contains(str, description) {
			t.Errorf("Want: %q, Got: %q", description, str)
		}

		if service.provisioned {
			t.Error("Should not have called provision")
		}
	}

	testData := testProvisionRequest
	testData.ServiceID = ""
	doTest("ServiceID", testData, "service_id does not reference a service")

	testData = testProvisionRequest
	testData.PlanID = ""
	doTest("PlanID", testData, "plan_id does not reference a services plan")

	testData = testProvisionRequest
	testData.OrganizationGUID = ""
	doTest("OrganizationGUID", testData, "organization_guid must not be blank")

	testData = testProvisionRequest
	testData.SpaceGUID = ""
	doTest("SpaceGUID", testData, "space_guid must not be blank")
}

func TestBroker_ProvisionInstanceJobExists(t *testing.T) {
	t.Parallel()

	instanceID := "instance-id"
	service, _, broker := setupTestBroker()

	// Setup job
	var b brokerJob
	b.InstanceID = instanceID
	b.JobResult.Status = jobStatusInProgress
	if err := putJob(broker.db, &b); err != nil {
		t.Error(err)
	}

	w, r := setupHTTPReq("GET", "/", strings.NewReader(testProvisionRequestJSON))
	p := httprouter.Params{{"instance_id", instanceID}}

	go broker.dispatcher()
	defer close(broker.jobs)

	if err := broker.provisionInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != statusTooManyRequests {
		t.Error("Wrong code:", w.Code)
	}

	if service.provisioned {
		t.Error("Should not have called provision")
	}
}

func TestBroker_DeprovisionInstance(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	w, r := setupHTTPReq("DELETE", "/?plan_id=planid&service_id=serviceid", nil)
	instanceID := "instance-id"
	p := httprouter.Params{{"instance_id", instanceID}}

	service, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}
	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	err := broker.deprovisionInstance(w, r, p)
	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	if err := verifyJobResult(broker, instanceID, "", jobStatusSucceeded); err != nil {
		t.Error(err)
	}

	if !service.deprovisioned {
		t.Error("Should have called deprovision")
	}
}

func TestBroker_DeprovisionMissing(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	w, r := setupHTTPReq("DELETE", "/?plan_id=planid&service_id=serviceid", nil)
	instanceID := "instance-id"
	p := httprouter.Params{{"instance_id", instanceID}}

	service, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	if err := broker.deprovisionInstance(w, r, p); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	if err := verifyJobResult(broker, instanceID, errMsgGone, jobStatusFailed); err != nil {
		t.Error(err)
	}

	if service.deprovisioned {
		t.Error("Should not have called deprovision")
	}
}

func TestBroker_UpdateInstance(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	updateReq := testUpdateProvisionRequest
	updateReq.ServiceID = testConfig.Services[0].ID
	updateReq.PlanID = testConfig.Services[0].Plans[1].ID
	data, err := json.Marshal(updateReq)
	if err != nil {
		t.Error(err)
	}

	w, r := setupHTTPReq("PATCH", "/", bytes.NewReader(data))
	instanceID := "instance-id"
	p := httprouter.Params{{"instance_id", instanceID}}

	service, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}
	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	err = broker.updateInstance(w, r, p)
	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusAccepted {
		t.Error("Wrong code:", w.Code)
	}

	<-hook

	if err := verifyJobResult(broker, instanceID, "", jobStatusSucceeded); err != nil {
		t.Error(err)
	}

	if !service.updated {
		t.Error("Should have called update")
	}
}

func TestBroker_UpdateProvisionInstanceMissingFields(t *testing.T) {
	t.Parallel()

	w, r := setupHTTPReq("PATCH", "/", strings.NewReader(`{"plan_id":""}`))
	instanceID := "instance-id"
	p := httprouter.Params{{"instance_id", instanceID}}

	service, _, broker := setupTestBroker()

	err := broker.updateInstance(w, r, p)
	if err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusBadRequest {
		t.Error("Wrong code:", w.Code)
	}

	if service.updated {
		t.Error("Should not have called update")
	}
}

func TestBroker_BindServiceInstance(t *testing.T) {
	t.Parallel()

	var err error
	w, r := setupHTTPReq("PATCH", "/", strings.NewReader(testBindingRequestJSON))
	instanceID := "instance-id"
	bindingID := "binding-id"
	p := httprouter.Params{{"instance_id", instanceID}, {"binding_id", bindingID}}

	service, _, broker := setupTestBroker()

	if err = putInstance(broker.db, instanceID, Instance{}); err != nil {
		t.Error(err)
	}

	if err = broker.bindInstance(w, r, p); err != nil {
		t.Error(err)
	}

	if w.Code != http.StatusCreated {
		t.Error("Wrong code:", w.Code)
	}

	var instance Instance
	if instance, err = getInstance(broker.db, instanceID); err != nil {
		t.Error(err)
	}

	found := false
	for _, binding := range instance.Bindings {
		if binding.BindingID == bindingID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Could not find binding: %#v", instance.Bindings)
	}

	if !service.bound {
		t.Error("Should have called bind")
	}
}

func TestBroker_BindServiceAppGUIDRequired(t *testing.T) {
	t.Parallel()

	var err error
	request := testBindingRequest
	request.AppGUID = ""

	w, r := setupHTTPReq("PATCH", "/", strings.NewReader(testMustEncode(request)))
	instanceID := "instance-id"
	bindingID := "binding-id"
	p := httprouter.Params{{"instance_id", instanceID}, {"binding_id", bindingID}}

	service, _, broker := setupTestBroker()
	broker.Options.RequireAppGUID = true

	if err = putInstance(broker.db, instanceID, Instance{}); err != nil {
		t.Error(err)
	}

	if err = broker.bindInstance(w, r, p); err != nil {
		t.Error(err)
	}

	if w.Code != statusUnprocessableEntity {
		t.Error("Wrong code:", w.Code)
	}

	if s := w.Body.String(); strings.TrimSpace(s) != errMsgRequireAppGUID {
		t.Error("Wrong message:", s)
	}

	if service.bound {
		t.Error("Should not have called bind")
	}
}

func TestBroker_BindServiceInstanceDup(t *testing.T) {
	t.Parallel()

	w, r := setupHTTPReq("PUT", "/", strings.NewReader(testBindingRequestJSON))
	instanceID := "instance-id"
	bindingID := "binding-id"
	p := httprouter.Params{{"instance_id", instanceID}, {"binding_id", bindingID}}

	service, _, broker := setupTestBroker()

	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}
	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	// Bind a first time
	if err := broker.bindInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusCreated {
		t.Error("Wrong code:", w.Code)
	}
	if !service.bound {
		t.Error("Should have called bind")
	}

	// Bind a second time to get the StatusOK on same duplicate
	w, r = setupHTTPReq("PUT", "/", strings.NewReader(testBindingRequestJSON))
	service.bound = false
	if err := broker.bindInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusOK {
		t.Error("Wrong code:", w.Code)
	}
	if service.bound {
		t.Error("Should not have called bind")
	}
}

func TestBroker_BindServiceInstanceDiff(t *testing.T) {
	t.Parallel()

	w, r := setupHTTPReq("PUT", "/", strings.NewReader(testBindingRequestJSON))
	instanceID := "instance-id"
	bindingID := "binding-id"
	p := httprouter.Params{{"instance_id", instanceID}, {"binding_id", bindingID}}

	service, _, broker := setupTestBroker()

	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}
	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	// Bind a first time
	if err := broker.bindInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusCreated {
		t.Error("Wrong code:", w.Code)
	}
	if !service.bound {
		t.Error("Should have called bind")
	}

	// Change a bit of data in the bind request
	newRequest := testBindingRequest
	newRequest.PlanID = testConfig.Services[0].Plans[1].ID
	out, err := json.Marshal(newRequest)
	if err != nil {
		t.Error(err)
	}

	// Bind a second time to get the StatusConflict on diff duplicate
	w, r = setupHTTPReq("PUT", "/", bytes.NewReader(out))
	service.bound = false
	if err := broker.bindInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusConflict {
		t.Error("Wrong code:", w.Code)
	}
	if service.bound {
		t.Error("Should not have called bind")
	}
}

func TestBroker_BindServiceMissingFields(t *testing.T) {
	t.Parallel()

	service, _, broker := setupTestBroker()
	broker.Options.RequireAppGUID = false

	doTest := func(key string, testData BindingRequest, description string) {
		data, err := json.Marshal(testData)
		if err != nil {
			t.Errorf("%s) Failed json marshal: %v %#v", key, err, testData)
		}
		w, r := setupHTTPReq("GET", "/", bytes.NewReader(data))
		p := httprouter.Params{{"instance_id", "instance-id"}, {"binding_id", "binding-id"}}
		if err = broker.bindInstance(w, r, p); err != nil {
			t.Errorf("%s) %v", key, err)
		}
		if w.Code != http.StatusBadRequest {
			t.Errorf("%s) Wrong code: %v", key, w.Code)
		}
		if str := w.Body.String(); !strings.Contains(str, description) {
			t.Errorf("Want: %q, Got: %q", description, str)
		}

		if service.bound {
			t.Error("Should not have called bind")
		}
	}

	testData := testBindingRequest
	testData.ServiceID = ""
	testData.AppGUID = ""
	doTest("ServiceID", testData, "service_id does not reference a service")

	testData = testBindingRequest
	testData.PlanID = ""
	testData.AppGUID = ""
	doTest("PlanID", testData, "plan_id does not reference a services plan")
}

func TestBroker_UnbindServiceInstance(t *testing.T) {
	t.Parallel()

	instanceID := "instance-id"
	bindingID := "binding-id"
	p := httprouter.Params{{"instance_id", instanceID}, {"binding_id", bindingID}}

	service, _, broker := setupTestBroker()

	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}
	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	w, r := setupHTTPReq("PUT", "/", strings.NewReader(testBindingRequestJSON))
	if err := broker.bindInstance(w, r, p); err != nil {
		t.Error(err)
	}

	if gotInstance, err := getInstance(broker.db, testInstanceID); err != nil {
		t.Error(err)
	} else if len(gotInstance.Bindings) == 0 {
		t.Error("Expected some bindings:", gotInstance.Bindings)
	}

	w, r = setupHTTPReq("DELETE", "/", nil)
	if err := broker.unbindInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusOK {
		t.Error("Wrong code:", w.Code)
	}
	if !service.unbound {
		t.Error("Should have called unbind")
	}

	if gotInstance, err := getInstance(broker.db, testInstanceID); err != nil {
		t.Error(err)
	} else if len(gotInstance.Bindings) != 0 {
		t.Error("Expected no bindings:", gotInstance.Bindings)
	}
}

func TestBroker_UnbindServiceInstanceMissing(t *testing.T) {
	t.Parallel()

	w, r := setupHTTPReq("PUT", "/", strings.NewReader(testBindingRequestJSON))
	instanceID := "instance-id"
	bindingID := "binding-id"
	p := httprouter.Params{{"instance_id", instanceID}, {"binding_id", bindingID}}

	service, _, broker := setupTestBroker()

	if err := broker.unbindInstance(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusGone {
		t.Error("Wrong code:", w.Code)
	}
	if service.unbound {
		t.Error("Should not have called unbind")
	}
}

func TestBroker_LastOperation(t *testing.T) {
	t.Parallel()

	w, r := setupHTTPReq("GET", "/", nil)
	p := httprouter.Params{{"instance_id", testInstanceID}}

	_, opener, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	job := testProvisionJob
	job.JobResult.Status = jobStatusInProgress
	opener.keyValues[mkJobKey(testInstanceID)] = testMustEncode(job)

	if err := broker.lastOperation(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusOK {
		t.Error("Wrong code:", w.Code)
	}
	if body := w.Body.String(); strings.TrimSpace(body) != `{"state":"in progress"}` {
		t.Error("Body was wrong:", body)
	}

	w, r = setupHTTPReq("GET", "/", nil)
	p = httprouter.Params{{"instance_id", "notarealid"}}
	if err := broker.lastOperation(w, r, p); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusNotFound {
		t.Error("Wrong code:", w.Code)
	}
}

func TestBroker_EnsureParam(t *testing.T) {
	t.Parallel()

	w, r := setupHTTPReq("GET", "/", nil)
	p := httprouter.Params{}

	if _, ok := ensureParam(w, r, p, paramInstanceID); ok {
		t.Error("Expected no instanceID")
	}
	if w.Code != http.StatusBadRequest {
		t.Error("Wrong code:", w.Code)
	}

	instanceID := "instance-id"
	w, r = setupHTTPReq("GET", "/", nil)
	p = httprouter.Params{{"instance_id", instanceID}}

	if id, ok := ensureParam(w, r, p, paramInstanceID); !ok {
		t.Error("Expected a instanceID")
	} else if id != instanceID {
		t.Error("Wrong value:", id)
	}
	if w.Code != http.StatusOK {
		t.Error("Wrong code:", w.Code)
	}
}

var testConfig = &Options{
	APIVersion:     "2.6",
	AuthUser:       "testuser",
	AuthPassword:   "testpassword",
	RequireAppGUID: true,
	Listen:         ":8080",
	EncryptionKey:  "12345678901234567890123456789012",
	Catalog: Catalog{[]Service{
		Service{
			ID:            "146BF770-252B-4F45-9368-0C2E0EA47A22",
			Bindable:      true,
			Name:          "tester-service",
			Description:   "Test Service",
			Tags:          []string{"test", "service"},
			PlanUpdatable: true,
			Plans: []Plan{
				Plan{
					Name:        "default",
					ID:          "69AD6AD7-5D3D-4D6D-A894-5F1AA0785AE0",
					Description: "Free test plan",
					Free:        true,
				},
				Plan{
					Name:        "second",
					ID:          "4F96C019-B43E-43E8-B278-DDDDDE00D8E5",
					Description: "Second test plan",
					Free:        false,
				}},
		},
		Service{
			ID:            "E5D80C1D-D17A-4E6A-8CB3-640216D5B89D",
			Bindable:      true,
			Name:          "unupdatable-tester-service",
			Description:   "Unupdatable Test Service",
			Tags:          []string{"test", "service"},
			PlanUpdatable: false,
			Plans: []Plan{{
				Name:        "default",
				ID:          "C26CD86E-C5CB-4806-BE1F-78212D54DFF0",
				Description: "Free test plan",
				Free:        true,
			}},
		}},
	},
}

var testConfigJSON = func() []byte {
	j, err := json.Marshal(testConfig)
	if err != nil {
		panic(err)
	}
	return j
}()
