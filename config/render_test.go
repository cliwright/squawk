package config

import (
	"strings"
	"testing"
)

func TestRenderWithMentions(t *testing.T) {
	tmpl := Template{
		Channel:  "#alerts",
		Mentions: []string{"U111", "U222"},
		Text:     "deploy failed on `{{ .branch }}`\n{{ .mentions }}",
	}

	result, err := tmpl.Render(map[string]string{"branch": "main"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "deploy failed on `main`") {
		t.Errorf("expected rendered text, got: %s", result)
	}

	if !strings.Contains(result, "<@U111> <@U222>") {
		t.Errorf("expected mentions in rendered output, got: %s", result)
	}
}

func TestRenderWithGroupMentions(t *testing.T) {
	tmpl := Template{
		Channel:  "#alerts",
		Mentions: []string{"U111", "S222"},
		Text:     "deploy failed\n{{ .mentions }}",
	}

	result, err := tmpl.Render(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "<@U111>") {
		t.Errorf("expected user mention, got: %s", result)
	}
	if !strings.Contains(result, "<!subteam^S222>") {
		t.Errorf("expected subteam mention, got: %s", result)
	}
}

func TestRenderNoMentions(t *testing.T) {
	tmpl := Template{
		Channel: "#alerts",
		Text:    "hello {{ .repo }}",
	}

	result, err := tmpl.Render(map[string]string{"repo": "squawk"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "hello squawk" {
		t.Errorf("expected 'hello squawk', got: %q", result)
	}
}
