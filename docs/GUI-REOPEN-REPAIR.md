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
