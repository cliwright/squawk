package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()

	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	write("a.yaml", `templates:
  deploy-failed:
    channel: "#alerts"
    text: "deploy failed"
`)
	write("b.yaml", `templates:
  test-failed:
    channel: "#ci"
    text: "tests failed"
`)

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfg.Templates) != 2 {
		t.Fatalf("expected 2 templates, got %d", len(cfg.Templates))
	}

	if cfg.Templates["deploy-failed"].Channel != "#alerts" {
		t.Errorf("wrong channel for deploy-failed")
	}
	if cfg.Templates["test-failed"].Channel != "#ci" {
		t.Errorf("wrong channel for test-failed")
	}
}

func TestLoadDuplicateKey(t *testing.T) {
	dir := t.TempDir()

	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	write("a.yaml", `templates:
  deploy-failed:
    channel: "#alerts"
    text: "first"
`)
	write("b.yaml", `templates:
  deploy-failed:
    channel: "#other"
    text: "duplicate"
`)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for duplicate template key")
	}
}
