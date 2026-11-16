package coreapp

import (
	"testing"

	"github.com/TIANLI0/THRM/internal/types"
)

func TestFanCallbackCannotBlockHIDReaderOnCoreLock(t *testing.T) {
	a := &CoreApp{}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	// This callback executes inline on the sole HID reader. A held Core lock
	// must leave its cached status untouched and immediately return to ACK I/O.
	a.onFanDataUpdate(&types.FanData{CurrentRPM: 2000})
}
