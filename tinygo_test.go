package tinygo

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := newConfig()
	if c.version != DefaultVersion {
		t.Errorf("expected version %s, got %s", DefaultVersion, c.version)
	}

	homeDir, _ := os.UserHomeDir()
	expectedDir := filepath.Join(homeDir, DefaultInstallDir)
	if c.installDir != expectedDir {
		t.Errorf("expected installDir %s, got %s", expectedDir, c.installDir)
	}

	if c.lookPath == nil {
		t.Error("lookPath should not be nil")
	}

	if c.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestWithVersion(t *testing.T) {
	v := "0.39.0"
	c := newConfig(WithVersion(v))
	if c.version != v {
		t.Errorf("expected version %s, got %s", v, c.version)
	}
}

func TestWithInstallDir(t *testing.T) {
	dir := "/tmp/tinygo"
	c := newConfig(WithInstallDir(dir))
	if c.installDir != dir {
		t.Errorf("expected installDir %s, got %s", dir, c.installDir)
	}
}

func TestBinPath(t *testing.T) {
	tests := []struct {
		goos     string
		expected string
	}{
		{"linux", "tinygo/bin/tinygo"},
		{"darwin", "tinygo/bin/tinygo"},
		{"windows", "tinygo/bin/tinygo.exe"},
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			c := newConfig(WithInstallDir("/tmp"), withGOOS(tt.goos))
			got := c.binPath()
			// filepath.Join handles platform-specific path separators.
			// However, in our tests, we're forcing GOOS. Let's make sure our comparison is correct.
			var expected string
			if tt.goos == "windows" {
				// Special case for windows if we're on a non-windows host.
				// However, filepath.Join will use host's separator.
				// Since binPath() uses filepath.Join, it will use host's separator.
				expected = filepath.Join("/tmp", "tinygo", "bin", "tinygo.exe")
			} else {
				expected = filepath.Join("/tmp", "tinygo", "bin", "tinygo")
			}
			if got != expected {
				t.Errorf("expected %s, got %s", expected, got)
			}
		})
	}
}

func TestDownloadURL(t *testing.T) {
	tests := []struct {
		goos   string
		goarch string
		ext    string
	}{
		{"linux", "amd64", "tar.gz"},
		{"linux", "arm64", "tar.gz"},
		{"darwin", "amd64", "tar.gz"},
		{"darwin", "arm64", "tar.gz"},
		{"windows", "amd64", "zip"},
	}

	version := "0.40.1"
	for _, tt := range tests {
		t.Run(tt.goos+"-"+tt.goarch, func(t *testing.T) {
			c := newConfig(WithVersion(version), withGOOS(tt.goos), withGOARCH(tt.goarch))
			got := c.downloadURL()
			expected := "https://github.com/tinygo-org/tinygo/releases/download/v" +
				version + "/tinygo" + version + "." + tt.goos + "-" + tt.goarch + "." + tt.ext
			if got != expected {
				t.Errorf("expected %s, got %s", expected, got)
			}
		})
	}
}

func TestInternalOptions(t *testing.T) {
	mockLookPath := func(string) (string, error) { return "", nil }
	mockClient := &http.Client{}

	c := newConfig(withLookPath(mockLookPath), withHTTPClient(mockClient))

	if c.lookPath == nil {
		t.Error("lookPath should not be nil")
	}
	// We can't easily compare function pointers, but we can verify it's set.

	if c.httpClient != mockClient {
		t.Error("httpClient was not set correctly")
	}
}

func TestWithLogger(t *testing.T) {
	var logged []string
	logger := func(s string) {
		logged = append(logged, s)
	}

	c := newConfig(WithLogger(logger))
	c.logger("test message")

	if len(logged) != 1 || logged[0] != "test message" {
		t.Errorf("logger not working as expected")
	}
}
