//go:build windows

package powernotify

import (
	"reflect"
	"testing"
)

func TestPowerEventsDistinguishUserPresence(t *testing.T) {
	var events []string
	pn := &notifier{
		onSuspend: func() { events = append(events, "suspend") },
		onResume: func(userPresent bool) {
			if userPresent {
				events = append(events, "user-resume")
			} else {
				events = append(events, "automatic-resume")
			}
		},
	}
	for _, event := range []uint32{pbtAPMSuspend, pbtAPMResumeAutomatic, 0, pbtAPMResumeSuspend} {
		pn.dispatch(event)
	}
	if want := []string{"suspend", "automatic-resume", "user-resume"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("power event dispatch = %v, want %v", events, want)
	}
	// Native callbacks may legitimately be absent.
	for _, event := range []uint32{pbtAPMSuspend, pbtAPMResumeAutomatic, pbtAPMResumeSuspend} {
		(&notifier{}).dispatch(event)
	}
}

func TestNativePowerRegistrationAndIdempotentRelease(t *testing.T) {
	stop, err := RegisterSuspendResumeNotifications(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	stop()
	stop()
}
