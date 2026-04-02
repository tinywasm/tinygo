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
