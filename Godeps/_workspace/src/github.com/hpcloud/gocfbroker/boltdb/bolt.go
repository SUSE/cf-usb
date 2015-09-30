package boltdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/hpcloud/gocfbroker"
)

const (
	suffixLock = "_lock"
)

type (
	errNoLockValue      string
	errDecodeLockValue  string
	errKeyHasLockSuffix string
)

func (e errNoLockValue) Error() string {
	return fmt.Sprintf("lock value not set for key: %s", string(e))
}

func (e errDecodeLockValue) Error() string {
	return fmt.Sprintf("failed to decode lock value from key: %s%s", string(e), suffixLock)
}

func (e errKeyHasLockSuffix) Error() string {
	return fmt.Sprintf("key has \"_lock\" suffix and may conflict with real lock values: %s", string(e))
}

// New creates a BoltDB filename with a specific bucket to avoid
// name collisions.
func New(filename, bucket string) (gocfbroker.Storer, error) {
	var err error
	b := &boltDB{
		filename: filename,
		bucket:   bucket,
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err = os.MkdirAll(dir, 0775); err != nil {
		return nil, err
	}

	b.db, err = bolt.Open(b.filename, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte(bucket))
		return e
	})

	return b, err
}

type boltDB struct {
	db       *bolt.DB
	filename string
	bucket   string
}

// Put key-value. In addition to the Storer implementation details, Put() will
// return errKeyHasLockSuffix when passed a key with _lock as a suffix in order
// to preserve the locking mechanism integrity.
func (b *boltDB) Put(key, val string, lock int) error {
	if strings.HasSuffix(key, suffixLock) {
		return errKeyHasLockSuffix(key)
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		var err error
		var gotLock int
		if lockVal := bucket.Get([]byte(mkLockKey(key))); lockVal == nil {
			if lock != gocfbroker.StoreNoLock {
				return errNoLockValue(key)
			}
		} else if gotLock, err = strconv.Atoi(string(lockVal)); err != nil {
			return errDecodeLockValue(key)
		} else if lock != gocfbroker.StoreNoLock && gotLock != lock {
			return gocfbroker.ErrStaleData
		}

		if err := bucket.Put([]byte(mkLockKey(key)), []byte(strconv.Itoa(gotLock+1))); err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(val))
	})
}

// Get value for a key.
func (b *boltDB) Get(key string) (val string, lock int, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		value := bucket.Get([]byte(key))
		if value == nil {
			return gocfbroker.ErrKeyNotExist(key)
		}

		if lockVal := bucket.Get([]byte(mkLockKey(key))); lockVal == nil {
			return errNoLockValue(key)
		} else if lock, err = strconv.Atoi(string(lockVal)); err != nil {
			return errDecodeLockValue(key)
		}

		val = string(value)
		return nil
	})

	return val, lock, err
}

// Del removes a key-value.
func (b *boltDB) Del(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		if val := bucket.Get([]byte(key)); val == nil {
			return gocfbroker.ErrKeyNotExist(key)
		}

		if err := bucket.Delete([]byte(key)); err != nil {
			return err
		}
		return bucket.Delete([]byte(mkLockKey(key)))
	})
}

// Keys returns all the keys with the suffix.
func (b *boltDB) Keys(suffix string) ([]string, error) {
	var keys []string
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.bucket))

		c := bucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			if key := string(k); strings.HasSuffix(key, suffix) {
				keys = append(keys, key)
			}
		}

		return nil
	})

	return keys, err
}

// Close the boltdb instance.
func (b *boltDB) Close() error {
	return b.db.Close()
}

func mkLockKey(key string) string {
	return key + suffixLock
}
