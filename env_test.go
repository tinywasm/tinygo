package tinygo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetRoot_SystemInstall(t *testing.T) {
	opts := []Option{
		withLookPath(func(name string) (string, error) {
			return "/usr/bin/tinygo", nil
		}),
	}
	root := GetRoot(opts...)
	if root != "" {
		t.Errorf("expected empty root for system install, got %q", root)
	}
}

func TestGetRoot_LocalInstall(t *testing.T) {
	tmpDir := t.TempDir()
	opts := []Option{
		withLookPath(func(name string) (string, error) {
			return "", fmt.Errorf("not found")
		}),
		WithInstallDir(tmpDir),
	}
	root := GetRoot(opts...)
	expected := filepath.Join(tmpDir, "tinygo")
	if root != expected {
		t.Errorf("expected %q, got %q", expected, root)
	}
}

func TestGetEnv_SystemInstall(t *testing.T) {
	opts := []Option{
		withLookPath(func(name string) (string, error) {
			return "/usr/bin/tinygo", nil
		}),
	}
	env := GetEnv(opts...)
	expected := os.Environ()
	if len(env) != len(expected) {
		t.Errorf("expected %d env vars, got %d", len(expected), len(env))
	}
	for i := range env {
		if env[i] != expected[i] {
			t.Errorf("mismatch at index %d: expected %q, got %q", i, expected[i], env[i])
		}
	}
}

func TestGetEnv_LocalInstall(t *testing.T) {
	tmpDir := t.TempDir()
	opts := []Option{
		withLookPath(func(name string) (string, error) {
			return "", fmt.Errorf("not found")
		}),
		WithInstallDir(tmpDir),
	}
	env := GetEnv(opts...)
	root := filepath.Join(tmpDir, "tinygo")

	foundTinyRoot := false
	foundPath := false
	for _, e := range env {
		key, _, found := strings.Cut(e, "=")
		if !found {
			continue
		}
		upperKey := strings.ToUpper(key)
		if upperKey == "TINYGOROOT" {
			foundTinyRoot = true
			if e != "TINYGOROOT="+root {
				t.Errorf("expected TINYGOROOT=%s, got %s", root, e)
			}
		}
		if upperKey == "PATH" {
			foundPath = true
			binDir := filepath.Join(root, "bin")
			if !strings.Contains(e, binDir) {
				t.Errorf("PATH does not contain bin dir %q: %q", binDir, e)
			}
			if !strings.HasPrefix(e, key+"="+binDir+string(os.PathListSeparator)) && e != key+"="+binDir {
				t.Errorf("PATH should start with %q: %q", binDir, e)
			}
		}
	}

	if !foundTinyRoot {
		t.Error("TINYGOROOT not found in env")
	}
	if !foundPath {
		t.Error("PATH not found in env")
	}
}

func TestGetEnv_OverridesExistingTINYGOROOT(t *testing.T) {
	// Set a dummy TINYGOROOT in the current process to ensure it's overridden
	os.Setenv("TINYGOROOT", "/stale/path")
	defer os.Unsetenv("TINYGOROOT")

	tmpDir := t.TempDir()
	opts := []Option{
		withLookPath(func(name string) (string, error) {
			return "", fmt.Errorf("not found")
		}),
		WithInstallDir(tmpDir),
	}
	env := GetEnv(opts...)
	root := filepath.Join(tmpDir, "tinygo")

	foundCount := 0
	for _, e := range env {
		key, _, found := strings.Cut(e, "=")
		if found && strings.ToUpper(key) == "TINYGOROOT" {
			foundCount++
			if e != key+"="+root {
				t.Errorf("expected TINYGOROOT=%s, got %s", root, e)
			}
		}
	}
	if foundCount != 1 {
		t.Errorf("expected exactly 1 TINYGOROOT entry, got %d", foundCount)
	}
}
