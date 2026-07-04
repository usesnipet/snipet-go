package runtime

type SourceManager struct {
	*Manager[ISourceDriver]
}

func NewSourceManager(registry *Registry[ISourceDriver]) *SourceManager {
	return &SourceManager{Manager: NewManager(registry)}
}
