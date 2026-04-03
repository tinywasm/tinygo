package tinygo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsInstalled(opts ...Option) bool {
	_, err := GetPath(opts...)
	return err == nil
}

func GetPath(opts ...Option) (string, error) {
	return getPath(newConfig(opts...))
}

func getPath(c *config) (string, error) {
	// 1. PATH
	if p, err := c.lookPath("tinygo"); err == nil {
		return p, nil
	}

	// 2. local
	bin := c.binPath()
	if _, err := os.Stat(bin); err == nil {
		return bin, nil
	}

	return "", fmt.Errorf("tinygo not found in PATH or in local installation")
}

func GetVersion(opts ...Option) (string, error) {
	p, err := GetPath(opts...)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(p, "version")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run tinygo version: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// installedVersion returns just the semver number from `tinygo version` output.
// e.g. "tinygo version 0.39.0 linux/amd64 ..." → "0.39.0"
func installedVersion(opts ...Option) (string, error) {
	full, err := GetVersion(opts...)
	if err != nil {
		return "", err
	}
	// output format: "tinygo version X.Y.Z ..."
	fields := strings.Fields(full)
	if len(fields) < 3 {
		return "", fmt.Errorf("unexpected tinygo version output: %q", full)
	}
	return fields[2], nil
}
