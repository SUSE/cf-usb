package gocfbroker

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"
)

// job timeout variables
var (
	jobTimeout  = 1 * time.Hour
	syncTimeout = 60 * time.Second
	errTimeout  = errors.New("the operation timed out")
)

// these variables control the cleanup routine
var (
	// cleanupTick every X duration
	cleanupTick = 5 * time.Minute
	// cleanupJobsAfter X duration
	cleanupJobsAfter = 1 * time.Hour
)

// jobStatus is the status of a given job
type jobStatus string

const (
	jobStatusInProgress jobStatus = "in progress"
	jobStatusSucceeded  jobStatus = "succeeded"
	jobStatusFailed     jobStatus = "failed"
)

const (
	// errMsgNoJobExists is used to tell an API handler when a job it was looking
	// for does not exist
	errMsgNoJobExists = "no job exists for this instance ID"
	// errMsgNoDiff despite being called an error, a response body should still be returned
	// it just helps the API handlers know that we're returning 200 not 201
	errMsgNoDiff = "the resource has already been created with same settings"
	// errMsgConflict tells the API Handler when a request for this
	// instance/binding ID had already come in, but the bits like PlanID etc
	// were slightly different.
	errMsgConflict = "the resource has already been created with different settings"
	// errMsgGone notes that a resource has already been deleted
	errMsgGone = "the resource has already been deleted"
	// errMsgPanic gives a vague error about a crash to make sure clients don't
	// see large stack traces and panic warnings.
	errMsgPanic = "the job failed unexpectedly"
)

type jobKind string

const (
	jobKindProvision   jobKind = "provision"
	jobKindDeprovision jobKind = "deprovision"
	jobKindUpdate      jobKind = "update"
)

// jobResult is used to tell the dispatcher about the result of an operation
// as well as make requests from the API layer into the dispatcher to discover
// the result of a job.
type jobResult struct {
	InstanceID string
	Status     jobStatus
	// ErrorMsg is one of the constants defined (errMsg*). It can also be an
	// arbitrary error message. The reason for not using error here is because
	// it must be serialized to JSON and after it comes back (even with a custom
	// error type) it can no longer be compared because it's not the same error.
	// So we use strings instead.
	ErrorMsg string
}

// brokerJob is a job given to the broker's dispatcher
type brokerJob struct {
	LastUpdated time.Time `json:"last_updated"`

	JobKind          jobKind          `json:"job_kind"`
	InstanceID       string           `json:"instance_id"`
	ServiceID        string           `json:"service_id"`
	PlanID           string           `json:"plan_id"`
	OrganizationGUID string           `json:"organization_guid"`
	SpaceGUID        string           `json:"space_guid"`
	Parameters       *json.RawMessage `json:"parameters,omitempty"`

	JobResult jobResult `json:"job_result"`

	// Version is the json schema version
	Version int64 `json:"version"`
}

func (b *brokerJob) setSuccess() {
	b.LastUpdated = time.Now().UTC()
	b.JobResult.Status = jobStatusSucceeded
}

func (b *brokerJob) setFailed(err string) {
	b.LastUpdated = time.Now().UTC()
	b.JobResult.ErrorMsg = err
	b.JobResult.Status = jobStatusFailed
}

func (b *brokerJob) update(result jobResult) {
	if result.ErrorMsg == errMsgNoDiff {
		// Special case for deducing that there was a diff error
		// so we can return 200 not 201
		b.setSuccess()
		b.JobResult.ErrorMsg = errMsgNoDiff
	} else if len(result.ErrorMsg) != 0 {
		b.setFailed(result.ErrorMsg)
	} else if result.Status == jobStatusSucceeded {
		b.setSuccess()
	} else {
		panic(fmt.Sprintf("invalid state transition: %#v %#v", b, result))
	}
	hookedUpdate()
}

// hookedUpdate is used for testing the update function in a synchronous way
var hookedUpdate = func() {}

// jobResultRequest is used to request the status of a job.
type jobResultRequest struct {
	instanceID string
	resultCh   chan jobResult
}

