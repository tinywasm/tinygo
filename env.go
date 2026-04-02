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
	newEnv := make([]string, 0, len(env)+1)
	rootSet := false

	for _, e := range env {
		if strings.HasPrefix(e, "TINYGOROOT=") {
			newEnv = append(newEnv, "TINYGOROOT="+root)
			rootSet = true
		} else if strings.HasPrefix(e, "PATH=") {
			path := e[len("PATH="):]
			newPath := filepath.Join(root, "bin") + string(os.PathListSeparator) + path
			newEnv = append(newEnv, "PATH="+newPath)
		} else {
			newEnv = append(newEnv, e)
		}
	}

	if !rootSet {
		newEnv = append(newEnv, "TINYGOROOT="+root)
	}

	return newEnv
}
