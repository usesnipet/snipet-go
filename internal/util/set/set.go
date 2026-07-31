package set

type Set[T comparable] map[T]struct{}

func new[T comparable]() Set[T] {
	return make(Set[T])
}

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
	s := new[T]()
	for _, value := range values {
		s.Add(value)
	}
	return s
}
