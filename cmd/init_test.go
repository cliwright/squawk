package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cliwright/squawk/config"
)

func TestRunInitCreatesDirectoryAndFile(t *testing.T) {
	dir := t.TempDir()
	squawkDir := filepath.Join(dir, ".squawk")

	if err := runInit(squawkDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(squawkDir)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected .squawk to be a directory")
	}

	configPath := filepath.Join(squawkDir, "squawk.yaml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}

	cfg, err := config.Load(squawkDir)
	if err != nil {
		t.Fatalf("expected config to be loadable: %v", err)
	}

	if _, ok := cfg.Templates["failure"]; !ok {
		t.Error("expected 'failure' template to exist")
	}
	if _, ok := cfg.Templates["success"]; !ok {
		t.Error("expected 'success' template to exist")
	}
}

func TestRunInitFailsIfDirectoryExists(t *testing.T) {
	dir := t.TempDir()
	squawkDir := filepath.Join(dir, ".squawk")

	if err := os.MkdirAll(squawkDir, 0o700); err != nil {
		t.Fatalf("setup: %v", err)
	}

	err := runInit(squawkDir)
	if err == nil {
		t.Fatal("expected error when directory already exists")
	}
}
