package gocfbroker

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

func hookUpdate() chan int {
	updateChan := make(chan int)
	hookedUpdate = func() {
		updateChan <- 0
	}
	return updateChan
}

func restoreUpdate() {
	hookedUpdate = func() {}
}

var testInstanceID = "instance-id"

var testProvisionJob = brokerJob{
	JobKind:          jobKindProvision,
	InstanceID:       testInstanceID,
	ServiceID:        testConfig.Services[0].ID,
	PlanID:           testConfig.Services[0].Plans[0].ID,
	OrganizationGUID: "someid",
	SpaceGUID:        "anotherid",
	JobResult: jobResult{
		Status: jobStatusInProgress,
	},
}

func TestJobs_Dispatch(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	service, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	if err := broker.queueJob(testProvisionJob); err != nil {
		t.Error(err)
	}

	<-hook

	if !service.provisioned {
		t.Error("Should have called Provision")
	}
}

func TestJobs_DispatchClose(t *testing.T) {
	t.Parallel()

	_, _, broker := setupTestBroker()

	go func() {
		close(broker.jobs)
	}()

	// If this doesn't hang the test, we've exited and succeeded
	broker.dispatcher()
}

func TestJobs_DispatchJobExists(t *testing.T) {
	t.Parallel()

	_, opener, broker := setupTestBroker()

	job := testProvisionJob
	job.JobResult.Status = jobStatusInProgress
	opener.keyValues[mkJobKey(testInstanceID)] = testMustEncode(job)

	go func() {
		broker.jobs <- &job
		if err := <-broker.jobAccepts; err != errJobExists {
			t.Error("Expected errJobExists")
		}
		close(broker.jobs)
	}()

	// If this doesn't hang the test, we've exited
	broker.dispatcher()
}

func TestJobs_DispatchError(t *testing.T) {
	hook := hookUpdate()
	defer restoreUpdate()

	service, _, broker := setupTestBroker()
	go broker.dispatcher()
	defer close(broker.jobs)

	if err := broker.queueJob(testProvisionJob); err != nil {
		t.Error(err)
	}

	<-hook

	if !service.provisioned {
		t.Error("Should have called Provision")
	}

	provReqDiffed := testProvisionJob
	provReqDiffed.PlanID = "notthesame"
	if err := broker.queueJob(provReqDiffed); err != nil {
		t.Error(err)
	}

	<-hook

	res := broker.jobResult(testInstanceID)
	if res.ErrorMsg != errMsgConflict {
		t.Error("Wrong error:", res.ErrorMsg)
	}
}

func TestJobs_DispatchCleanup(t *testing.T) {
	// No parallel to ensure these variables aren't changed in a race-y way.
	saveCleanupTick := cleanupTick
	saveCleanupJobsAfter := cleanupJobsAfter
	cleanupTick = time.Microsecond
	cleanupJobsAfter = time.Microsecond

	defer func() {
		cleanupTick = saveCleanupTick
		cleanupJobsAfter = saveCleanupJobsAfter
	}()

	_, opener, broker := setupTestBroker()
	job := testProvisionJob
	job.JobResult.Status = jobStatusSucceeded
	job.LastUpdated = time.Now().AddDate(-1, 0, 0)
	opener.keyValues[mkJobKey(testInstanceID)] = testMustEncode(job)

	// Dispatch
	go broker.dispatcher()
	defer close(broker.jobs)

	for {
		// This will never exit unless delete succeeds
		res := broker.jobResult(testInstanceID)
		if res.ErrorMsg == errMsgNoJobExists {
			break
		}
	}
}

func TestJobs_Provision(t *testing.T) {
	t.Parallel()

	service, _, broker := setupTestBroker()

	ch := make(chan jobResult)
	go broker.provisionAsync(testInstanceID, testProvisionRequest, ch)

	result, ok := <-ch
	if !ok {
		t.Fatal("Failed to get anything")
	}
	if result.Status != jobStatusSucceeded {
		t.Error("Job failed:", result.Status)
	}
	if len(result.ErrorMsg) != 0 {
		t.Error("Unexpected error:", result.ErrorMsg)
	}

	if _, err := getInstance(broker.db, testInstanceID); err != nil {
		t.Error("It should have saved the instance")
	}

	if !service.provisioned {
		t.Error("Should have called Provision")
	}
}

func TestJobs_Deprovision(t *testing.T) {
	t.Parallel()

	service, _, broker := setupTestBroker()

	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}
	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	ch := make(chan jobResult)
	go broker.deprovisionAsync(testInstanceID, deprovisionRequest{}, ch)

	result, ok := <-ch
	if !ok {
		t.Fatal("Failed to get anything")
	}
	if result.Status != jobStatusSucceeded {
		t.Error("Job failed:", result.Status)
	}
	if len(result.ErrorMsg) != 0 {
		t.Error("Unexpected error:", result.ErrorMsg)
	}

	if _, err := getInstance(broker.db, testInstanceID); !IsKeyNotExist(err) {
		t.Error("Want keynotexist error:", err)
	}

	if !service.deprovisioned {
		t.Error("Should have called Deprovision")
	}
}

func TestJobs_DeprovisionMissing(t *testing.T) {
	t.Parallel()

	service, _, broker := setupTestBroker()

	ch := make(chan jobResult)
	go broker.deprovisionAsync(testInstanceID, deprovisionRequest{}, ch)

	result, ok := <-ch
	if !ok {
		t.Fatal("Failed to get anything")
	}
	if result.ErrorMsg != errMsgGone {
		t.Error("Wrong error:", result.ErrorMsg)
	}

	if service.deprovisioned {
		t.Error("Should not have called Deprovision")
	}
}

