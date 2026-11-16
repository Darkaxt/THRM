package device

import "testing"

func TestRepeatedLightProgramDoesNotUploadOrCommitAgain(t *testing.T) {
	program := newLightProgram(0, lightSpeedMedium, 100)
	m := NewManager(testLogger{})
	m.lastLightProgram, m.hasLightProgram = program, true
	m.hasRGBEnabled, m.rgbEnabled = true, true
	if err := m.applyLightFramesLocked(program); err != nil {
		t.Fatal(err)
	}
	if m.linkHealth.txCount != 0 || !m.lastFlashWriteAt.IsZero() {
		t.Fatal("same-program replay wrote to device")
	}
	m.resetLightStateCacheLocked()
	if m.hasLightProgram {
		t.Fatal("reconnect must invalidate program cache")
	}
}

func TestSmartConfigurationReplayPreservesCurrentPreset(t *testing.T) {
	m := NewManager(testLogger{})
	m.lightConfig.Mode, m.hasLightConfig = "smart_temp", true
	m.smartLightPreset, m.hasSmartLightPreset, m.smartLightBandIndex = 3, true, 2
	if err := m.setLightSmartTempLocked(m.lightConfig); err != nil {
		t.Fatal(err)
	}
	if m.linkHealth.txCount != 0 || m.smartLightPreset != 3 || m.smartLightBandIndex != 2 {
		t.Fatal("unchanged configuration restarted animation")
	}
}
