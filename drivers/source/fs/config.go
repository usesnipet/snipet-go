package fs

type Config struct {
	BasePath string   `json:"basePath"`
	Ignore   []string `json:"ignore"`
}
