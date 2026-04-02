# Stage 3 — Detect: IsInstalled, GetPath, GetVersion, EnsureInstalled

### Goal
Implement TinyGo binary detection functions and the main orchestrator `EnsureInstalled`.

### Dependency
- Stage 2 completed (functional Install)

### Steps

- [ ] Create `detect.go` with:
  - `IsInstalled() bool` — returns true if tinygo is available (PATH or local).
    - First looks in `PATH` with `exec.LookPath("tinygo")`.
    - Then looks in local path `~/.tinywasm/tinygo/bin/tinygo[.exe]`.
  - `GetPath() (string, error)` — returns the path to the tinygo binary.
    - Same logic: PATH first, then local.
    - Error if not found in either location.
  - `GetVersion() (string, error)` — runs `tinygo version` and returns the output.
    - Uses `GetPath()` to find the binary.
    - Executes the command and parses output.

- [ ] Create `ensure.go` with:
  - `EnsureInstalled(opts ...Option) (string, error)` — main function.
    - If `GetPath()` finds the binary, returns the path directly.
    - If not, calls `Install(opts...)` and then `GetPath()`.
    - Returns the path to the binary ready for use.

- [ ] Create `detect_test.go` with:
  - `IsInstalled` returns true when `withLookPath` finds tinygo in PATH.
  - `IsInstalled` returns true when binary exists at local `binPath()` (create fake binary in temp dir).
  - `IsInstalled` returns false when both PATH and local are missing.
  - `GetPath` returns PATH result with priority over local.
  - `GetPath` falls back to local path when PATH misses.
  - `GetPath` returns error when not found anywhere.
  - `GetVersion` executes real `tinygo version` if available (skip if not installed).

- [ ] Create `ensure_test.go` with:
  - `EnsureInstalled` returns existing path without calling Install (mock lookPath to succeed).
  - `EnsureInstalled` triggers Install when not found (use `httptest.Server` + `withLookPath` that fails).

### Files
- [detect.go](../../detect.go)
- [ensure.go](../../ensure.go)
- [detect_test.go](../../detect_test.go)
