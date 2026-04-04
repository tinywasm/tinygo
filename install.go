package tinygo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Install(opts ...Option) error {
	c := newConfig(opts...)

	if c.goos == "windows" {
		return installWindows(c, opts...)
	}
	return installUnix(c, opts...)
}

// installUnix installs TinyGo on Linux/macOS by extracting the official tarball
// into the install directory (default: /usr/local), following
// https://tinygo.org/getting-started/install/
func installUnix(c *config, opts ...Option) error {
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

	if err := extractTarGz(tmpFile, c.installDir); err != nil {
		os.RemoveAll(filepath.Join(c.installDir, "tinygo"))
		return fmt.Errorf("failed to extract tinygo: %w", err)
	}

	if _, err := os.Stat(bin); err != nil {
		return fmt.Errorf("tinygo binary not found after extraction: %w", err)
	}

	ver, err := GetVersion(opts...)
	if err != nil {
		os.RemoveAll(filepath.Join(c.installDir, "tinygo"))
		return fmt.Errorf("tinygo verification failed: %w", err)
	}
	c.logger("TinyGo verified: " + ver)

	// Create symlink in installDir/bin/ so tinygo is reachable from PATH
	// without shell modifications:
	//   /usr/local     → /usr/local/bin/tinygo      (system install, sudo)
	//   ~/.local       → ~/.local/bin/tinygo         (user install, no sudo)
	symlinkDir := filepath.Join(c.installDir, "bin")
	symlink := filepath.Join(symlinkDir, "tinygo")
	if bin != symlink {
		if err := os.MkdirAll(symlinkDir, 0755); err != nil {
			c.logger("Warning: could not create bin dir " + symlinkDir + ": " + err.Error())
		} else {
			os.Remove(symlink)
			if err := os.Symlink(bin, symlink); err != nil {
				c.logger("Warning: could not create symlink " + symlink + ": " + err.Error())
			} else {
				c.logger("Symlink created: " + symlink + " → " + bin)
			}
		}
	}

	return nil
}

// installWindows installs TinyGo on Windows via Scoop, following
// https://tinygo.org/getting-started/install/
//
// If Scoop is not installed, it installs it first using the official
// PowerShell one-liner, then runs `scoop install tinygo`.
func installWindows(c *config, opts ...Option) error {
	if err := ensureScoop(c); err != nil {
		return err
	}

	c.logger("Installing TinyGo via scoop...")
	if out, err := exec.Command("scoop", "install", "tinygo").CombinedOutput(); err != nil {
		return fmt.Errorf("scoop install tinygo: %w\n%s", err, out)
	}

	ver, err := GetVersion(opts...)
	if err != nil {
		return fmt.Errorf("tinygo verification failed: %w", err)
	}
	c.logger("TinyGo verified: " + ver)

	return nil
}

// removeExisting removes a TinyGo installation before replacing it.
//
//   - Linux: if installed via apt/dpkg, runs `sudo apt-get remove -y tinygo`.
//     Otherwise removes the tinygo directory and any symlink in /usr/local/bin.
//   - macOS: removes the tinygo directory under installDir.
//   - Windows: uses `scoop uninstall tinygo` if scoop is available.
func removeExisting(c *config, binPath string) error {
	switch c.goos {
	case "windows":
		if _, err := c.lookPath("scoop"); err == nil {
			c.logger("Removing old TinyGo via scoop...")
			if out, err := exec.Command("scoop", "uninstall", "tinygo").CombinedOutput(); err != nil {
				return fmt.Errorf("scoop uninstall tinygo: %w\n%s", err, out)
			}
		}
		return nil

	default: // linux / darwin
		// Check if installed via apt/dpkg.
		if c.goos == "linux" {
			if out, err := exec.Command("dpkg", "-s", "tinygo").CombinedOutput(); err == nil && len(out) > 0 {
				c.logger("Removing old TinyGo apt package...")
				if out, err := exec.Command("sudo", "apt-get", "remove", "-y", "tinygo").CombinedOutput(); err != nil {
					return fmt.Errorf("apt-get remove tinygo: %w\n%s", err, out)
				}
				return nil
			}
		}
		// Tarball install: remove the tinygo subdirectory and its symlink.
		tinygoDir := filepath.Join(c.installDir, "tinygo")
		c.logger("Removing " + tinygoDir + "...")
		if err := os.RemoveAll(tinygoDir); err != nil {
			return fmt.Errorf("failed to remove %s: %w", tinygoDir, err)
		}
		symlink := filepath.Join(c.installDir, "bin", "tinygo")
		if _, err := os.Lstat(symlink); err == nil {
			os.Remove(symlink)
		}
		return nil
	}
}

// ensureScoop installs Scoop if it is not already available in PATH.
func ensureScoop(c *config) error {
	if _, err := c.lookPath("scoop"); err == nil {
		return nil
	}

	c.logger("Scoop not found. Installing Scoop...")

	// Official Scoop install: https://scoop.sh
	script := `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser; ` +
		`Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression`

	if out, err := exec.Command("powershell", "-NoProfile", "-Command", script).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install Scoop: %w\n%s", err, out)
	}

	// Scoop installs to %USERPROFILE%\scoop\shims which is not in the current
	// process PATH. Prepend it so lookPath can find scoop without a new shell.
	if home, err := os.UserHomeDir(); err == nil {
		scoopShims := filepath.Join(home, "scoop", "shims")
		os.Setenv("PATH", scoopShims+string(os.PathListSeparator)+os.Getenv("PATH"))
	}

	if _, err := c.lookPath("scoop"); err != nil {
		return fmt.Errorf("scoop not found after installation; restart your shell and retry")
	}

	c.logger("Scoop installed.")
	return nil
}
