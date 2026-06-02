package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiURL = "https://slack.com/api/chat.postMessage"

type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// Token returns the Slack token from the environment.
// Checks SQUAWK_SLACK_TOKEN first, then SLACK_TOKEN.
func Token() (string, error) {
	if t := os.Getenv("SQUAWK_SLACK_TOKEN"); t != "" {
		return t, nil
	}
	if t := os.Getenv("SLACK_TOKEN"); t != "" {
		return t, nil
	}
	return "", fmt.Errorf("no slack token found: set SQUAWK_SLACK_TOKEN or SLACK_TOKEN")
}

// Send posts a message to a Slack channel.
func Send(token string, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshalling message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	var r response
	if err := json.Unmarshal(respBody, &r); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}

	if !r.OK {
		return fmt.Errorf("slack API error: %s", r.Error)
	}

	return nil
}
