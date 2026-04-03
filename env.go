package tinygo

import (
	"os"
	"path/filepath"
	"strings"
)

// GetRoot returns the TINYGOROOT for a local install, or "" if tinygo is in PATH.
func GetRoot(opts ...Option) string {
	c := newConfig(opts...)
	if _, err := c.lookPath("tinygo"); err == nil {
		return ""
	}
	return filepath.Join(c.installDir, "tinygo")
}

// GetEnv returns os.Environ() + TINYGOROOT + prepended PATH for local installs.
// Safe to assign directly to exec.Command.Env.
func GetEnv(opts ...Option) []string {
	root := GetRoot(opts...)
	if root == "" {
		return os.Environ()
	}

	env := os.Environ()
	// We might add TINYGOROOT if not present, and we definitely update PATH.
	newEnv := make([]string, 0, len(env)+1)

	foundTinyRoot := false
	foundPath := false
	pathKey := "PATH" // Default to PATH, but will match actual case from env

	for _, e := range env {
		key, val, found := strings.Cut(e, "=")
		if !found {
			newEnv = append(newEnv, e)
			continue
		}

		upperKey := strings.ToUpper(key)
		if upperKey == "TINYGOROOT" {
			newEnv = append(newEnv, key+"="+root)
			foundTinyRoot = true
		} else if upperKey == "PATH" {
			pathKey = key
			foundPath = true
			binDir := filepath.Join(root, "bin")
			newEnv = append(newEnv, key+"="+binDir+string(os.PathListSeparator)+val)
		} else {
			newEnv = append(newEnv, e)
		}
	}

	if !foundTinyRoot {
		newEnv = append(newEnv, "TINYGOROOT="+root)
	}

	if !foundPath {
		binDir := filepath.Join(root, "bin")
		newEnv = append(newEnv, pathKey+"="+binDir)
	}

	return newEnv
}
