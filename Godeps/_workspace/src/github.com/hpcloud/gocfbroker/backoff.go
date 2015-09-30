package gocfbroker

import (
	"errors"
	"math/rand"
	"time"
)

// Time constants in milliseconds
const (
	backoffMaxSleep = 3000
	backoffBase     = 10
)

var (
	errBackoffRetry = errors.New("retry operation")
	backoffRandInt  = rand.Intn
	backoffAfter    = time.After
)

// exponentialJitter implements exponential backoff with a jitter factor.
// If the error returned from fn is backoffRetryErr it will attempt the operation
// again. If not it will simply consider it a permanent failure and return the
// error returned from  fn.
func exponentialJitter(fn func() error) error {
	attempts := uint(0)

	for {
		err := fn()
		if err == nil {
			break
		} else if err != errBackoffRetry {
			return err
		}

		backoff := backoffMaxSleep
		if calc := backoffBase * (2 << attempts); calc < backoff {
			backoff = calc
		}

		sleep := time.Duration(backoffRandInt(backoff+1)) * time.Millisecond
		<-backoffAfter(sleep)

		attempts++
	}

	return nil
}
