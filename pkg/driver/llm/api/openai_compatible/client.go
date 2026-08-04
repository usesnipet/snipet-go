package openaicompatible

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// httpClient is shared across requests to reuse connections; it carries no
// per-call state.
var httpClient = &http.Client{
	Timeout: 5 * time.Minute, // TODO: make this configurable
}

// chatCompletionsURL builds the /chat/completions endpoint under baseURL.
func chatCompletionsURL(baseURL string) string {
	return baseURL + "/chat/completions"
}

// doChatCompletion sends a non-streaming chat completions request and
// decodes the JSON response, surfacing both API-level errors (parsed.Error)
// and HTTP-level errors (non-2xx status).
func doChatCompletion(ctx context.Context, defaultBaseURL string, cfg Config, body chatRequest) (*chatResponse, error) {
	baseURL, err := resolveBaseURL(defaultBaseURL, cfg)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionsURL(baseURL), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	setHeaders(req, cfg)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chat completions request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var parsed chatResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("decode response (status %d): %w: %s", res.StatusCode, err, truncate(string(data), 512))
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("api error: %w", parsed.Error)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("chat completions returned %d: %s", res.StatusCode, truncate(string(data), 512))
	}
	return &parsed, nil
}

// openChatCompletionStream opens a Server-Sent Events stream for a chat
// completions request. The caller owns the returned response body and must
// close it; on a non-2xx status the body is drained and closed here instead,
// and an error is returned.
func openChatCompletionStream(ctx context.Context, defaultBaseURL string, cfg Config, body chatRequest) (*http.Response, error) {
	baseURL, err := resolveBaseURL(defaultBaseURL, cfg)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, chatCompletionsURL(baseURL), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	setHeaders(req, cfg)
	req.Header.Set("Accept", "text/event-stream")

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chat completions stream request: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		defer res.Body.Close()
		data, _ := io.ReadAll(res.Body)
		var parsed chatResponse
		if err := json.Unmarshal(data, &parsed); err == nil && parsed.Error != nil {
			return nil, fmt.Errorf("api error: %w", parsed.Error)
		}
		return nil, fmt.Errorf("chat completions stream returned %d: %s", res.StatusCode, truncate(string(data), 512))
	}
	return res, nil
}

// setHeaders applies the standard JSON content headers and, if configured,
// a Bearer authorization header.
func setHeaders(req *http.Request, cfg Config) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	}
}

// truncate shortens s to at most n bytes, appending "..." when it was cut,
// so error messages don't embed unbounded response bodies.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
