package coreapp

import (
	"fmt"
	"os"
	"testing"

	"github.com/TIANLI0/THRM/internal/ipc"
)

func TestShowWindowRequestLaunchesGUIWithTelemetryObserver(t *testing.T) {
	t.Setenv(ipc.EnvPipeName, fmt.Sprintf("THRM-reopen-test-%d", os.Getpid()))
	server := ipc.NewServer(func(ipc.Request) ipc.Response { return ipc.Response{Success: true} }, nil)
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	observer := ipc.NewClient(nil)
	if err := observer.Connect(); err != nil {
		t.Fatal(err)
	}
	defer observer.Close()
	if _, err := observer.SendRequest(ipc.ReqGetDeviceStatus, nil); err != nil {
		t.Fatal(err)
	}
	if !server.HasClients() {
		t.Fatal("observer fixture is not connected")
	}
	app := &CoreApp{ipcServer: server}
	launches := 0
	app.showWindowRequest(func() error { launches++; return nil })
	if launches != 1 {
		t.Fatalf("GUI launches = %d, want 1 with observer connected", launches)
	}
	if _, err := observer.SendRequest(ipc.ReqGetDeviceStatus, nil); err != nil {
		t.Fatalf("observer lost after reopen: %v", err)
	}
}
