# Stage 5 — Fix wasmbuild.go: use GetEnv() instead of os.Setenv

### Goal

Replace the broken `os.Setenv("PATH", ...)` pattern in `client/wasmbuild.go` with
`tinygo.GetEnv()` so that TINYGOROOT is correctly injected into the compiler subprocess.

### Dependency
- Stage 4 completed (`env.go`: `GetRoot`, `GetEnv`)

### Context

`gobuild/compiler.go` already supports env injection via `h.config.Env`:

```go
// compiler.go:25-27
if len(h.config.Env) > 0 {
    comp.cmd.Env = append(os.Environ(), h.config.Env...)
}
```

The fix is to stop mutating the process env in `wasmbuild.go` and instead populate
`h.config.Env` with the output of `tinygo.GetEnv()` before calling `w.Compile()`.

### Steps

- [ ] In `client/wasmbuild.go`:
  - Remove lines 24-28 (`os.Setenv("PATH", newPath)`).
  - After `EnsureInstalled()` succeeds, call `tinygo.GetEnv()` and store the result.
  - Pass the env slice to the `WasmClient` config before `Compile()`.

- [ ] Verify `gobuild` config has a field for env injection (already exists as `config.Env`).

- [ ] Add/update test in `client/tests/wasmbuild_test.go` to verify that:
  - `TINYGOROOT` is present in the env passed to the subprocess when using local install.
  - `os.Environ()` is not mutated.

### Files
- [client/wasmbuild.go](../../../../client/wasmbuild.go)
- [client/tests/wasmbuild_test.go](../../../../client/tests/wasmbuild_test.go)
