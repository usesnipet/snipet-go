package runtime

import "errors"

var (
	ErrSourceDriverNotFound = errors.New("source driver not found")
	ErrInvalidConfiguration = errors.New("invalid configuration")
)