// dispatcher has several responsibilities:
// listen for job requests and dispatch them to worker goroutines
// listen for job status requests and do the job
// listen for job results and update the job information with it
// clean expired jobs
func (b *Broker) dispatcher() {
	results := make(chan jobResult)
	ticker := time.NewTicker(cleanupTick)

	for {
		select {
		case job, ok := <-b.jobs:
			if !ok {
				ticker.Stop()
				return
			}

			if err := putJob(b.db, job); err != nil {
				b.jobAccepts <- err
				continue
			}

			go b.dispatch(job, results)
			b.jobAccepts <- nil
		case resultReq := <-b.jobResultReqs:
			j, err := getJob(b.db, resultReq.instanceID)
			if IsKeyNotExist(err) {
				resultReq.resultCh <- jobResult{ErrorMsg: errMsgNoJobExists}
				continue
			} else if err != nil {
				resultReq.resultCh <- jobResult{ErrorMsg: err.Error()}
				continue
			}

			resultReq.resultCh <- j.JobResult
		case result := <-results:
			logger.Printf("updating job (%s): %v", result.InstanceID, result.Status)
			if err := updateJob(b.db, result); err != nil {
				logger.Printf("error updating job id: %s, err: %v", result.InstanceID, err)
			}
		case <-ticker.C:
			logger.Println("reaping jobs")
			if err := reapJobs(b.db); err != nil {
				logger.Println("error reaping jobs:", err)
			}
		}
	}
}

// dispatch a job to a handler
func (b *Broker) dispatch(job *brokerJob, results chan<- jobResult) {
	defer dispatchRecover(job, results)

	switch job.JobKind {
	case jobKindProvision:
		preq := ProvisionRequest{
			ServiceID:        job.ServiceID,
			PlanID:           job.PlanID,
			OrganizationGUID: job.OrganizationGUID,
			SpaceGUID:        job.SpaceGUID,
		}
		logger.Printf("dispatching (%s) [%s]: ServiceID: %q, PlanID: %q OrganizationGUID: %q SpaceGUID: %q",
			job.InstanceID, job.JobKind, preq.ServiceID, preq.PlanID, preq.OrganizationGUID, preq.SpaceGUID)
		b.provisionAsync(job.InstanceID, preq, results)
	case jobKindDeprovision:
		dpreq := deprovisionRequest{
			serviceID: job.ServiceID,
			planID:    job.PlanID,
		}
		logger.Printf("dispatching (%s) [%s]: ServiceID: %q, PlanID: %q",
			job.InstanceID, job.JobKind, dpreq.serviceID, dpreq.planID)
		b.deprovisionAsync(job.InstanceID, dpreq, results)
	case jobKindUpdate:
		upreq := UpdateProvisionRequest{
			ServiceID: job.ServiceID,
			PlanID:    job.PlanID,
		}
		logger.Printf("dispatching (%s) [%s]: ServiceID: %q, PlanID: %q",
			job.InstanceID, job.JobKind, upreq.ServiceID, upreq.PlanID)
		b.updateAsync(job.InstanceID, upreq, results)
	}
}

func dispatchRecover(job *brokerJob, results chan<- jobResult) {
	r := recover()
	if r == nil {
		return
	}

	var recovered interface{}
	var stack string

	if perr, ok := r.(panicErr); ok {
		recovered = perr.recovered
		stack = perr.stack
	} else {
		recovered = r
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		stack = string(buf[:n])
	}

	logger.Printf("job panicked: \"%v\" job: %#v\n%s", recovered, job, stack)

	// Fail the job
	var instanceID string
	if job != nil {
		instanceID = job.InstanceID
	}
	results <- jobResult{
		InstanceID: instanceID,
		Status:     jobStatusFailed,
		ErrorMsg:   errMsgPanic,
	}
}

// queueJob creates a job and puts it into the dispatcher
func (b *Broker) queueJob(job brokerJob) error {
	job.Version = brokerJSONSchemaVersion
	job.LastUpdated = time.Now().UTC()
	b.jobs <- &job
	return <-b.jobAccepts
}

// jobResult retrieves the result the last job for the instanceID
func (b *Broker) jobResult(instanceID string) jobResult {
	resultRequest := jobResultRequest{
		instanceID: instanceID,
		resultCh:   make(chan jobResult),
	}
	b.jobResultReqs <- resultRequest
	return <-resultRequest.resultCh
}

