package boltdb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hpcloud/gocfbroker"
)

func isNoLockValueErr(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(errNoLockValue)
	return ok
}

func isDecodeLockValueErr(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(errDecodeLockValue)
	return ok
}

var (
	testBucket   = []byte("testbucket")
	tmpDir       string
	tmpFileCount chan int
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	dir, err := ioutil.TempDir("", "gocfbroker_bolt_test")
	if err != nil {
		panic(err)
	}
	tmpDir = dir

	tmpFileCount = make(chan int)
	go func() {
		i := 0
		for {
			tmpFileCount <- i
			i++
		}
	}()

	code := m.Run()

	os.RemoveAll(tmpDir)
	os.Exit(code)
}

type testFileStore struct {
	gocfbroker.Storer
	filename string
}

func (t *testFileStore) Close() error {
	err := t.Storer.Close()
	os.Remove(t.filename)
	return err
}

func randFilename() string {
	return filepath.Join(tmpDir, fmt.Sprintf("gocfbroker_bolt_test_%d.bolt", <-tmpFileCount))
}

func openHelper() gocfbroker.Storer {
	tmpFile := randFilename()

	storer, err := New(tmpFile, string(testBucket))

	storer = &testFileStore{Storer: storer, filename: tmpFile}

	if err != nil {
		panic("could not open bolt db: " + err.Error())
	}

	return storer
}

func closeHelper(s gocfbroker.Storer) {
	if err := s.Close(); err != nil {
		panic("could not close bolt db: " + err.Error())
	}
}

func TestOpen(t *testing.T) {
	t.Parallel()

	// Test that the interfaces are matched.
	var storer gocfbroker.Storer
	var err error

	tmpFile := randFilename()
	defer os.Remove(tmpFile)

	storer, err = New(tmpFile, string(testBucket))

	if err != nil {
		t.Error(err)
	}
	storer.Close()
}

func TestPut(t *testing.T) {
	t.Parallel()

	bdb := openHelper()
	defer closeHelper(bdb)

	if err := bdb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	fileStorer := bdb.(*testFileStore)
	actualDB := fileStorer.Storer.(*boltDB)
	actualDB.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(testBucket)
		if got := string(bucket.Get([]byte("key"))); got != "value" {
			t.Error("Wrong value:", got)
		}
		return nil
	})
}

func TestPutErrors(t *testing.T) {
	t.Parallel()

	bdb := openHelper()
	defer closeHelper(bdb)

	ensureMissing := func(key string) error {
		if _, _, err := bdb.Get(key); !gocfbroker.IsKeyNotExist(err) {
			return err
		}
		return nil
	}

	if err := bdb.Put("key", "value", 10); !isNoLockValueErr(err) {
		t.Error("Wrong error type:", err)
	}

	if err := ensureMissing("key"); err != nil {
		t.Error(err)
	}

	actualDB := bdb.(*testFileStore).Storer.(*boltDB)
	err := actualDB.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(testBucket)
		return bucket.Put([]byte(mkLockKey("key")), []byte("notAnInteger"))
	})
	if err != nil {
		t.Error(err)
	}

	if err := bdb.Put("key", "newValue", 10); !isDecodeLockValueErr(err) {
		t.Error(err)
	}

	if err := ensureMissing("key"); err != nil {
		t.Error(err)
	}

	err = bdb.Put("key_lock", "value", 10)
	if _, ok := err.(errKeyHasLockSuffix); !ok {
		t.Errorf("Expected a has lock suffix error: (%T) %v", err, err)
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	bdb := openHelper()
	defer closeHelper(bdb)

	if err := bdb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	if val, _, err := bdb.Get("key"); err != nil {
		t.Error(err)
	} else if val != "value" {
		t.Error("Wrong value:", val)
	}

	if _, _, err := bdb.Get("cake"); !gocfbroker.IsKeyNotExist(err) {
		t.Error(err)
	}
}

func TestGetErrors(t *testing.T) {
	t.Parallel()

	bdb := openHelper()
	defer closeHelper(bdb)

	if err := bdb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	actualDB := bdb.(*testFileStore).Storer.(*boltDB)
	err := actualDB.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(testBucket)
		return bucket.Delete([]byte(mkLockKey("key")))
	})
	if err != nil {
		t.Error(err)
	}

	if _, _, err := bdb.Get("key"); !isNoLockValueErr(err) {
		t.Error(err)
	}

	err = actualDB.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(testBucket)
		return bucket.Put([]byte(mkLockKey("key")), []byte("notAnInteger"))
	})
	if err != nil {
		t.Error(err)
	}

	if _, _, err := bdb.Get("key"); !isDecodeLockValueErr(err) {
		t.Error(err)
	}
}

func TestCompareAndSwap(t *testing.T) {
	t.Parallel()

	bdb := openHelper()
	defer closeHelper(bdb)

	if err := bdb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	_, lock, err := bdb.Get("key")
	if err != nil {
		t.Error(err)
	}

	if err := bdb.Put("key", "value2", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	if err := bdb.Put("key", "value3", lock); err != gocfbroker.ErrStaleData {
		t.Error("Expected a stale data error back:", err)
	}

	var val string
	val, lock, err = bdb.Get("key")
	if err != nil {
		t.Error(err)
	}
	if val != "value2" {
		t.Error("Wrong value:", val, "it should have not been updated by second put")
	}
}

func TestDel(t *testing.T) {
	t.Parallel()

	bdb := openHelper()
	defer closeHelper(bdb)

	if err := bdb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	if err := bdb.Del("key"); err != nil {
		t.Error(err)
	}

	if _, _, err := bdb.Get("key"); !gocfbroker.IsKeyNotExist(err) {
		t.Error(err)
	}

	if err := bdb.Del("key"); !gocfbroker.IsKeyNotExist(err) {
		t.Error(err)
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()

	bdb := openHelper()
	defer closeHelper(bdb)

	tests := map[string]bool{
		"a_suffix":    true,
		"b_suffix":    true,
		"c_notsuffix": false,
	}

	for key := range tests {
		if err := bdb.Put(key, "value", gocfbroker.StoreNoLock); err != nil {
			t.Error(err)
		}
	}

	keys, err := bdb.Keys("_suffix")
	if err != nil {
		t.Error(err)
	}

	for key, shouldFind := range tests {
		found := false
		for _, fromDBKey := range keys {
			if key == fromDBKey {
				if found {
					t.Error("Found key more than once:", key)
				}
				found = true
			}
		}

		if shouldFind && !found {
			t.Error("Did not find key", key, "but should have found it.")
		} else if !shouldFind && found {
			t.Error("Found key", key, "but should not have found it.")
		}
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		Err error
		Str string
	}{
		{Err: errNoLockValue("key"), Str: "lock value not set for key: key"},
		{Err: errDecodeLockValue("key"), Str: "failed to decode lock value from key: key_lock"},
		{Err: errKeyHasLockSuffix("key_lock"), Str: `key has "_lock" suffix and may conflict with real lock values: key_lock`},
	}

	for i, test := range tests {
		if e := test.Err.Error(); e != test.Str {
			t.Errorf("%d) Want: %q got: %q", i, test.Str, e)
		}
	}
}
