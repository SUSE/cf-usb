package gocfbroker

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	errJobFinished = errors.New("job with this id is already finished")
	errJobExists   = errors.New("job for this instance id is already in progress")
	errNoJobExists = errors.New("job not found for this id")
)

const (
	// suffixJob is used to suffix instance ID's to differentiate them from the
	// actual stored instance.
	suffixJob = "_job"
)

const (
	// StoreNoLock is used to signify that we didn't do a Get() to retrieve
	// a lock value and we want to ignore the compare-and-swap semantics.
	StoreNoLock = -1
)

// Storer implements a key-value storage unit. It supports compare-and-swap
// to ensure consistency.
type Storer interface {
	// Close this database connection
	Close() error

	// Get the value for the given key. If the key was not found, the
	// implementation should return an error made with ErrKeyNotExist(key).
	// The lock value is used as input to Put() and Del() to ensure that the
	// value hasn't changed between get and sets.
	Get(key string) (value string, lock int, err error)
	// Put a new value for the given key, overwrites if key already exists.
	// Creates if not already exists. Takes a lock value from Get() or StoreNoLock
	// in order to ensure (or not) that no other writes happened to the data
	// between the last read and the write (compare-and-swap semantic).
	Put(key, value string, lock int) error
	// Delete a key from the StoreWriter. If the key was not found, the
	// implementation should return an error made with ErrKeyNotExist(key).
	// This is to avoid data races within the API, as delete must return
	// a "StatusGone" if the service/binding is already gone.
	Del(key string) error
	// Keys returns all the keys with the given suffix.
	Keys(suffix string) ([]string, error)
}

// NewJSONStorer wraps a storer in JSON helper methods.
func NewJSONStorer(storer Storer) JSONStorer {
	return JSONStorer{Storer: storer}
}

// JSONStorer wraps the Storer interfaces above with the ability to
// read and write json directly using all the above functionality.
type JSONStorer struct {
	Storer
}

// Put a JSON object into the internal storer.
func (j JSONStorer) Put(key string, obj interface{}, lock int) error {
	out, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return j.Storer.Put(key, string(out), lock)
}

// Get a JSON object from the internal storer.
func (j JSONStorer) Get(key string, obj interface{}) (int, error) {
	out, lock, err := j.Storer.Get(key)
	if err != nil {
		return StoreNoLock, err
	}

	return lock, json.Unmarshal([]byte(out), obj)
}

func putInstance(s JSONStorer, instanceID string, instance Instance) error {
	instance.Version = brokerJSONSchemaVersion
	return s.Put(instanceID, instance, StoreNoLock)
}

func delInstance(s JSONStorer, instanceID string) error {
	return s.Del(instanceID)
}

func getInstance(s JSONStorer, instanceID string) (instance Instance, err error) {
	if _, err := s.Get(instanceID, &instance); err != nil {
		return instance, err
	}

	if err == nil {
		err = checkJSONVersion(instance.Version)
	}

	return instance, err
}

func putBinding(s JSONStorer, secretKey []byte, instanceID, bindingID string, brq BindingRequest, brs BindingResponse) error {
	bindingJSON, err := json.Marshal(brs)
	if err != nil {
		return err
	}

	// Encrypt the binding data before storing in the db
	encryptedBindingResponse, err := encrypt(secretKey, bindingJSON)
	if err != nil {
		return err
	}

	return exponentialJitter(func() error {
		var instance Instance
		lock, err := s.Get(instanceID, &instance)
		if err != nil {
			return err
		}

		if err = checkJSONVersion(instance.Version); err != nil {
			return err
		}

		instance.Bindings = append(instance.Bindings, Binding{
			BindingID:  bindingID,
			BindingReq: brq,
			BindingRes: string(encryptedBindingResponse),
		})

		err = s.Put(instanceID, &instance, lock)
		if err == ErrStaleData {
			return errBackoffRetry
		}
		return err
	})
}

func delBinding(s JSONStorer, instanceID, bindingID string) error {
	return exponentialJitter(func() error {
		var instance Instance
		lock, err := s.Get(instanceID, &instance)
		if err != nil {
			return err
		}

		if err = checkJSONVersion(instance.Version); err != nil {
			return err
		}

		found := -1
		for i, binding := range instance.Bindings {
			if binding.BindingID == bindingID {
				found = i
				break
			}
		}

		if found < 0 {
			return nil
		}

		ln := len(instance.Bindings) - 1
		instance.Bindings[found] = instance.Bindings[ln]
		instance.Bindings = instance.Bindings[:ln]

		err = s.Put(instanceID, &instance, lock)
		if err == ErrStaleData {
			return errBackoffRetry
		}
		return err
	})
}

