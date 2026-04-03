package tinygo

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGetRoot_SystemInstall(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "/usr/bin/tinygo", nil
	}

	root := GetRoot(withLookPath(lookPath))
	if root != "" {
		t.Errorf("expected empty string for system install, got %s", root)
	}
}

func TestGetRoot_LocalInstall(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "", fmt.Errorf("not found")
	}
	tmpDir := "/tmp/tinywasm"
	expected := filepath.Join(tmpDir, "tinygo")

	root := GetRoot(withLookPath(lookPath), WithInstallDir(tmpDir))
	if root != expected {
		t.Errorf("expected %s, got %s", expected, root)
	}
}

func TestGetEnv_SystemInstall(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "/usr/bin/tinygo", nil
	}
	mockEnv := []string{"FOO=BAR", "PATH=/usr/bin"}

	env := GetEnv(withLookPath(lookPath), withEnviron(mockEnv))
	if !reflect.DeepEqual(env, mockEnv) {
		t.Errorf("expected environment to be unchanged, got %v", env)
	}
}

func TestGetEnv_LocalInstall(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "", fmt.Errorf("not found")
	}
	tmpDir := "/tmp/tinywasm"
	mockEnv := []string{"FOO=BAR", "PATH=/usr/bin"}
	root := filepath.Join(tmpDir, "tinygo")

	env := GetEnv(withLookPath(lookPath), WithInstallDir(tmpDir), withEnviron(mockEnv), withGOOS("linux"))

	foundTinygoRoot := false
	foundPath := false
	for _, e := range env {
		if e == "TINYGOROOT="+root {
			foundTinygoRoot = true
		}
		if strings.HasPrefix(e, "PATH=") {
			foundPath = true
			expectedPath := filepath.Join(root, "bin") + string(filepath.ListSeparator) + "/usr/bin"
			if e != "PATH="+expectedPath {
				t.Errorf("expected PATH=%s, got %s", expectedPath, e)
			}
		}
	}

	if !foundTinygoRoot {
		t.Error("TINYGOROOT not found in environment")
	}
	if !foundPath {
		t.Error("PATH not found in environment")
	}

	// Verify original env not mutated
	if mockEnv[1] != "PATH=/usr/bin" {
		t.Error("original environment was mutated")
	}
}

func TestGetEnv_OverridesExistingTINYGOROOT(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "", fmt.Errorf("not found")
	}
	tmpDir := "/tmp/tinywasm"
	mockEnv := []string{"TINYGOROOT=/stale/path", "PATH=/usr/bin"}
	root := filepath.Join(tmpDir, "tinygo")

	env := GetEnv(withLookPath(lookPath), WithInstallDir(tmpDir), withEnviron(mockEnv), withGOOS("linux"))

	foundNew := false
	foundOld := false
	for _, e := range env {
		if e == "TINYGOROOT="+root {
			foundNew = true
		}
		if e == "TINYGOROOT=/stale/path" {
			foundOld = true
		}
	}

	if !foundNew {
		t.Error("new TINYGOROOT not found")
	}
	if foundOld {
		t.Error("old TINYGOROOT still present")
	}
}

func TestGetEnv_Windows(t *testing.T) {
	lookPath := func(string) (string, error) {
		return "", fmt.Errorf("not found")
	}
	tmpDir := `C:\Users\test\.tinywasm`
	mockEnv := []string{"FOO=BAR", "Path=C:\\Windows\\system32"}
	root := filepath.Join(tmpDir, "tinygo")

	env := GetEnv(withLookPath(lookPath), WithInstallDir(tmpDir), withEnviron(mockEnv), withGOOS("windows"))

	foundPath := false
	for _, e := range env {
		if strings.HasPrefix(e, "Path=") {
			foundPath = true
			expectedPath := filepath.Join(root, "bin") + string(os.PathListSeparator) + `C:\Windows\system32`
			if e != "Path="+expectedPath {
				t.Errorf("expected Path=%s, got %s", expectedPath, e)
			}
		}
	}

	if !foundPath {
		t.Error("Path not found in environment")
	}
}
