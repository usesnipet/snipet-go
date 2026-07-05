package rag

type Config struct {
	Milvus struct {
		Address     string `json:"address"`
		Dimension   int    `json:"dimension"`
		MaxPoolSize int    `json:"maxPoolSize"`
		Collection  string `json:"collection"`
	} `json:"milvus"`
	ChunkSize int `json:"chunkSize"`
	Embedder  struct {
		Provider string `json:"provider"`
		Model    string `json:"model"`
		APIKey   string `json:"apiKey"`
	} `json:"embedder"`
}
