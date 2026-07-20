package registry

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type R[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

func New[T any]() *R[T] {
	return &R[T]{
		items: make(map[string]T),
	}
}

func (r *R[T]) Register(name string, value T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[name]; exists {
		return fmt.Errorf("registry: %q already registered", name)
	}

	r.items[name] = value
	return nil
}

func (r *R[T]) MustRegister(name string, value T) {
	if err := r.Register(name, value); err != nil {
		panic(err)
	}
}

func (r *R[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	value, ok := r.items[name]
	return value, ok
}

func (r *R[T]) MustGet(name string) T {
	value, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("registry: %q not found", name))
	}
	return value
}

func (r *R[T]) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.items[name]
	return ok
}

func (r *R[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := slices.Collect(maps.Keys(r.items))
	slices.Sort(names)

	return names
}
