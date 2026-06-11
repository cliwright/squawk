package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenPrimary(t *testing.T) {
	t.Setenv("SQUAWK_SLACK_TOKEN", "primary-token")
	t.Setenv("SLACK_TOKEN", "fallback-token")

	tok, err := Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "primary-token" {
		t.Errorf("expected primary token, got %q", tok)
	}
}

func TestTokenFallback(t *testing.T) {
	t.Setenv("SQUAWK_SLACK_TOKEN", "")
	t.Setenv("SLACK_TOKEN", "fallback-token")

	tok, err := Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "fallback-token" {
		t.Errorf("expected fallback token, got %q", tok)
	}
}

func TestTokenMissing(t *testing.T) {
	t.Setenv("SQUAWK_SLACK_TOKEN", "")
	t.Setenv("SLACK_TOKEN", "")

	_, err := Token()
	if err == nil {
		t.Fatal("expected error when no token is set")
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage("#alerts", "hello world", "#CC0000")

	if msg.Channel != "#alerts" {
		t.Errorf("expected channel #alerts, got %q", msg.Channel)
	}
	if len(msg.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(msg.Attachments))
	}
	if msg.Attachments[0].Text != "hello world" {
		t.Errorf("expected text 'hello world', got %q", msg.Attachments[0].Text)
	}
	if msg.Attachments[0].Color != "#CC0000" {
		t.Errorf("expected color #CC0000, got %q", msg.Attachments[0].Color)
	}
}

func TestSendSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization header Bearer test-token, got %q", r.Header.Get("Authorization"))
		}

		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if msg.Channel != "#test" {
			t.Errorf("expected channel #test, got %q", msg.Channel)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{OK: true})
	}))
	defer server.Close()

	originalURL := apiURL
	// We can't easily patch the package-level apiURL constant, so this test
	// documents the expected behaviour against a real-looking response.
	_ = originalURL
	_ = server

	// Since apiURL is a const, we test the parsing logic with a mocked round-tripper
	// by replacing the default client temporarily.
	oldClient := http.DefaultClient
	defer func() { http.DefaultClient = oldClient }()

	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != apiURL {
				return nil, fmt.Errorf("unexpected URL: %s", req.URL.String())
			}
			rec := httptest.NewRecorder()
			rec.Header().Set("Content-Type", "application/json")
			rec.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(rec).Encode(response{OK: true})
			return rec.Result(), nil
		}),
	}

	msg := NewMessage("#test", "hello", "")
	if err := Send("test-token", msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendAPIError(t *testing.T) {
	oldClient := http.DefaultClient
	defer func() { http.DefaultClient = oldClient }()

	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			rec.Header().Set("Content-Type", "application/json")
			rec.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(rec).Encode(response{OK: false, Error: "channel_not_found"})
			return rec.Result(), nil
		}),
	}

	msg := NewMessage("#bad", "hello", "")
	err := Send("test-token", msg)
	if err == nil {
		t.Fatal("expected error for failed slack API response")
	}
	if err.Error() != "slack API error: channel_not_found" {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
