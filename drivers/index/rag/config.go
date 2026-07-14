package rag

import (
	"github.com/usesnipet/snipet/drivers/index/rag/chunker"
	"github.com/usesnipet/snipet/drivers/index/rag/embedder"
	"github.com/usesnipet/snipet/drivers/index/rag/store"
)

type Config struct {
	Store    store.Config    `json:"store"`
	Chunker  chunker.Config  `json:"chunker"`
	Embedder embedder.Config `json:"embedder"`
}

func (c Config) Validate() error {
	if err := c.Store.Validate(); err != nil {
		return err
	}
	if err := c.ChunkerConfig().Validate(); err != nil {
		return err
	}
	if err := c.Embedder.Validate(); err != nil {
		return err
	}
	return nil
}

func (c Config) ChunkerConfig() chunker.Config {
	cfg := c.Chunker
	def := chunker.DefaultConfig()
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = def.ChunkSize
	}
	if cfg.Overlap <= 0 {
		cfg.Overlap = def.Overlap
	}
	if !cfg.TrimWhitespace && c.Chunker.ChunkSize <= 0 {
		cfg.TrimWhitespace = def.TrimWhitespace
	}
	return cfg
}

func (c Config) EmbedderConfig() embedder.Config {
	return c.Embedder
}

func (c Config) StoreConfig() store.Config {
	return c.Store
}
