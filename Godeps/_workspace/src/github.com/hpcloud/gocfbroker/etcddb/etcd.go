package etcddb

import (
	"path"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/hpcloud/gocfbroker"
)

const (
	errEtcdKeyNotFound   = 100
	errEtcdCompareFailed = 101
)

// New creates an Etcd client that will connect to the machines supplied
// and prefix all keys with the given directory name to avoid collisions.
func New(directory string, machines ...string) (gocfbroker.Storer, error) {
	e := &etcdDB{
		machines:  machines,
		directory: directory,
	}

	e.client = etcd.NewClient(machines)

	// Create initial directory
	_, err := e.client.CreateDir(directory, 0)
	if err != nil {
		return nil, err
	}

	return e, nil
}

type etcdDB struct {
	machines  []string
	directory string

	client *etcd.Client
}

// Put new value for key, optional lock value (use gocfbroker.StoreNoLock for
// no locking mechanism.
func (e *etcdDB) Put(key, val string, lock int) error {
	var err error

	key = e.mkKey(key)
	if lock != gocfbroker.StoreNoLock {
		_, err = e.client.CompareAndSwap(key, val, 0, "", uint64(lock))
	} else {
		_, err = e.client.Set(key, val, 0)
	}

	if etcdErr, ok := err.(*etcd.EtcdError); ok && etcdErr.ErrorCode == errEtcdCompareFailed {
		return gocfbroker.ErrStaleData
	}

	return err
}

// Get a key's current value and lock value.
func (e *etcdDB) Get(key string) (val string, lock int, err error) {
	var resp *etcd.Response

	resp, err = e.client.Get(e.mkKey(key), false, false)
	if err == nil {
		val = resp.Node.Value
		lock = int(resp.Node.ModifiedIndex)
		return val, lock, nil
	}

	if etcdErr, ok := err.(*etcd.EtcdError); ok && etcdErr.ErrorCode == errEtcdKeyNotFound {
		return val, lock, gocfbroker.ErrKeyNotExist(key)
	}

	return val, lock, err
}

// Del a key
func (e *etcdDB) Del(key string) error {
	_, err := e.client.Delete(e.mkKey(key), false)

	if etcdErr, ok := err.(*etcd.EtcdError); ok && etcdErr.ErrorCode == errEtcdKeyNotFound {
		return gocfbroker.ErrKeyNotExist(key)
	}

	return err
}

// Keys looks up the keys if they have the specific suffix (but also limited
// to the directory)
func (e *etcdDB) Keys(suffix string) ([]string, error) {
	var keys []string

	resp, err := e.client.Get(e.directory, false, false)
	if err != nil {
		return nil, err
	}

	for _, node := range resp.Node.Nodes {
		if strings.HasSuffix(node.Key, suffix) {
			keys = append(keys, path.Base(node.Key))
		}
	}

	return keys, nil
}

// Close the etcd client (no-op since etcd client has no persistent connection)
func (e *etcdDB) Close() error {
	return nil
}

func (e etcdDB) mkKey(key string) string {
	return path.Join(e.directory, key)
}
