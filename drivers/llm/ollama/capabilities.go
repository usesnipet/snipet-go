package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/usesnipet/snipet/internal/util"
	llm "github.com/usesnipet/snipet/pkg/driver/llm"
	openaicompatible "github.com/usesnipet/snipet/pkg/driver/llm/openai_compatible"
)

var showHTTPClient = &http.Client{Timeout: 30 * time.Second}

// showRequest/showResponse mirror the subset of Ollama's native
// POST /api/show response used here. Every locally installed model reports
// its own capabilities, so there's no list to maintain as new models are
// pulled.
type showRequest struct {
	Model string `json:"model"`
}

type showResponse struct {
	Capabilities []string `json:"capabilities"`
}

// capabilities reports the configured model's capabilities via Ollama's
// native /api/show endpoint, not the OpenAI-compatible surface used for
// generation, resolving the base URL the same way defaultBaseURL/Config do.
func capabilities(ctx context.Context, defaultBaseURL string, config util.JSONMap) (llm.Capabilities, error) {
	cfg, err := openaicompatible.NewConfig(config)
	if err != nil {
		return llm.Capabilities{}, err
	}

	base := cfg.Endpoint
	if base == "" {
		base = defaultBaseURL
	}
	base = strings.TrimSuffix(strings.TrimRight(base, "/"), "/v1")

	payload, err := json.Marshal(showRequest{Model: cfg.Model})
	if err != nil {
		return llm.Capabilities{}, fmt.Errorf("marshal show request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/show", bytes.NewReader(payload))
	if err != nil {
		return llm.Capabilities{}, fmt.Errorf("create show request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := showHTTPClient.Do(req)
	if err != nil {
		return llm.Capabilities{}, fmt.Errorf("show request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return llm.Capabilities{}, fmt.Errorf("read show response: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return llm.Capabilities{}, fmt.Errorf("show returned %d: %s", res.StatusCode, data)
	}

	var parsed showResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return llm.Capabilities{}, fmt.Errorf("decode show response: %w", err)
	}

	if slices.Contains(parsed.Capabilities, "tools") {
		return llm.Capabilities{ToolCall: true}, nil
	}
	return llm.Capabilities{ToolCall: false}, nil
}
