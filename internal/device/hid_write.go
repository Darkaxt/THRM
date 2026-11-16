package device

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/TIANLI0/THRM/internal/deviceproto"
	"github.com/TIANLI0/THRM/internal/types"
	"github.com/sstallion/go-hid"
)

// Backported from upstream 7ae35fe/da0c69e without the unrelated Linux USB
// descriptor migration. Keep this branch's established Pro-model report format.
func writeHIDReport(writer interface{ Write([]byte) (int, error) }, report []byte) error {
	n, err := writer.Write(report)
	// Windows hidapi returns the descriptor-padded byte count. Only a smaller
	// count is a short write, including when the transport is Bluetooth HID.
	if err == nil && n < len(report) {
		return fmt.Errorf("HID 写入 %d/%d 字节: %w", n, len(report), io.ErrShortWrite)
	}
	return err
}

func (m *Manager) writeHIDReportLocked(frame []byte, reportLen int) ([]byte, error) {
	if m.hidWritePoisoned != nil {
		return nil, fmt.Errorf("HID 上一次写入仍可能挂起，等待更换句柄: %w", m.hidWritePoisoned)
	}
	report := deviceproto.BuildReport(frame, reportLen)
	m.recordDebugFrame("tx", types.DeviceTypeHID, report)
	m.linkHealth.record(m.connectionGen.Load(), true, report, time.Now())
	err := writeHIDReport(m.device, report)
	if isPendingHIDWriteFailure(err) {
		m.hidWritePoisoned = err
		if m.pendingHIDWrites == nil {
			m.pendingHIDWrites = make(map[*hid.Device][]byte)
		}
		m.pendingHIDWrites[m.device] = report
		m.realtimeWriteRecoveryScheduled = true
		generation := m.connectionGen.Load()
		m.logError("Windows HID 写入等待失败，停止复用挂起的 I/O 状态并恢复句柄: %v", err)
		go m.disconnectGeneration(true, false, generation)
	}
	return report, err
}

func isPendingHIDWriteFailure(err error) bool {
	return err != nil && strings.Contains(err.Error(), "hid_write/WaitForSingleObject")
}