func (b *Broker) provisionAsync(instanceID string, req ProvisionRequest, ch chan<- jobResult) {
	result := jobResult{InstanceID: instanceID}

	// Check if a service instance already exists
	instance, err := getInstance(b.db, instanceID)
	if err == nil {
		if instance.ProvisionReq.Equal(req) {
			result.Status = jobStatusSucceeded
			result.ErrorMsg = errMsgNoDiff
		} else {
			result.ErrorMsg = errMsgConflict
		}
		ch <- result
		return
	} else if !IsKeyNotExist(err) {
		result.ErrorMsg = err.Error()
		ch <- result
		return
	}

	// Provision a service
	instance.ProvisionReq = req
	err = timeoutable(jobTimeout, func() error {
		var innerErr error
		instance.ProvisionRes, innerErr = b.provisioner.Provision(instanceID, instance.ProvisionReq)
		return innerErr
	})

	if err != nil {
		result.ErrorMsg = err.Error()
		ch <- result
		return
	}

	// Store the service provisioning in the database.
	if err = putInstance(b.db, instanceID, instance); err != nil {
		errMsg := fmt.Sprintf("storing provisioned instance failed (there may be an orphaned service) for id: %v, err: %v", instanceID, err)
		logger.Print(errMsg)
		result.ErrorMsg = errMsg
		ch <- result
		return
	}

	result.Status = jobStatusSucceeded
	ch <- result
}

func (b *Broker) provisionSync(instanceID string, req ProvisionRequest) (result jobResult) {
	result = jobResult{InstanceID: instanceID}

	// Check if a service instance already exists
	instance, err := getInstance(b.db, instanceID)
	if err == nil {
		if instance.ProvisionReq.Equal(req) {
			result.Status = jobStatusSucceeded
			result.ErrorMsg = errMsgNoDiff
		} else {
			result.ErrorMsg = errMsgConflict
		}
		return
	} else if !IsKeyNotExist(err) {
		result.ErrorMsg = err.Error()
		return
	}

	// Provision a service
	instance.ProvisionReq = req
	provionError := make(chan (error))
	go func() {
		var innerErr error
		instance.ProvisionRes, innerErr = b.provisioner.Provision(instanceID, instance.ProvisionReq)
		provionError <- innerErr
	}()

	//Wait for provisioner to finish
	select {
	case err := <-provionError:
		if err != nil {
			result.ErrorMsg = err.Error()
		}
	case <-time.After(syncTimeout):
		return jobResult{Status: jobStatusFailed, ErrorMsg: "Timed out provisioning service", InstanceID: instanceID}
	}

	if err != nil {
		result.ErrorMsg = err.Error()
		return
	}

	// Store the service provisioning in the database.
	if err = putInstance(b.db, instanceID, instance); err != nil {
		errMsg := fmt.Sprintf("storing provisioned instance failed (there may be an orphaned service) for id: %v, err: %v", instanceID, err)
		logger.Print(errMsg)
		result.ErrorMsg = errMsg
		return
	}

	result.Status = jobStatusSucceeded
	return result
}

func (b *Broker) deprovisionAsync(instanceID string, req deprovisionRequest, ch chan<- jobResult) {
	result := jobResult{InstanceID: instanceID}

	// Delete the instance from the database
	err := delInstance(b.db, instanceID)
	if IsKeyNotExist(err) {
		result.ErrorMsg = errMsgGone
	} else if err != nil {
		result.ErrorMsg = err.Error()
	}

	if len(result.ErrorMsg) != 0 {
		ch <- result
		return
	}

	err = timeoutable(jobTimeout, func() error {
		return b.provisioner.Deprovision(instanceID, req.serviceID, req.planID)
	})

	if err != nil {
		errMsg := fmt.Sprintf("failed to deprovision instance, this may require manual deprovisioning id: %s, err: %v", instanceID, err)
		logger.Print(errMsg)
		result.ErrorMsg = errMsg
		ch <- result
		return
	}

	result.Status = jobStatusSucceeded
	ch <- result
}

