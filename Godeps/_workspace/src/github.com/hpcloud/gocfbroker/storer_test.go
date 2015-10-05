package gocfbroker

import (
	"strings"
	"testing"
	"time"
)

func newTestStore() testStore {
	return testStore{make(map[string]string)}
}

type testStore struct {
	keyValues map[string]string
}

func (t testStore) Get(key string) (val string, lock int, err error) {
	val, ok := t.keyValues[key]
	if !ok {
		return val, StoreNoLock, ErrKeyNotExist(key)
	}
	return val, 1, nil
}

func (t testStore) Put(key, value string, lock int) error {
	t.keyValues[key] = value
	return nil
}

func (t testStore) Del(key string) error {
	if _, ok := t.keyValues[key]; !ok {
		return ErrKeyNotExist(key)
	}
	delete(t.keyValues, key)
	return nil
}

func (t testStore) Keys(suffix string) ([]string, error) {
	keys := make([]string, len(t.keyValues))
	for key := range t.keyValues {
		if strings.HasSuffix(key, suffix) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (t testStore) Close() error { return nil }

type jitterStore struct {
	testStore
	failedOnce bool
}

func (j *jitterStore) Put(key, val string, lock int) (err error) {
	if j.failedOnce {
		j.testStore.Put(key, val, lock)
		return nil
	}

	j.failedOnce = true
	return ErrStaleData
}

func TestStorer_PutInstance(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	instance := Instance{ProvisionReq: testProvisionRequest, ProvisionRes: testProvisionResponse}
	if err := putInstance(s, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	gotInstance, err := getInstance(s, testInstanceID)
	if err != nil {
		t.Error(err)
	}
	if !gotInstance.ProvisionReq.Equal(testProvisionRequest) {
		t.Errorf("Request not stored properly: %#v", gotInstance.ProvisionReq)
	}

	if gotInstance.ProvisionRes != testProvisionResponse {
		t.Errorf("Response not stored properly: %#v", gotInstance.ProvisionRes)
	}
}

func TestStorer_GetInstance(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	_, err := getInstance(s, testInstanceID)
	if !IsKeyNotExist(err) {
		t.Error("Expected key not found error:", err)
	}

	instance := Instance{ProvisionReq: testProvisionRequest, ProvisionRes: testProvisionResponse}
	if err := putInstance(s, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	gotInstance, err := getInstance(s, testInstanceID)
	if err != nil {
		t.Error(err)
	}
	if !gotInstance.ProvisionReq.Equal(testProvisionRequest) {
		t.Errorf("Request not stored properly: %#v", gotInstance.ProvisionReq)
	}

	if gotInstance.ProvisionRes != testProvisionResponse {
		t.Errorf("Response not stored properly: %#v", gotInstance.ProvisionRes)
	}
}

func TestStorer_DelInstance(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	err := delInstance(s, testInstanceID)
	if !IsKeyNotExist(err) {
		t.Error("Expected key not found error:", err)
	}

	instance := Instance{ProvisionReq: testProvisionRequest, ProvisionRes: testProvisionResponse}
	if err := putInstance(s, testInstanceID, instance); err != nil {
		t.Error(err)
	}

	err = delInstance(s, testInstanceID)
	if err != nil {
		t.Error(err)
	}

	if !IsKeyNotExist(delInstance(s, testInstanceID)) {
		t.Error(err)
	}
}

func TestStorer_PutBinding(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	bindingID := "binding"
	key := []byte(strings.Repeat("0", 32))

	if err := putInstance(s, testInstanceID, Instance{}); err != nil {
		t.Error(err)
	}

	if err := putBinding(s, key, testInstanceID, bindingID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}

	instance, err := getInstance(s, testInstanceID)
	if err != nil {
		t.Error(err)
	}
	if !instance.Bindings[0].BindingReq.Equal(testBindingRequest) {
		t.Error("BindRequest saved wrong:", instance.Bindings[0].BindingReq)
	}
	if len(instance.Bindings[0].BindingRes) == 0 {
		t.Error("BindResponse not saved:", instance.Bindings[0].BindingRes)
	}

	err = putBinding(s, key, "no-instance-id", bindingID, testBindingRequest, testBindingResponse)
	if !IsKeyNotExist(err) {
		t.Error(err)
	}
}

func TestStorer_PutBindingJitter(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	bindingID := "binding"
	key := []byte(strings.Repeat("0", 32))

	if err := putInstance(s, testInstanceID, Instance{}); err != nil {
		t.Error(err)
	}

	jitterStorer := &jitterStore{testStore: store}
	jitterJSON := NewJSONStorer(jitterStorer)
	if err := putBinding(jitterJSON, key, testInstanceID, bindingID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}

	if !jitterStorer.failedOnce {
		t.Error("Should have failed once.")
	}
}

func TestStorer_DelBinding(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	binding1ID := "binding1"
	binding2ID := "binding2"
	key := []byte(strings.Repeat("0", 32))

	if err := putInstance(s, testInstanceID, Instance{}); err != nil {
		t.Error(err)
	}

	if err := putBinding(s, key, testInstanceID, binding1ID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}
	if err := putBinding(s, key, testInstanceID, binding2ID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}

	if err := delBinding(s, testInstanceID, binding1ID); err != nil {
		t.Error(err)
	}

	if err := delBinding(s, testInstanceID, "notarealbindingid"); err != nil {
		t.Error("It should not produce an error when given a binding that doesn't exist:", err)
	}

	instance, err := getInstance(s, testInstanceID)
	if err != nil {
		t.Error(err)
	}
	if len(instance.Bindings) != 1 {
		t.Error("One binding should have been deleted, one should still be there.")
	}
	if instance.Bindings[0].BindingID != binding2ID {
		t.Error("Binding2 should have been preserved.")
	}
}

func TestStorer_DelBindingJitter(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	bindingID := "binding"
	key := []byte(strings.Repeat("0", 32))

	if err := putInstance(s, testInstanceID, Instance{}); err != nil {
		t.Error(err)
	}
	if err := putBinding(s, key, testInstanceID, bindingID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}

	jitterStorer := &jitterStore{testStore: store}
	jitterJSON := NewJSONStorer(jitterStorer)
	if err := delBinding(jitterJSON, testInstanceID, bindingID); err != nil {
		t.Error(err)
	}

	if !jitterStorer.failedOnce {
		t.Error("Should have failed once.")
	}
}

func TestStorer_DelBindingMissingInstance(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	if err := delBinding(s, testInstanceID, "binding-id"); !IsKeyNotExist(err) {
		t.Error(err)
	}
}

func TestStorer_DelBindingMissingBinding(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	instance1ID := "instance1"
	instance2ID := "instance2"
	binding1ID := "binding1"
	binding2ID := "binding2"
	key := []byte(strings.Repeat("0", 32))

	if err := putInstance(s, instance1ID, Instance{}); err != nil {
		t.Error(err)
	}
	if err := putInstance(s, instance2ID, Instance{}); err != nil {
		t.Error(err)
	}

	if err := putBinding(s, key, instance1ID, binding1ID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}
	if err := putBinding(s, key, instance2ID, binding2ID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}

	if err := delBinding(s, instance2ID, binding1ID); err != nil {
		t.Error(err)
	}
	if err := delBinding(s, instance1ID, binding2ID); err != nil {
		t.Error(err)
	}

	if instance, err := getInstance(s, instance1ID); err != nil {
		t.Error(err)
	} else if len(instance.Bindings) != 1 {
		t.Error("No binding should have been deleted:", instance.Bindings)
	} else if id := instance.Bindings[0].BindingID; id != binding1ID {
		t.Error("Binding ID wrong:", id)
	}
	if instance, err := getInstance(s, instance2ID); err != nil {
		t.Error(err)
	} else if len(instance.Bindings) != 1 {
		t.Error("No binding should have been deleted:", instance.Bindings)
	} else if id := instance.Bindings[0].BindingID; id != binding2ID {
		t.Error("Binding ID wrong:", id)
	}
}

func TestStorer_GetBinding(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	bindingID := "binding"
	key := []byte(strings.Repeat("0", 32))

	if err := putInstance(s, testInstanceID, Instance{}); err != nil {
		t.Error(err)
	}

	if err := putBinding(s, key, testInstanceID, bindingID, testBindingRequest, testBindingResponse); err != nil {
		t.Error(err)
	}

	brq, brs, err := getBinding(s, key, testInstanceID, bindingID)
	if err != nil {
		t.Error(err)
	}
	if !brq.Equal(testBindingRequest) {
		t.Errorf("Want: %#v Got: %#v", testBindingRequest, brq)
	}
	if !brs.Equal(testBindingResponse) {
		t.Errorf("Want: %#v Got: %#v", testBindingResponse, brs)
	}
}

func TestStorer_PutJob(t *testing.T) {
	t.Parallel()

	var b brokerJob
	b.InstanceID = testInstanceID
	store := newTestStore()
	s := NewJSONStorer(store)

	if err := putJob(s, &b); err != nil {
		t.Error(err)
	}

	if _, ok := store.keyValues[mkJobKey(testInstanceID)]; !ok {
		t.Error("Expected something in the DB")
	}
}

func TestStorer_PutJobWhenIncompleteJobPresent(t *testing.T) {
	t.Parallel()

	var b brokerJob
	b.InstanceID = testInstanceID
	b.JobResult.Status = jobStatusInProgress
	store := newTestStore()
	s := NewJSONStorer(store)

	if err := putJob(s, &b); err != nil {
		t.Error(err)
	}

	if err := putJob(s, &b); err != errJobExists {
		t.Error("Expected exists error:", err)
	}
}

func TestStorer_PutJobJitter(t *testing.T) {
	t.Parallel()

	var b brokerJob
	b.InstanceID = testInstanceID

	store := newTestStore()
	jitterStorer := &jitterStore{testStore: store}
	jitterJSON := NewJSONStorer(jitterStorer)

	if err := putJob(jitterJSON, &b); err != nil {
		t.Error(err)
	}

	if !jitterStorer.failedOnce {
		t.Error("Should have failed once.")
	}
}

func TestStorer_UpdateJob(t *testing.T) {
	t.Parallel()

	var b brokerJob
	b.InstanceID = testInstanceID
	b.JobResult.Status = jobStatusInProgress
	b.JobResult.InstanceID = testInstanceID

	store := newTestStore()
	s := NewJSONStorer(store)

	if err := updateJob(s, b.JobResult); err != errNoJobExists {
		t.Error("Expect it to fail trying to update non existent thing:", err)
	}
	if err := putJob(s, &b); err != nil {
		t.Error(err)
	}

	b.JobResult.Status = jobStatusSucceeded
	if err := updateJob(s, b.JobResult); err != nil {
		t.Error(err)
	}
}

func TestStorer_UpdateJobJitter(t *testing.T) {
	t.Parallel()

	var b brokerJob
	b.InstanceID = testInstanceID
	b.JobResult.Status = jobStatusInProgress
	b.JobResult.InstanceID = testInstanceID

	store := newTestStore()
	s := NewJSONStorer(store)
	jitterStorer := &jitterStore{testStore: store}
	jitterJSON := NewJSONStorer(jitterStorer)

	if err := putJob(s, &b); err != nil {
		t.Error(err)
	}

	b.JobResult.Status = jobStatusSucceeded
	if err := updateJob(jitterJSON, b.JobResult); err != nil {
		t.Error(err)
	}

	if !jitterStorer.failedOnce {
		t.Error("Should have failed once.")
	}
}

func TestStorer_UpdateJobUpdateWhenAlreadyComplete(t *testing.T) {
	t.Parallel()

	var b brokerJob
	b.InstanceID = testInstanceID
	b.JobResult.Status = jobStatusSucceeded
	b.JobResult.InstanceID = testInstanceID
	store := newTestStore()
	s := NewJSONStorer(store)

	if err := putJob(s, &b); err != nil {
		t.Error(err)
	}

	if err := updateJob(s, b.JobResult); err != errJobFinished {
		t.Error("Expected finished error:", err)
	}
}

func TestStorer_GetJob(t *testing.T) {
	t.Parallel()

	var b brokerJob
	b.InstanceID = testInstanceID
	store := newTestStore()
	s := NewJSONStorer(store)

	if err := putJob(s, &b); err != nil {
		t.Error(err)
	}

	if _, ok := store.keyValues[mkJobKey(testInstanceID)]; !ok {
		t.Error("Expected something in the DB")
	}

	if job, err := getJob(s, testInstanceID); err != nil {
		t.Error(err)
	} else if job.InstanceID != testInstanceID {
		t.Error("Instance ID was wrong:", job.InstanceID)
	}

	if _, err := getJob(s, "nonexistentid"); !IsKeyNotExist(err) {
		t.Error(err)
	}
}

func TestStorer_ReapJobs(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	instance1ID := "instance-id1"
	instance2ID := "instance-id2"
	instance3ID := "instance-id3"
	instance4ID := "instance-id4"

	var b1, b2, b3, b4 brokerJob
	b1.InstanceID = instance1ID
	b1.LastUpdated = time.Now().UTC().AddDate(0, 0, -1)
	b1.JobResult.Status = jobStatusSucceeded

	b2.InstanceID = instance2ID
	b2.JobResult.Status = jobStatusSucceeded
	b2.LastUpdated = time.Now().UTC()

	b3.InstanceID = instance3ID
	b3.JobResult.Status = jobStatusInProgress
	b3.LastUpdated = time.Now().UTC().AddDate(0, 0, -1)

	b4.InstanceID = instance4ID
	b4.Version = 9999
	b4.LastUpdated = time.Now().UTC().AddDate(0, 0, -1)
	b4.JobResult.Status = jobStatusSucceeded

	if err := putJob(s, &b1); err != nil {
		t.Error(err)
	}
	if err := putJob(s, &b2); err != nil {
		t.Error(err)
	}
	if err := putJob(s, &b3); err != nil {
		t.Error(err)
	}
	if err := s.Put(mkJobKey(instance4ID), b4, StoreNoLock); err != nil {
		t.Error(err)
	}

	if err := reapJobs(s); err != nil {
		t.Error(err)
	}

	if _, ok := store.keyValues[mkJobKey(instance1ID)]; ok {
		t.Error("This instance should have been reaped (finished and old)")
	}
	if _, ok := store.keyValues[mkJobKey(instance2ID)]; !ok {
		t.Error("This instance should have been spared (finished but still new)")
	}
	if _, ok := store.keyValues[mkJobKey(instance3ID)]; !ok {
		t.Error("This instance should have been spared (not finished, even though old)")
	}
	if _, ok := store.keyValues[mkJobKey(instance4ID)]; !ok {
		t.Error("This instance should have been skipped (bad json version value)")
	}
}

func TestStorer_InstanceVersionMismatching(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)

	store.keyValues[testInstanceID] = `{"version": 9999}`

	if _, err := getInstance(s, testInstanceID); err == nil {
		t.Error("Expected a failure")
	} else if !strings.Contains(err.Error(), "json version mismatch") {
		t.Error("Want message about json version mismatch:", err)
	}
}

func TestStorer_InstanceBindingJSONVersionMismatching(t *testing.T) {
	t.Parallel()

	store := newTestStore()
	s := NewJSONStorer(store)
	var err error

	instance := Instance{Version: 9999}
	if err := s.Put(testInstanceID, instance, StoreNoLock); err != nil {
		t.Error(err)
	}

	_, err = getInstance(s, testInstanceID)
	if !isJSONVersionMismatchErr(err) {
		t.Error(err)
	}

	err = delBinding(s, testInstanceID, "binding-id")
	if !isJSONVersionMismatchErr(err) {
		t.Error(err)
	}

	key := []byte(strings.Repeat("0", 32))
	_, _, err = getBinding(s, key, testInstanceID, "binding-id")
	if !isJSONVersionMismatchErr(err) {
		t.Error(err)
	}

	err = putBinding(s, key, testInstanceID, "binding-id", BindingRequest{}, BindingResponse{})
	if !isJSONVersionMismatchErr(err) {
		t.Error(err)
	}
}

func TestStorer_JobVersionMismatching(t *testing.T) {
	t.Parallel()

	var err error
	store := newTestStore()
	s := NewJSONStorer(store)

	store.keyValues[mkJobKey(testInstanceID)] = `{"version": 9999}`

	_, err = getJob(s, testInstanceID)
	if !isJSONVersionMismatchErr(err) {
		t.Error(err)
	}
	err = updateJob(s, jobResult{InstanceID: testInstanceID})
	if !isJSONVersionMismatchErr(err) {
		t.Error(err)
	}
}

func isJSONVersionMismatchErr(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(errJSONVersionMismatch)
	return ok
}
