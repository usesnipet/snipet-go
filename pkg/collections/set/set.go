package set

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) Size() int {
	return len(s)
}

func New[T comparable](values ...T) Set[T] {
	s := make(Set[T])
	for _, value := range values {
		s.Add(value)
	}
	return s
}
