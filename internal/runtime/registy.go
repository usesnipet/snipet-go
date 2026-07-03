package runtime

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type Registry[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		items: make(map[string]T),
	}
}

func (r *Registry[T]) Register(name string, value T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[name]; exists {
		return fmt.Errorf("registry: %q already registered", name)
	}

	r.items[name] = value
	return nil
}

func (r *Registry[T]) MustRegister(name string, value T) {
	if err := r.Register(name, value); err != nil {
		panic(err)
	}
}

func (r *Registry[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.items[name]
	return value, ok
}

func (r *Registry[T]) MustGet(name string) T {
	value, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("registry: %q not found", name))
	}
	return value
}

func (r *Registry[T]) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.items[name]
	return ok
}

func (r *Registry[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := slices.Collect(maps.Keys(r.items))
	slices.Sort(names)

	return names
}
