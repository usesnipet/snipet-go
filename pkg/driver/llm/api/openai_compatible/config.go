package openaicompatible

import (
	"fmt"
	"strings"

	"github.com/usesnipet/snipet/internal/util"
)

// Config holds OpenAI-compatible provider settings from the driver configuration schema.
type Config struct {
	APIKey      string  `json:"api_key"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	TopP        float64 `json:"top_p"`
	Endpoint    string  `json:"endpoint"`
}

// NewConfig parses and validates a driver config map into a Config. It fails
// if the map doesn't match Config's shape or required fields (e.g. Model)
// are missing.
func NewConfig(config util.JSONMap) (Config, error) {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}
	err = cfg.validate()
	return cfg, err
}

func (c Config) validate() error {
	if c.Model == "" {
		return fmt.Errorf("model is required")
	}
	return nil
}

// resolveBaseURL prefers a runtime endpoint override from config, otherwise
// uses the base URL injected when the API was created.
func resolveBaseURL(defaultBaseURL string, cfg Config) (string, error) {
	base := cfg.Endpoint
	if base == "" {
		base = defaultBaseURL
	}
	base = strings.TrimRight(base, "/")
	if base == "" {
		return "", fmt.Errorf("base url is required")
	}
	return base, nil
}
