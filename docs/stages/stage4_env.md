# Stage 4 — Env: GetRoot, GetEnv

### Goal

Expose the environment configuration required for a locally installed TinyGo to work,
without modifying the user's shell or global process environment.

### Dependency
- Stage 3 completed (detect: IsInstalled, GetPath, GetVersion, EnsureInstalled)

### Background

When TinyGo is installed into a custom directory, the `TINYGOROOT` variable must point
to `{installDir}/tinygo/` so TinyGo can locate its standard library, targets, and
compiler support files. Consumers (e.g. `gobuild`) must inject this into any
`exec.Command` that invokes tinygo.

### Steps

- [ ] Create `env.go` with:
  - `GetRoot(opts ...Option) string`
    - Builds config from opts.
    - If tinygo is found in PATH (`lookPath` succeeds) → returns `""` (system install, user controls env).
    - Otherwise → returns `filepath.Join(c.installDir, "tinygo")`.
  - `GetEnv(opts ...Option) []string`
    - Calls `GetRoot()`.
    - If root is `""` → returns `os.Environ()` unchanged.
    - Otherwise:
      - Copies `os.Environ()`.
      - Replaces existing `TINYGOROOT=...` entry if present, or appends new one.
      - Prepends `{root}/bin` to the `PATH` entry.
      - Returns the modified slice.

- [ ] Create `env_test.go` with:
  - `TestGetRoot_SystemInstall` — mock `withLookPath` to succeed → assert returns `""`.
  - `TestGetRoot_LocalInstall` — mock `withLookPath` to fail, set `WithInstallDir(tmpDir)` → assert returns `filepath.Join(tmpDir, "tinygo")`.
  - `TestGetEnv_SystemInstall` — mock lookPath to succeed → assert returned env equals `os.Environ()`.
  - `TestGetEnv_LocalInstall` — mock lookPath to fail → assert:
    - Slice contains `TINYGOROOT={root}`.
    - `PATH` entry contains `{root}/bin` prepended.
    - Does NOT mutate `os.Environ()` (copy check).
  - `TestGetEnv_OverridesExistingTINYGOROOT` — set env with stale `TINYGOROOT`, assert it's replaced correctly.

### Files
- [env.go](../../env.go)
- [env_test.go](../../env_test.go)

### Consumer Example

```go
// In gobuild or client:
binPath, _ := tinygo.EnsureInstalled()

cmd := exec.Command(binPath, "build", "-o", "out.wasm", ".")
cmd.Env = tinygo.GetEnv()  // TINYGOROOT + PATH injected here
out, err := cmd.CombinedOutput()
```
