package rag

type Config struct {
	Milvus struct {
		Address     string `json:"address"`
		Dimension   int    `json:"dimension"`
		MaxPoolSize int    `json:"maxPoolSize"`
		Collection  string `json:"collection"`
	} `json:"milvus"`
}
