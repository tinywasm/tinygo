package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/tinywasm/tinygo"
)

func main() {
	var (
		version = flag.String("version", tinygo.DefaultVersion, "TinyGo version to install (Linux/macOS only; ignored on Windows via Scoop)")
		dir     = flag.String("dir", "", "Installation directory (Linux/macOS only; default: /usr/local)")
		verbose = flag.Bool("v", false, "Verbose output")
		help    = flag.Bool("h", false, "Show help")
	)

	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	opts := []tinygo.Option{}

	if *verbose {
		opts = append(opts, tinygo.WithLogger(func(msg string) {
			fmt.Println("[tinygo]", msg)
		}))
	}

	// -version and -dir only apply on Linux/macOS (tarball install).
	// On Windows, Scoop manages version and directory automatically.
	if runtime.GOOS != "windows" {
		opts = append(opts, tinygo.WithVersion(*version))
		if *dir != "" {
			opts = append(opts, tinygo.WithInstallDir(*dir))
		}
	}

	binPath, err := tinygo.EnsureInstalled(opts...)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Printf("✓ TinyGo is ready\n")
	fmt.Printf("  Binary: %s\n", binPath)
	fmt.Printf("\nRun 'hash -r' or open a new terminal to use the updated tinygo.\n")
}

func printUsage() {
	w := os.Stderr
	fmt.Fprintf(w, `Usage: tinygoinstall [options]

Installs TinyGo following the official instructions at
https://tinygo.org/getting-started/install/

  Linux/macOS  Extracts the official tarball into /usr/local (requires sudo).
               Use -dir to install into a user-writable directory instead.
  Windows      Installs via Scoop (installs Scoop first if not present; no sudo needed).

Options:
  -version string   TinyGo version to install, e.g. 0.35.0 (Linux/macOS only, default: %s)
  -dir string       Installation directory (Linux/macOS only, default: /usr/local)
  -v                Verbose output
  -h                Show this help message

Examples:
  sudo tinygoinstall                    # Linux/macOS — install to /usr/local
  tinygoinstall -dir ~/.local           # Linux/macOS — no sudo, user directory
  sudo tinygoinstall -version 0.35.0    # specific version
  tinygoinstall                         # Windows — uses Scoop, no sudo needed
`, tinygo.DefaultVersion)
}
