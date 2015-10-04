package etcddb

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/hpcloud/gocfbroker"
)

var (
	testEtcdDir = "broker_etcd_test_dir"
	tmpDir      string

	etcdProcesses   = make(chan *exec.Cmd)
	etcdProcessList []*exec.Cmd
)

// Test that the interfaces are implemented.
var (
	_ gocfbroker.Storer = &etcdDB{}
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	dir, err := ioutil.TempDir("", "gocfbroker_etcd_test")
	if err != nil {
		panic(err)
	}
	tmpDir = dir

	code := m.Run()
	_ = os.RemoveAll(tmpDir)
	os.Exit(code)
}

type testEtcdStore struct {
	gocfbroker.Storer
	port          int
	runningDaemon *exec.Cmd
}

func (e *testEtcdStore) Close() error {
	e.Storer.Close()

	if err := e.runningDaemon.Process.Kill(); err != nil {
		msg := fmt.Sprintf("failed to kill running daemon: %s %v %d",
			e.runningDaemon.Path, e.runningDaemon.Args, e.port)
		panic(msg)
	}
	return nil
}

func tmpPort() int {
	listen, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP: net.IPv4(127, 0, 0, 1), // Leave port as 0
	})
	if err != nil {
		panic(err)
	}

	port := listen.Addr().(*net.TCPAddr).Port

	if err = listen.Close(); err != nil {
		panic(err)
	}

	return port
}

func openHelper() gocfbroker.Storer {
	tmpClientPort := tmpPort()
	tmpPeerPort := tmpPort()

	cmd := exec.Command("etcd",
		"--force-new-cluster",
		"--listen-peer-urls",
		fmt.Sprintf("http://localhost:%d", tmpPeerPort),
		"--listen-client-urls",
		fmt.Sprintf("http://localhost:%d", tmpClientPort),
		"--data-dir",
		fmt.Sprintf("testetcd%d.etcd", tmpClientPort),
	)
	cmd.Dir = tmpDir
	if err := cmd.Start(); err != nil {
		panic("failed to start etcd db: " + err.Error())
	}

	// Wait for ETCD to come up (takes some time)
	cli := etcd.NewClient([]string{
		fmt.Sprintf("http://127.0.0.1:%d", tmpClientPort),
	})
	for i := 0; i < 100; i++ {
		_, err := cli.Get("ping", false, false)
		if err == nil {
			break
		}

		if etcdErr, ok := err.(*etcd.EtcdError); ok {
			if etcdErr.ErrorCode == errEtcdKeyNotFound {
				break
			}
		}
		<-time.After(time.Millisecond * 10)
	}

	etcdDBStorer, err := New(testEtcdDir, fmt.Sprintf("http://127.0.0.1:%d", tmpClientPort))
	if err != nil {
		panic("could not open create storer: " + err.Error())
	}

	return &testEtcdStore{Storer: etcdDBStorer, port: tmpClientPort, runningDaemon: cmd}
}

func closeHelper(t *testing.T, s gocfbroker.Storer) {
	if r := recover(); r != nil {
		t.Fail()
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		t.Logf("panic: %v\n%s\n", r, buf[:n])
	}

	if err := s.Close(); err != nil {
		panic("could not close etcd db: " + err.Error())
	}
}

func TestOpen(t *testing.T) {
	t.Parallel()

	storer := openHelper()
	defer closeHelper(t, storer)
}

func TestPut(t *testing.T) {
	t.Parallel()

	edb := openHelper()
	defer closeHelper(t, edb)

	if err := edb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	etcdStorer := edb.(*testEtcdStore)
	etcddb := etcdStorer.Storer.(*etcdDB)
	resp, err := etcddb.client.Get(etcddb.mkKey("key"), false, false)
	if err != nil {
		t.Fatal("failed to get the key:", err)
	}
	if resp.Node.Value != "value" {
		t.Error("wrong value:", resp.Node.Value)
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	edb := openHelper()
	defer closeHelper(t, edb)

	if err := edb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	if val, _, err := edb.Get("key"); err != nil {
		t.Error(err)
	} else if val != "value" {
		t.Error("Wrong value:", val)
	}

	if _, _, err := edb.Get("cake"); err == nil {
		t.Error("Non-existent key should produce ErrKeyNotExist")
	} else if _, ok := err.(gocfbroker.ErrKeyNotExist); !ok {
		t.Error(err)
	}
}

func TestCompareAndSwap(t *testing.T) {
	t.Parallel()

	edb := openHelper()
	defer closeHelper(t, edb)

	if err := edb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	_, lock, err := edb.Get("key")
	if err != nil {
		t.Error(err)
	}

	if err := edb.Put("key", "value2", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	if err := edb.Put("key", "value3", lock); err != gocfbroker.ErrStaleData {
		t.Error("Expected a stale data error back:", err)
	}

	var val string
	val, lock, err = edb.Get("key")
	if err != nil {
		t.Error(err)
	}
	if val != "value2" {
		t.Error("Wrong value:", val, "it should have not been updated by second put")
	}
}

func TestDel(t *testing.T) {
	t.Parallel()

	edb := openHelper()
	defer closeHelper(t, edb)

	if err := edb.Put("key", "value", gocfbroker.StoreNoLock); err != nil {
		t.Error(err)
	}

	if err := edb.Del("key"); err != nil {
		t.Error(err)
	}

	if _, _, err := edb.Get("key"); err == nil {
		t.Error("Non-existent key should produce ErrKeyNotExist")
	} else if _, ok := err.(gocfbroker.ErrKeyNotExist); !ok {
		t.Error(err)
	}

	if err := edb.Del("key"); err == nil {
		t.Error("Non-existent key should produce ErrKeyNotExist")
	} else if _, ok := err.(gocfbroker.ErrKeyNotExist); !ok {
		t.Error(err)
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()

	edb := openHelper()
	defer closeHelper(t, edb)

	tests := map[string]bool{
		"a_suffix":    true,
		"b_suffix":    true,
		"c_notsuffix": false,
	}

	for key := range tests {
		if err := edb.Put(key, "value", gocfbroker.StoreNoLock); err != nil {
			t.Error(err)
		}
	}

	keys, err := edb.Keys("_suffix")
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
