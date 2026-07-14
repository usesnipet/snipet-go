package embedder

import "fmt"

type Config struct {
	Model  string `json:"model"`
	APIKey string `json:"apiKey"`
}

func (c Config) Validate() error {
	if c.Model == "" {
		return fmt.Errorf("embedder: model is required")
	}
	return nil
}
