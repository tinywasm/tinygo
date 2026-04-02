package tinygo

import (
	"fmt"
	"io"
	"os"
)

func (c *config) download() (string, error) {
	url := c.downloadURL()
	c.logger(fmt.Sprintf("Downloading tinygo from %s", url))

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download tinygo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to download tinygo: status code %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "tinygo-*.tmp")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save tinygo to temp file: %w", err)
	}

	return tmpFile.Name(), nil
}
