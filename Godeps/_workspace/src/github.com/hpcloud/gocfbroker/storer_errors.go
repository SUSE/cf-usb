package gocfbroker

import (
	"errors"
	"fmt"
)

// ErrStaleData is returned when the lock value from Get is not consistent
// with the value in the data store during a put.
var ErrStaleData = errors.New("data is stale; updated since last get")

// ErrKeyNotExist is returned when the key could not be found. The value of
// ErrKeyNotExist should be the key given.
type ErrKeyNotExist string

// Error implements error interface for ErrKeyNotExist
func (e ErrKeyNotExist) Error() string {
	return fmt.Sprintf("key not found: %s", string(e))
}

// IsKeyNotExist checks to see if the given error is an instance of the
// ErrKeyNotExist error type.
func IsKeyNotExist(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(ErrKeyNotExist)
	return ok
}

// errJSONVersionMismatch is returned when the json version of the object stored
// does not match the version of the broker
type errJSONVersionMismatch int

// Error returns human readable string.
func (e errJSONVersionMismatch) Error() string {
	return fmt.Sprintf(`json version mismatch, stored version: %d, broker version: %d`, int(e), brokerJSONSchemaVersion)
}
