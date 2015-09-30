package gocfbroker

import (
	"sync"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	saveRandInt := backoffRandInt
	saveAfter := backoffAfter

	// Take the random out, max is exclusive
	backoffRandInt = func(max int) int {
		return max - 1
	}

	// Take the sleep out
	durChan := make(chan time.Duration)
	backoffAfter = func(dur time.Duration) <-chan time.Time {
		after := make(chan time.Time)
		go func() {
			durChan <- dur
			after <- time.Time{}
		}()
		return after
	}

	// Set everything back to normal
	defer func() {
		backoffRandInt = saveRandInt
		backoffAfter = saveAfter
	}()

	// Until the last attempt, keep retrying
	attempts := 1
	failyFunction := func() error {
		if attempts > 10 {
			return nil
		}
		attempts++
		return errBackoffRetry
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Use a goroutine to siphon the sleep values
	go func() {
		wantSleep := []int64{
			20, 40, 80, 160, 320, 640, 1280, 2560,
			backoffMaxSleep, backoffMaxSleep,
		}

		for i := 0; i < 10; i++ {
			if sleep := int64(<-durChan / time.Millisecond); sleep != wantSleep[i] {
				t.Errorf("Attempt: %d, Want sleep: %d, Got sleep: %d", i+1, wantSleep[i], sleep)
			}
		}

		wg.Done()
	}()

	if err := exponentialJitter(failyFunction); err != nil {
		t.Error(err)
	}

	wg.Wait()
}
