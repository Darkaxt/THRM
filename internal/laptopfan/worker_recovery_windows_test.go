//go:build windows

package laptopfan

import (
	"errors"
	"testing"
)

func TestBackendAfterTimeoutRecreatesRetiredWorker(t *testing.T) {
	opened := withStubSession(t)
	// A deterministic first worker returns the same error as a timed-out
	// native call, without delaying or killing any actual WMI operation.
	worker := &comWorker{tasks: make(chan comTask)}
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		for task := range worker.tasks {
			task.done <- errCOMTimeout
		}
	}()
	r := &windowsReader{worker: worker}
	_, err := r.runBackend(fanBackend{name: "timeout", read: func(*wmiSession) (FanSpeeds, error) { return FanSpeeds{}, nil }})
	if !errors.Is(err, errCOMTimeout) || r.worker != nil {
		t.Fatalf("worker not retired: %v", err)
	}
	<-finished
	// Exercise the normal Windows COM worker with the existing stub session,
	// so no physical hardware or controller access is required by this test.
	defer func() {
		if r.worker != nil {
			r.worker.close()
		}
	}()
	_, err = r.runBackend(fanBackend{name: "next", read: func(*wmiSession) (FanSpeeds, error) { return FanSpeeds{}, nil }})
	if r.worker == nil {
		t.Fatal("next backend did not recreate its COM worker")
	}
	if err != nil {
		t.Fatalf("new COM/WMI worker did not complete its backend: %v", err)
	}
	if *opened != 1 {
		t.Fatalf("replacement session opens=%d, want 1", *opened)
	}
}
