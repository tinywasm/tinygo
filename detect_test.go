package tinygo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetPath_PATH(t *testing.T) {
	mockPath := "/usr/bin/tinygo"
	lookPath := func(string) (string, error) {
		return mockPath, nil
	}

	c := newConfig(withLookPath(lookPath))
	p, err := getPath(c)
	if err != nil {
		t.Fatal(err)
	}

	if p != mockPath {
		t.Errorf("expected %s, got %s", mockPath, p)
	}
}

func TestGetPath_Local(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	tmpDir, _ := os.MkdirTemp("", "tinygo-local-*")
	defer os.RemoveAll(tmpDir)

	binPath := filepath.Join(tmpDir, "tinygo/bin/tinygo")
	os.MkdirAll(filepath.Dir(binPath), 0755)
	os.WriteFile(binPath, []byte("fake"), 0755)

	c := newConfig(withLookPath(lookPath), WithInstallDir(tmpDir), withGOOS("linux"))
	p, err := getPath(c)
	if err != nil {
		t.Fatal(err)
	}

	if p != binPath {
		t.Errorf("expected %s, got %s", binPath, p)
	}
}

func TestIsInstalled(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "/usr/bin/tinygo", nil
	}

	if !IsInstalled(withLookPath(lookPath)) {
		t.Errorf("expected IsInstalled to be true when mock lookPath succeeds")
	}

	lookPathFails := func(string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	tmpDir, _ := os.MkdirTemp("", "tinygo-not-installed-*")
	defer os.RemoveAll(tmpDir)

	if IsInstalled(withLookPath(lookPathFails), WithInstallDir(tmpDir)) {
		t.Errorf("expected IsInstalled to be false when tinygo is missing")
	}
}

func TestGetVersion(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "tinygo-version-*")
	defer os.RemoveAll(tmpDir)

	binPath := filepath.Join(tmpDir, "tinygo")
	if _, err := os.Stat("/bin/bash"); err == nil {
		// Create a shell script as a mock binary
		script := "#!/bin/bash\necho \"tinygo version 0.40.1 linux/amd64 (using go version go1.22.0 and LLVM version 18.1.2)\"\n"
		os.WriteFile(binPath, []byte(script), 0755)

		lookPath := func(string) (string, error) {
			return binPath, nil
		}

		version, err := GetVersion(withLookPath(lookPath))
		if err != nil {
			t.Fatalf("GetVersion failed: %v", err)
		}

		expected := "tinygo version 0.40.1 linux/amd64 (using go version go1.22.0 and LLVM version 18.1.2)"
		if version != expected {
			t.Errorf("expected %s, got %s", expected, version)
		}
	} else {
		t.Skip("Skipping GetVersion test because /bin/bash is not available for mocking")
	}
}
