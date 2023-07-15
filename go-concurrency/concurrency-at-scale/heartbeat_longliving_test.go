package main

import (
	"testing"
	"time"
)

// go test heartbeat_longliving_test.go heartbeat_longliving.go
//

func TestDoWorkLongLiving(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	exp := []int{0, 1, 2, 3, 5}
	heartbeat, results := DoWorkLongLiving(done, exp...)

	<-heartbeat

	i := 0

	for r := range results {
		if expected := exp[i]; r != expected {
			t.Errorf("index %v: expected %v, got %v", i, expected, r)
		}

		i++
	}
}

func TestDoWorkLongLivingLabel(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	exp := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := DoWorkLongLivingLabel(done, timeout, exp...)

	<-heartbeat

	i := 0
	for {
		select {
		case r, ok := <-results:
			if !ok {
				return
			}

			if expected := exp[i]; expected != r {
				t.Errorf("index %v: expected %v, got %v", i, expected, r)
			}

			i++
		case <-heartbeat:
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
