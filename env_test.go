package tinygo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetRoot(t *testing.T) {
	t.Run("SystemInstall", func(t *testing.T) {
		root := GetRoot(withLookPath(func(string) (string, error) {
			return "/usr/bin/tinygo", nil
		}))
		if root != "" {
			t.Errorf("expected empty root for system install, got %q", root)
		}
	})

	t.Run("LocalInstall", func(t *testing.T) {
		tmpDir := t.TempDir()
		root := GetRoot(
			WithInstallDir(tmpDir),
			withLookPath(func(string) (string, error) {
				return "", fmt.Errorf("not found")
			}),
		)
		expected := filepath.Join(tmpDir, "tinygo")
		if root != expected {
			t.Errorf("expected root %q, got %q", expected, root)
		}
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("SystemInstall", func(t *testing.T) {
		env := GetEnv(withLookPath(func(string) (string, error) {
			return "/usr/bin/tinygo", nil
		}))
		// Should be equal to os.Environ()
		expected := os.Environ()
		if len(env) != len(expected) {
			t.Errorf("expected env length %d, got %d", len(expected), len(env))
		}
	})

	t.Run("LocalInstall", func(t *testing.T) {
		tmpDir := t.TempDir()
		root := filepath.Join(tmpDir, "tinygo")
		env := GetEnv(
			WithInstallDir(tmpDir),
			withLookPath(func(string) (string, error) {
				return "", fmt.Errorf("not found")
			}),
		)

		var tinygoRoot, path string
		for _, e := range env {
			if strings.HasPrefix(e, "TINYGOROOT=") {
				tinygoRoot = e[len("TINYGOROOT="):]
			} else if strings.HasPrefix(e, "PATH=") {
				path = e[len("PATH="):]
			}
		}

		if tinygoRoot != root {
			t.Errorf("expected TINYGOROOT %q, got %q", root, tinygoRoot)
		}

		expectedPathPrefix := filepath.Join(root, "bin") + string(os.PathListSeparator)
		if !strings.HasPrefix(path, expectedPathPrefix) {
			t.Errorf("expected PATH to start with %q, got %q", expectedPathPrefix, path)
		}

		// Verify os.Environ() was not mutated
		found := false
		for _, e := range os.Environ() {
			if strings.HasPrefix(e, "TINYGOROOT=") && e[len("TINYGOROOT="):] == root {
				found = true
				break
			}
		}
		if found {
			t.Error("os.Environ() was mutated")
		}
	})

	t.Run("OverridesExistingTINYGOROOT", func(t *testing.T) {
		t.Setenv("TINYGOROOT", "/stale/path")
		tmpDir := t.TempDir()
		root := filepath.Join(tmpDir, "tinygo")
		env := GetEnv(
			WithInstallDir(tmpDir),
			withLookPath(func(string) (string, error) {
				return "", fmt.Errorf("not found")
			}),
		)

		count := 0
		var tinygoRoot string
		for _, e := range env {
			if strings.HasPrefix(e, "TINYGOROOT=") {
				tinygoRoot = e[len("TINYGOROOT="):]
				count++
			}
		}

		if count != 1 {
			t.Errorf("expected 1 TINYGOROOT entry, got %d", count)
		}
		if tinygoRoot != root {
			t.Errorf("expected TINYGOROOT %q, got %q", root, tinygoRoot)
		}
	})
}
