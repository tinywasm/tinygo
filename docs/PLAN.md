# PLAN: tinygo package — Environment Variables for Local Installations

## Problem

`TINYGOROOT` must point to the TinyGo root directory so the compiler can locate
`src/`, `targets/`, and `lib/`. Without it, `tinygo version` works but `tinygo build`
fails with `could not load target: targets/wasm.json: no such file`.

System installs (apt, brew, winget) set `TINYGOROOT` automatically.
Local installs (`~/.tinywasm/tinygo/`) do not — this package must handle it.

## Design Decisions

- **No shell modification**: do not write to `~/.bashrc`, `~/.zshrc`, registry, or call `setx`.
- **No `os.Setenv`**: mutating the global process env is wrong for a library.
- **Inject at exec time**: expose `GetEnv()` so consumers pass it to `exec.Command.Env`.
- **`GetRoot()` as single point of truth**: root is always `{installDir}/tinygo/` for local installs.
- **No-op for system installs**: if tinygo is found in PATH, `GetEnv()` returns `os.Environ()` unchanged.

## Public API (additions to tinygo package)

```go
// GetRoot returns the TINYGOROOT for a local install, or "" if tinygo is in PATH.
func GetRoot(opts ...Option) string

// GetEnv returns os.Environ() + TINYGOROOT + prepended PATH for local installs.
// Safe to assign directly to exec.Command.Env.
func GetEnv(opts ...Option) []string
```

## Consumer Pattern (gobuild / client)

```go
binPath, _ := tinygo.EnsureInstalled()
cmd := exec.Command(binPath, "build", "-o", "out.wasm", ".")
cmd.Env = tinygo.GetEnv()   // ← replaces os.Setenv("PATH", ...) in wasmbuild.go
out, err := cmd.CombinedOutput()
```

## Current Bug in client/wasmbuild.go

```go
// Line 26-27 — sets PATH only, ignores TINYGOROOT, mutates global process env:
newPath := filepath.Dir(tinyGoPath) + string(os.PathListSeparator) + os.Getenv("PATH")
os.Setenv("PATH", newPath)  // ← BUG: wrong approach, fixed by Stage 5
```

## Files

| File | Role |
|------|------|
| [env.go](../env.go) | `GetRoot()` + `GetEnv()` |
| [env_test.go](../env_test.go) | unit tests |
| [client/wasmbuild.go](../../../client/wasmbuild.go) | consumer fix (Stage 5) |

## Stages

| Stage | Description | Dependency | Completed |
|-------|-------------|-------------|-----------|
| 4 | [env.go: GetRoot, GetEnv](stages/stage4_env.md) | Stage 3 | [x] |
| 5 | [Fix wasmbuild.go: use GetEnv() instead of os.Setenv](stages/stage5_wasmbuild_fix.md) | Stage 4 | [ ] |
