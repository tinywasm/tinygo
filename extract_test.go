package tinygo

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractTarGz(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "extract-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	archiveFile, err := os.CreateTemp("", "archive-*.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(archiveFile.Name())
	defer archiveFile.Close()

	gw := gzip.NewWriter(archiveFile)
	tw := tar.NewWriter(gw)

	files := []struct {
		Name, Body string
	}{
		{"tinygo/bin/tinygo", "fake binary"},
		{"tinygo/README.md", "some info"},
	}

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0755,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	gw.Close()

	if err := extractTarGz(archiveFile.Name(), tmpDir); err != nil {
		t.Errorf("failed to extract tar.gz: %v", err)
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file.Name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", path)
		}
	}
}

func TestExtractZip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "extract-zip-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	archiveFile, err := os.CreateTemp("", "archive-*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(archiveFile.Name())
	defer archiveFile.Close()

	zw := zip.NewWriter(archiveFile)

	files := []struct {
		Name, Body string
	}{
		{"tinygo/bin/tinygo.exe", "fake binary"},
		{"tinygo/README.md", "some info"},
	}

	for _, file := range files {
		f, err := zw.Create(file.Name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(file.Body)); err != nil {
			t.Fatal(err)
		}
	}
	zw.Close()

	if err := extractZip(archiveFile.Name(), tmpDir); err != nil {
		t.Errorf("failed to extract zip: %v", err)
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file.Name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", path)
		}
	}
}

func TestZipSlip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipslip-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test tar.gz
	archiveFile, _ := os.CreateTemp("", "slip-*.tar.gz")
	gw := gzip.NewWriter(archiveFile)
	tw := tar.NewWriter(gw)
	tw.WriteHeader(&tar.Header{Name: "../../etc/passwd", Size: 0})
	tw.Close()
	gw.Close()
	archiveFile.Close()
	err = extractTarGz(archiveFile.Name(), tmpDir)
	if err == nil || !strings.Contains(err.Error(), "illegal file path") {
		t.Errorf("expected illegal file path error for tar.gz, got %v", err)
	}

	// Test zip
	archiveFile2, _ := os.CreateTemp("", "slip-*.zip")
	zw := zip.NewWriter(archiveFile2)
	zw.Create("../../etc/passwd")
	zw.Close()
	archiveFile2.Close()
	err = extractZip(archiveFile2.Name(), tmpDir)
	if err == nil || !strings.Contains(err.Error(), "illegal file path") {
		t.Errorf("expected illegal file path error for zip, got %v", err)
	}
}
