package tinygo

import (
	"fmt"
	"os"
	"path/filepath"
)

func Install(opts ...Option) error {
	c := newConfig(opts...)

	bin := c.binPath()
	if _, err := os.Stat(bin); err == nil {
		c.logger("TinyGo already installed at " + bin)
		return nil
	}

	tmpFile, err := c.download()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	c.logger("Extracting tinygo to " + c.installDir)
	if err := os.MkdirAll(c.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install dir: %w", err)
	}

	if c.goos == "windows" {
		err = extractZip(tmpFile, c.installDir)
	} else {
		err = extractTarGz(tmpFile, c.installDir)
	}

	if err != nil {
		// cleanup partial install - only the 'tinygo' directory within installDir
		os.RemoveAll(filepath.Join(c.installDir, "tinygo"))
		return fmt.Errorf("failed to extract tinygo: %w", err)
	}

	if _, err := os.Stat(bin); err != nil {
		return fmt.Errorf("tinygo binary not found after extraction: %w", err)
	}

	return nil
}
