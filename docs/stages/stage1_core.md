# Stage 1 — Core: config, options, constants

### Goal
Establish the base structures of the package: configuration, functional options, and constants.

### Steps

- [ ] Create `tinygo.go` with:
  - Constant `DefaultVersion = "0.40.1"`
  - Constant `DefaultInstallDir = ".tinywasm"` (relative to home)
  - Struct `config` (private): `version string`, `installDir string`, `logger func(string)`
  - Function `newConfig(opts ...Option) *config` that applies defaults + options
  - Type `Option func(*config)`
  - Functions: `WithVersion(v string) Option`, `WithInstallDir(dir string) Option`, `WithLogger(f func(string)) Option`
  - Helper function `(c *config) binPath() string` that returns the full path to the tinygo binary according to the OS:
    - Linux/macOS: `{installDir}/tinygo/bin/tinygo`
    - Windows: `{installDir}/tinygo/bin/tinygo.exe`
  - Helper function `(c *config) downloadURL() string` that builds the download URL based on `runtime.GOOS` and `runtime.GOARCH`:
    - Linux/macOS: `.tar.gz`
    - Windows: `.zip`

  - Private fields for test injection: `lookPath func(string) (string, error)` (default: `exec.LookPath`), `httpClient *http.Client` (default: `http.DefaultClient`)
  - Internal options (not exported, for tests only): `withLookPath(f) Option`, `withHTTPClient(c) Option`

### Validation via `tinygo_test.go`
- `newConfig()` without options returns correct defaults (version, installDir, lookPath, httpClient).
- `WithVersion("0.39.0")` overrides the version.
- `binPath()` returns `.exe` on Windows, and no extension on Linux/Darwin.
- `downloadURL()` generates valid URLs for all 6 combinations (Linux/Darwin/Windows x amd64/arm64).
- `withLookPath` and `withHTTPClient` override defaults correctly.

### Files
- [tinygo.go](../../tinygo.go)
- [tinygo_test.go](../../tinygo_test.go)
