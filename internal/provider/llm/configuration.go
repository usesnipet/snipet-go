package llm

type Configuration[T any] struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Config   T      `json:"config"`
}
