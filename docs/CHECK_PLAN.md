# PLAN: tinygo package — Automated Cross-platform Installer

## Goal
Standalone package to manage automated TinyGo installation in the tinywasm ecosystem.
Full support for Linux, macOS, and Windows — no admin privileges, no user interaction.

## Design Decisions
- **No Interface**: direct public functions. Function injection via `config` for testability.
- **No Package Managers**: local tarball/zip for all 3 OS (sudo-free + unattended + cross-platform).
- **Hardcoded Version**: `0.40.1` (latest stable release).
- **Installation Path**: `~/.tinywasm/tinygo/`.
- **Formats**: `.tar.gz` (Linux/macOS), `.zip` (Windows).

## Diagrams
- [Install Flow](diagrams/install_flow.md) — EnsureInstalled → Install → download → extract → verify

## Public API

```go
package tinygo

func IsInstalled() bool
func GetPath() (string, error)
func GetVersion() (string, error)
func EnsureInstalled(opts ...Option) (string, error)
func Install(opts ...Option) error
```

### Options

```go
func WithVersion(v string) Option
func WithInstallDir(dir string) Option
func WithLogger(f func(string)) Option
```

## Testing Strategy

Function injection via private `config` fields — no interfaces, no network, no tinygo required.

| Function | Mock strategy |
|----------|---------------|
| `IsInstalled`, `GetPath` | `withLookPath` (internal option) + temp dir with fake binary |
| `GetVersion` | Real binary if available, skip otherwise |
| `Install` | `httptest.Server` serves fake archive + `WithInstallDir(tempDir)` |
| `EnsureInstalled` | Combines lookPath mock + httptest |
| `extractTarGz`, `extractZip` | Programmatically created test archives, Zip Slip protection |
| `downloadURL`, `binPath` | Direct unit test on config |

## Release URLs Reference

```
Linux amd64:   tinygo{V}.linux-amd64.tar.gz
Linux arm64:   tinygo{V}.linux-arm64.tar.gz
macOS amd64:   tinygo{V}.darwin-amd64.tar.gz
macOS arm64:   tinygo{V}.darwin-arm64.tar.gz
Windows amd64: tinygo{V}.windows-amd64.zip
```
Base: `https://github.com/tinygo-org/tinygo/releases/download/v{V}/`

## Target File Structure

- [tinygo.go](../tinygo.go) — config, options, constants
- [download.go](../download.go) — HTTP download
- [extract.go](../extract.go) — untar (Linux/macOS) + unzip (Windows)
- [install.go](../install.go) — Install orchestrator
- [detect.go](../detect.go) — IsInstalled, GetPath, GetVersion
- [ensure.go](../ensure.go) — EnsureInstalled
- [tinygo_test.go](../tinygo_test.go) — config, options, binPath, downloadURL
- [extract_test.go](../extract_test.go) — extraction + Zip Slip
- [install_test.go](../install_test.go) — Install with httptest.Server
- [detect_test.go](../detect_test.go) — detect with injected lookPath
- [ensure_test.go](../ensure_test.go) — EnsureInstalled orchestration

## Stages

| Stage | Description | Dependency | Completed |
|-------|-------------|-------------|-----------|
| 1 | [Core: config, options, constants](stages/stage1_core.md) | — | [x] |
| 2 | [Install: download, extract, Install()](stages/stage2_install.md) | Stage 1 | [x] |
| 3 | [Detect: IsInstalled, GetPath, GetVersion, EnsureInstalled](stages/stage3_detect.md) | Stage 2 | [x] |

> Integration and cleanup live in [`client/docs/PLAN.md`](../../../client/docs/PLAN.md) — that package owns its own refactoring.
