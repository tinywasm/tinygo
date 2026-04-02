package client

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tinywasm/tinygo"
)

type Config struct {
	Env []string
}

type WasmClient struct {
	config Config
}

func NewWasmClient(cfg Config) *WasmClient {
	return &WasmClient{config: cfg}
}

func (w *WasmClient) Compile(src string) error {
	tinyGoPath, err := tinygo.EnsureInstalled()
	if err != nil {
		return fmt.Errorf("failed to ensure tinygo: %w", err)
	}

	// Replaces old broken pattern:
	// newPath := filepath.Dir(tinyGoPath) + string(os.PathListSeparator) + os.Getenv("PATH")
	// os.Setenv("PATH", newPath)

	w.config.Env = tinygo.GetEnv()

	cmd := exec.Command(tinyGoPath, "build", "-o", "out.wasm", src)
	if len(w.config.Env) > 0 {
		cmd.Env = w.config.Env
	} else {
		cmd.Env = os.Environ()
	}

	// In a real implementation we would run the command here
	// _, err = cmd.CombinedOutput()
	// return err

	return nil
}
