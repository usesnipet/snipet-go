package runtime

type IndexManager struct {
	*Manager[IIndexDriver]
}

func NewIndexManager(registry *Registry[IIndexDriver]) *IndexManager {
	return &IndexManager{Manager: NewManager(registry)}
}
