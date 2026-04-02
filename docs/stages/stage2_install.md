# Stage 2 — Install: download, extract, Install()

### Goal
Implement the download and installation logic for the TinyGo binary for all three operating systems.

### Dependency
- Stage 1 completed (config, options, constants)

### Steps

- [x] Create `download.go` with:
  - Function `(c *config) download() (string, error)` — downloads the file to a temporary directory.
  - Uses `net/http.Get` for downloading.
  - Validates HTTP status 200.
  - Returns the path to the downloaded temporary file.
  - Calls `c.logger` with progress messages if a logger is configured.

- [x] Create `extract.go` with:
  - Function `extractTarGz(src, dst string) error` — extracts .tar.gz files (Linux/macOS).
  - Function `extractZip(src, dst string) error` — extracts .zip files (Windows).
  - Both functions must:
    - Create the destination directory if it doesn't exist.
    - Protect against Zip Slip (path traversal).
    - Preserve file permissions.

- [x] Create `install.go` with:
  - Public function `Install(opts ...Option) error`.
  - Orchestrates: build config -> download -> extract -> verify binary exists.
  - Download must use `c.httpClient` (not `http.Get` directly) to enable test injection.
  - If the installation already exists in the destination path, return without doing anything (idempotent).
  - In case of an error during extraction, clean up partial files.

- [x] Create `extract_test.go` with:
  - tar.gz extraction test with programmatically created test archive.
  - zip extraction test with programmatically created test archive.
  - Zip Slip protection test (path traversal must return error).
  - Test verifying permissions of extracted files are preserved.

- [x] Create `install_test.go` with:
  - End-to-end test using `httptest.Server` to serve a fake tar.gz/zip archive.
  - Uses `WithInstallDir(tempDir)` to avoid polluting real filesystem.
  - Verifies binary exists at expected `binPath()` after install.
  - Tests idempotency: second `Install()` call returns immediately.
  - Tests cleanup on extraction error.

### Implementation Notes
- DO NOT use `sudo` or package managers.
- Download to temp, extract to `~/.tinywasm/tinygo/`, verify binary, clean up temp.
- On Windows use `archive/zip`, on Linux/macOS use `archive/tar` + `compress/gzip`.
- Reuse existing `untar()` logic from `client/tinygo_installer.go` as a reference.

### Files
- [download.go](../../download.go)
- [extract.go](../../extract.go)
- [install.go](../../install.go)
- [extract_test.go](../../extract_test.go)
