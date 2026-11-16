package device

import (
	"errors"
	"io"
	"testing"
)

type reportWriter struct {
	count int
	err   error
}

func (w reportWriter) Write([]byte) (int, error) { return w.count, w.err }

func TestWindowsPaddedWriteIsCompleteButShortWriteFails(t *testing.T) {
	report := make([]byte, 25)
	for _, n := range []int{25, 32, 65} {
		if err := writeHIDReport(reportWriter{count: n}, report); err != nil {
			t.Fatalf("complete/padded count %d rejected: %v", n, err)
		}
	}
	if err := writeHIDReport(reportWriter{count: 24}, report); !errors.Is(err, io.ErrShortWrite) {
		t.Fatalf("short write not detected: %v", err)
	}
	native := errors.New("native write failed")
	if err := writeHIDReport(reportWriter{count: 65, err: native}, report); !errors.Is(err, native) {
		t.Fatalf("native error lost: %v", err)
	}
}
