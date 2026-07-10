package rag

import (
	"github.com/usesnipet/snipet/drivers/index/rag/chunker"
	"github.com/usesnipet/snipet/drivers/index/rag/milvus"
)

type Config struct {
	Milvus   milvus.Config  `json:"milvus"`
	Chunker  chunker.Config `json:"chunker"`
	Embedder struct {
		Provider string `json:"provider"`
		Model    string `json:"model"`
		APIKey   string `json:"apiKey"`
	} `json:"embedder"`
}

func (c Config) ChunkerConfig() chunker.Config {
	cfg := c.Chunker
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = chunker.DefaultConfig().ChunkSize
	}
	return cfg
}
