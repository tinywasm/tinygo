package tinygo

import (
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
	c := newConfig(opts...)
	root := GetRoot(opts...)
	if root == "" {
		return c.environ
	}

	env := make([]string, 0, len(c.environ)+1)
	pathKey := "PATH"
	if c.goos == "windows" {
		pathKey = "Path"
	}

	tinygoRootEntry := "TINYGOROOT=" + root
	pathFound := false
	rootFound := false

	for _, e := range c.environ {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key := pair[0]
		if strings.EqualFold(key, "TINYGOROOT") {
			env = append(env, tinygoRootEntry)
			rootFound = true
			continue
		}
		if strings.EqualFold(key, pathKey) {
			newPath := filepath.Join(root, "bin") + string(filepath.ListSeparator) + pair[1]
			env = append(env, key+"="+newPath)
			pathFound = true
			continue
		}
		env = append(env, e)
	}

	if !rootFound {
		env = append(env, tinygoRootEntry)
	}
	if !pathFound {
		env = append(env, pathKey+"="+filepath.Join(root, "bin"))
	}

	return env
}
