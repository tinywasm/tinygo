package tinygo

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	DefaultVersion    = "0.40.1"
	DefaultInstallDir = ".tinywasm"
)

type config struct {
	version    string
	installDir string
	logger     func(string)
	lookPath   func(string) (string, error)
	httpClient *http.Client
}

type Option func(*config)

func WithVersion(v string) Option {
	return func(c *config) {
		c.version = v
	}
}

func WithInstallDir(dir string) Option {
	return func(c *config) {
		c.installDir = dir
	}
}

func WithLogger(f func(string)) Option {
	return func(c *config) {
		c.logger = f
	}
}

func newConfig(opts ...Option) *config {
	homeDir, _ := os.UserHomeDir()
	c := &config{
		version:    DefaultVersion,
		installDir: filepath.Join(homeDir, DefaultInstallDir),
		logger:     func(string) {},
		lookPath:   exec.LookPath,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// internal options for testing (not exported)
func withLookPath(f func(string) (string, error)) Option {
	return func(c *config) {
		c.lookPath = f
	}
}

func withHTTPClient(cl *http.Client) Option {
	return func(c *config) {
		c.httpClient = cl
	}
}

func (c *config) binPath() string {
	bin := filepath.Join(c.installDir, "tinygo", "bin", "tinygo")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	return bin
}

func (c *config) downloadURL() string {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	// eg: tinygo0.40.1.linux-amd64.tar.gz
	return "https://github.com/tinygo-org/tinygo/releases/download/v" +
		c.version + "/tinygo" + c.version + "." + runtime.GOOS + "-" + runtime.GOARCH + "." + ext
}
