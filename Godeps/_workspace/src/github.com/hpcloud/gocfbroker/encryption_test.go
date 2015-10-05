package gocfbroker

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncryption(t *testing.T) {
	t.Parallel()

	pt := []byte("hello world")
	key := []byte(strings.Repeat("a", 32))
	failKey := []byte(strings.Repeat("b", 32))
	ct, err := encrypt(key, pt)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(ct, pt) == 0 {
		t.Error("Not encrypted")
	}

	decryptPt, err := decrypt(key, string(ct))
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(decryptPt, pt) != 0 {
		t.Error("Decrypt did not yield correct pt:", decryptPt)
	}

	failPt, err := decrypt(failKey, string(ct))
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(failPt, pt) == 0 {
		t.Error("It should not decrypt correctly with wrong key")
	}
}