func getBinding(s JSONStorer, secretKey []byte, instanceID, bindingID string) (brq BindingRequest, brs BindingResponse, err error) {
	var instance Instance
	if _, err := s.Get(instanceID, &instance); err != nil {
		return brq, brs, err
	}

	if err := checkJSONVersion(instance.Version); err != nil {
		return brq, brs, err
	}

	found := -1
	var encryptedBindingResponse string
	for i, binding := range instance.Bindings {
		if binding.BindingID == bindingID {
			found = i
			brq = binding.BindingReq
			encryptedBindingResponse = binding.BindingRes
			break
		}
	}

	if found < 0 {
		return brq, brs, ErrKeyNotExist(bindingID)
	}

	decrypted, err := decrypt(secretKey, encryptedBindingResponse)
	if err != nil {
		return brq, brs, err
	}
	if err := json.Unmarshal(decrypted, &brs); err != nil {
		return brq, brs, err
	}

	return brq, brs, err
}

// putJob should fail when there is a job that exists that is not complete
func putJob(s JSONStorer, job *brokerJob) error {
	jobKey := mkJobKey(job.InstanceID)

	var currentJob brokerJob

	return exponentialJitter(func() error {
		lock, err := s.Get(jobKey, &currentJob)
		switch {
		case err == nil && currentJob.JobResult.Status == jobStatusInProgress:
			return errJobExists
		case err != nil && !IsKeyNotExist(err):
			return err
		}

		err = s.Put(jobKey, job, lock)
		if err == ErrStaleData {
			return errBackoffRetry
		}
		return err
	})
}

// updateJob should fail when there is a job that is finished being updated,
// or when the job doesn't exist
func updateJob(s JSONStorer, result jobResult) error {
	jobKey := mkJobKey(result.InstanceID)

	var currentJob brokerJob

	return exponentialJitter(func() error {
		lock, err := s.Get(jobKey, &currentJob)

		if IsKeyNotExist(err) {
			return errNoJobExists
		} else if err != nil {
			return err
		}

		if err = checkJSONVersion(currentJob.Version); err != nil {
			return err
		}

		if currentJob.JobResult.Status != jobStatusInProgress {
			return errJobFinished
		}

		currentJob.update(result)

		err = s.Put(jobKey, currentJob, lock)
		if err == ErrStaleData {
			return errBackoffRetry
		}
		return err
	})
}

func getJob(s JSONStorer, instanceID string) (job *brokerJob, err error) {
	jobKey := mkJobKey(instanceID)

	job = &brokerJob{}
	if _, err := s.Get(jobKey, &job); err != nil {
		return nil, err
	}

	if err = checkJSONVersion(job.Version); err != nil {
		return nil, err
	}

	return job, err
}

func reapJobs(s JSONStorer) error {
	keys, err := s.Keys(suffixJob)
	if err != nil {
		return err
	}

	for _, key := range keys {
		var job brokerJob
		if _, err = s.Get(key, &job); err != nil {
			logger.Printf("skipping reaping job: %s, %v\n%#v", key, err, job)
			continue
		}

		if err = checkJSONVersion(job.Version); err != nil {
			logger.Printf("skipping reaping job: %s, %v\n%#v", key, err, job)
			continue
		}

		timeSinceLastUpdate := time.Now().UTC().Sub(job.LastUpdated)
		if job.JobResult.Status == jobStatusInProgress || timeSinceLastUpdate < cleanupJobsAfter {
			continue
		}

		logger.Println("reaping job:", key)
		if err = s.Del(key); err != nil {
			logger.Printf("failed reaping (delete) job: %s, %v\n%#v", key, err, job)
		}
	}

	return nil
}

// mkJobID returns a key for a job for the instanceID
func mkJobKey(instanceID string) string {
	return instanceID + suffixJob
}

// checkJSONVersion to make sure it's appropriate to be decoded
func checkJSONVersion(version int64) error {
	if version > brokerJSONSchemaVersion {
		return errJSONVersionMismatch(version)
	}
	return nil
}