func TestJobs_Update(t *testing.T) {
	t.Parallel()

	service, _, broker := setupTestBroker()

	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}

	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	planID := testConfig.Services[0].Plans[1].ID
	ch := make(chan jobResult)
	go broker.updateAsync(testInstanceID, UpdateProvisionRequest{PlanID: planID}, ch)

	result, ok := <-ch
	if !ok {
		t.Fatal("Failed to get anything")
	}
	if result.Status != jobStatusSucceeded {
		t.Error("Job failed:", result.Status)
	}
	if len(result.ErrorMsg) != 0 {
		t.Error("Unexpected error:", result.ErrorMsg)
	}

	if gotInstance, err := getInstance(broker.db, testInstanceID); err != nil {
		t.Error(err)
	} else if gotInstance.ProvisionReq.PlanID != planID {
		t.Error("Plan ID was not updated:", gotInstance.ProvisionReq.PlanID)
	}

	if !service.updated {
		t.Error("Should have called update")
	}
}

func TestJobs_UpdateNoDiff(t *testing.T) {
	t.Parallel()

	service, _, broker := setupTestBroker()

	planID := testConfig.Services[0].Plans[0].ID
	instance := Instance{
		ProvisionReq: testProvisionRequest,
		ProvisionRes: testProvisionResponse,
	}
	instance.ProvisionReq.PlanID = planID

	if err := putInstance(broker.db, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	ch := make(chan jobResult)
	go broker.updateAsync(testInstanceID, UpdateProvisionRequest{PlanID: planID}, ch)

	result, ok := <-ch
	if !ok {
		t.Fatal("Failed to get anything")
	}
	if result.Status != jobStatusSucceeded {
		t.Error("Job failed:", result.Status)
	}
	if result.ErrorMsg != errMsgNoDiff {
		t.Error("Unexpected error:", result.ErrorMsg)
	}

	if service.updated {
		t.Error("Should not have called update")
	}
}

func TestJobs_TimeoutHandling(t *testing.T) {
	saveTimeout := jobTimeout
	jobTimeout = time.Duration(0)

	defer func() {
		jobTimeout = saveTimeout
	}()

	_, _, broker := setupTestBroker()
	broker.provisioner = timeoutProvisioner{}
	job := &brokerJob{
		InstanceID: testInstanceID,
		JobKind:    jobKindProvision,
	}
	results := make(chan jobResult)

	go broker.dispatch(job, results)

	result := <-results

	if result.ErrorMsg != errTimeout.Error() {
		t.Error("Want timed out message:", result.ErrorMsg)
	}
	if result.InstanceID != testInstanceID {
		t.Error("Wrong instance id:", result.InstanceID)
	}
}

func TestJobs_DispatchPanicHandling(t *testing.T) {
	buf := &bytes.Buffer{}
	saveLogger := logger
	logger = log.New(buf, "", 0)

	defer func() {
		logger = saveLogger
	}()

	_, _, broker := setupTestBroker()
	broker.provisioner = panicProvisioner{}
	job := &brokerJob{
		InstanceID: testInstanceID,
		JobKind:    jobKindProvision,
	}
	results := make(chan jobResult)

	go broker.dispatch(job, results)

	result := <-results

	if result.ErrorMsg != errMsgPanic {
		t.Error("Expected an explanation that it panicked:", result.ErrorMsg)
	}
	if result.InstanceID != testInstanceID {
		t.Error("Wrong instance id:", result.InstanceID)
	}
	if result.Status != jobStatusFailed {
		t.Error("Should mark the job as failed:", result.Status)
	}

	logMsg := buf.String()
	if !strings.Contains(logMsg, testInstanceID) {
		t.Error("Expected the instance ID to be in the log msg:", result.ErrorMsg)
	}
	if !strings.Contains(logMsg, `job panicked: "panic provisioner"`) {
		t.Error("The log message should contain explanation about panic:", logMsg)
	}
	if !strings.Contains(logMsg, `gocfbroker/jobs.go:`) {
		t.Error("Want stack trace:", logMsg)
	}
}

// panicProvisioner panics when Provision is called
type panicProvisioner struct{}

func (panicProvisioner) Provision(instanceID string, pr ProvisionRequest) (p ProvisionResponse, e error) {
	panic("panic provisioner")
	return p, e
}
func (panicProvisioner) Deprovision(instanceID, serviceID, planID string) error     { return nil }
func (panicProvisioner) Update(instanceID string, upr UpdateProvisionRequest) error { return nil }
func (panicProvisioner) Bind(instanceID, bindingID string, br BindingRequest) (b BindingResponse, e error) {
	return b, e
}
func (panicProvisioner) Unbind(instanceID, bindingID, serviceID, planID string) error { return nil }

// timeoutProvisioner blocks forever when Provision is called
type timeoutProvisioner struct{}

func (timeoutProvisioner) Provision(instanceID string, pr ProvisionRequest) (p ProvisionResponse, e error) {
	select {}
	return p, e
}
func (timeoutProvisioner) Deprovision(instanceID, serviceID, planID string) error     { return nil }
func (timeoutProvisioner) Update(instanceID string, upr UpdateProvisionRequest) error { return nil }
func (timeoutProvisioner) Bind(instanceID, bindingID string, br BindingRequest) (b BindingResponse, e error) {
	return b, e
}
func (timeoutProvisioner) Unbind(instanceID, bindingID, serviceID, planID string) error { return nil }
