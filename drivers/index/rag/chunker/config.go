package chunker

import "fmt"

type Config struct {
	ChunkSize      int      `json:"chunkSize"`
	Overlap        int      `json:"overlap"`
	MinChunkSize   int      `json:"minChunkSize"`
	TrimWhitespace bool     `json:"trimWhitespace"`
	Separators     []string `json:"separators"`
}

func DefaultConfig() Config {
	return Config{
		ChunkSize:      1000,
		Overlap:        200,
		MinChunkSize:   0,
		TrimWhitespace: true,
	}
}

func (c Config) Validate() error {
	if c.ChunkSize <= 0 {
		return fmt.Errorf("chunker: chunk size must be greater than zero")
	}
	if c.Overlap < 0 {
		return fmt.Errorf("chunker: overlap cannot be negative")
	}
	if c.Overlap >= c.ChunkSize {
		return fmt.Errorf("chunker: overlap must be smaller than chunk size")
	}
	if c.MinChunkSize < 0 {
		return fmt.Errorf("chunker: min chunk size cannot be negative")
	}
	if c.MinChunkSize > c.ChunkSize {
		return fmt.Errorf("chunker: min chunk size cannot exceed chunk size")
	}
	return nil
}

func (c Config) separators() []string {
	if len(c.Separators) > 0 {
		return c.Separators
	}
	return defaultSeparators
}
