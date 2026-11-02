# Windows GUI reopening repair

Authorized by the user on 2026-10-24 as part of the Taskbar Thermals / BS2 Pro integration. The authoritative integration specification is maintained in the sibling taskbar-thermals checkout at docs/CUSTOMIZATION-SPEC.md, requirements R6/R8/R11.

THRM 3.7.3's core uses the presence of any IPC client to decide that a GUI exists. A persistent telemetry observer therefore prevents a closed GUI from being relaunched. The GUI already provides a Wails single-instance lock whose second-instance callback restores and unminimises the existing window.

Repair: route every explicit show-window request through that existing GUI executable/single-instance path. The GUI owns its own lifetime and restoration; telemetry connections no longer stand in for GUI presence. Do not add polling, guessed liveness deadlines, fan-control writes or duplicate persistent GUI instances. Preserve all user settings and the existing telemetry protocol.

Acceptance: reproduce the observer-only failure before the fix; pass focused core/IPC regression checks; verify the real tray show-window action after closing the GUI with Taskbar Thermals connected, repeat after closing again, and restore an already-open/minimised GUI without creating a second persistent GUI. Verify live BS2 Pro monitoring/control and user settings after activation. Publish and independently download a complete Windows portable repair package with component provenance and install the verified package. The unchanged GUI/bridge remain the verified upstream 3.7.3 binaries; the patched core has version 3.7.3.1. This is a complete core repair, with no controller/Vader input claim.

No unrelated upstream changes are included: the repair branch starts at tag v3.7.3 / commit 2131102526c8772211896ce23dd24a93065a37d3.

## Build and release

On Windows x64 with Go 1.27.3, a C compiler and go-winres v0.3.3, run `go test ./internal/coreapp ./internal/ipc -count=1`, then `scripts/package-core-repair.ps1 -OutputDirectory <new-folder> -WinresPath <go-winres.exe>`. The script verifies the upstream archive before reusing its unchanged GUI and temperature bridge, builds the repaired core and records component hashes and source revisions. Publish the verified ZIP and SHA256SUMS.txt as the normal `repair-v3.7.3.1` release after real tray reopening/restore and live telemetry verification. This core-only repair does not use the upstream all-platform release workflow.

## Install and recovery

Extract the complete portable package into its own version directory. Quit the old core normally before launching the new `THRM Core.exe`; keep the existing user profile. If Windows startup is enabled, update the existing THRM task's executable to that core while preserving its arguments, triggers, principal and settings. The GUI and its About version remain 3.7.3; the core's Windows file version is 3.7.3.1.

To roll back, quit the repaired core, launch the retained original 3.7.3 core and restore the existing startup action's executable. User settings remain in the same profile. Upstream updates can replace this repair; recheck the reopening fix before updating. No new updater or device-control protocol is introduced.

## Weekly Windows maintenance, 2026-11-02

Authorization: existing weekly THRM/Taskbar Thermals upstream integration and installed-deployment automation. Preserve the current Windows 3.7.3 GUI/bridge and observer-aware reopening repair. Upstream main is 67091d5dea620f06585e5793b1f1eab68c96746e (16 commits beyond the integrated 2131102 baseline); the latest official normal release remains v3.7.3. Apply compatible Windows core maintenance from efa8176 (accurate caller logging), 172568b (one blocking native hotkey message loop) and ac3acb7 (avoid unnecessary tray icon redraws and repeated debug logs). Assess conflict resolution against the original GUI repair, hotkey lifecycle, local observer contract and unchanged device-control/config behavior.

The unreleased 3.8 GUI/Wails3/firmware/driver feature migration is not part of this core-maintenance package: it requires a different GUI/bridge pairing and includes new firmware capabilities outside this maintenance outcome. Preserve the verified current pairing rather than mixing a 3.8 core with 3.7.3 GUI. BS1/Linux-only changes are inapplicable to this Windows BS2 Pro installation. General frontend/dependency/CI migrations accompany the unreleased GUI rather than this compatible backport. These are applicability decisions, not required work deferred from this maintenance package. No firmware, driver, security, voltage, WSL or profile changes.

| Stage | Status | Acceptance |
|---|---|---|
| M1: Compatible Windows core integration | ACTIVE | Integrate the three identified upstream changes without losing reopening; verify affected native/core/logger/tray/IPC contracts, race/lifecycle where applicable, and build a complete 3.7.3.2 portable package reusing verified 3.7.3 GUI/bridge. No blockers or tracked deferrals currently. |
| M2: Normal release and installed delivery | NOT STARTED | After M1, publish repair-v3.7.3.2 with ZIP/checksums, independently verify source/version/component hashes, update existing THRM startup and installation once through the broker, verify actual core/observer telemetry and unchanged settings, retain 3.7.3.1 rollback, and clean disposable output. |

Only M1 is ACTIVE. Independent Taskbar Thermals upstream review follows THRM closure or justified parking. Required source/package verification must precede publication; any unavailable administrator approval or actual hardware acceptance remains an explicit delivery blocker, with the working installed core preserved. The separate Vader input check is unchanged and unverified.