func (b *Broker) deprovisionSync(instanceID string, req deprovisionRequest) (result jobResult) {
	result = jobResult{InstanceID: instanceID}

	// Delete the instance from the database
	err := delInstance(b.db, instanceID)
	if IsKeyNotExist(err) {
		result.ErrorMsg = errMsgGone
	} else if err != nil {
		result.ErrorMsg = err.Error()
	}

	if len(result.ErrorMsg) != 0 {
		return result
	}

	deprovisionError := make(chan (error))

	go func() {
		deprovisionError <- b.provisioner.Deprovision(instanceID, req.serviceID, req.planID)
	}()

	//Wait for provisioner to finish
	select {
	case err := <-deprovisionError:
		if err != nil {
			result.ErrorMsg = err.Error()
		}
	case <-time.After(syncTimeout):
		return jobResult{Status: jobStatusFailed, ErrorMsg: "Timed out deprovisioning service instance", InstanceID: instanceID}
	}

	result.Status = jobStatusSucceeded
	return result
}

func (b *Broker) updateAsync(instanceID string, req UpdateProvisionRequest, ch chan<- jobResult) {
	result := jobResult{InstanceID: instanceID}

	// Check if a service instance already exists
	instance, err := getInstance(b.db, instanceID)
	if err == nil {
		if instance.ProvisionReq.PlanID == req.PlanID {
			result.ErrorMsg = errMsgNoDiff
			result.Status = jobStatusSucceeded
			ch <- result
			return
		}
	} else {
		result.ErrorMsg = err.Error()
		ch <- result
		return
	}

	err = timeoutable(jobTimeout, func() error {
		return b.provisioner.Update(instanceID, req)
	})

	if err != nil {
		result.ErrorMsg = err.Error()
		ch <- result
		return
	}

	instance.ProvisionReq.PlanID = req.PlanID
	if err = putInstance(b.db, instanceID, instance); err != nil {
		errMsg := fmt.Sprintf("failed to update instance, inconsistency may exist between actual instance and what the broker knows for id: %v, err: %v", instanceID, err)
		logger.Print(errMsg)
		result.ErrorMsg = errMsg
		ch <- result
		return
	}

	result.Status = jobStatusSucceeded
	ch <- result
}

func (b *Broker) updateSync(instanceID string, req UpdateProvisionRequest) (result jobResult) {
	result = jobResult{InstanceID: instanceID}

	// Check if a service instance already exists
	instance, err := getInstance(b.db, instanceID)
	if err == nil {
		if instance.ProvisionReq.PlanID == req.PlanID {
			result.ErrorMsg = errMsgNoDiff
			result.Status = jobStatusSucceeded
			return
		}
	} else {
		result.ErrorMsg = err.Error()
		return
	}

	updateError := make(chan (error))

	go func() {
		updateError <- b.provisioner.Update(instanceID, req)
	}()

	//Wait for provisioner to finish
	select {
	case err := <-updateError:
		if err != nil {
			result.ErrorMsg = err.Error()
		}
	case <-time.After(syncTimeout):
		return jobResult{Status: jobStatusFailed, ErrorMsg: "Timed out updating service instance", InstanceID: instanceID}
	}

	instance.ProvisionReq.PlanID = req.PlanID
	if err = putInstance(b.db, instanceID, instance); err != nil {
		errMsg := fmt.Sprintf("failed to update instance, inconsistency may exist between actual instance and what the broker knows for id: %v, err: %v", instanceID, err)
		logger.Print(errMsg)
		result.ErrorMsg = errMsg
		return
	}

	result.Status = jobStatusSucceeded
	return result
}

type panicErr struct {
	recovered interface{}
	stack     string
}

func (p panicErr) Error() string {
	return errMsgPanic
}

// timeoutable calls a callback in a goroutine and ensures it completes within
// a certain amount of time. Returns errTimeout if a timeout occured, else it
// returns what was returned from callback().
func timeoutable(wait time.Duration, callback func() error) (err error) {
	ch := make(chan error, 1)

	// Run callback
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				var buf [4096]byte
				n := runtime.Stack(buf[:], false)
				ch <- panicErr{recovered: rec, stack: string(buf[:n])}
			}
		}()

		ch <- callback()
	}()

	timer := time.NewTimer(wait)

	select {
	case <-timer.C:
		return errTimeout
	case err = <-ch:
	}

	timer.Stop()

	// Panic this up to something that has a clue about what acutally panic'd:
	// the dispatch call.
	if perr, ok := err.(panicErr); ok {
		panic(perr)
	}

	return err
}
