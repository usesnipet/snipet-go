package runtime

import "errors"

var (
	ErrDriverNotFound       = errors.New("driver not found")
	ErrInvalidConfiguration = errors.New("invalid configuration")
	ErrConnectionFailed     = errors.New("connection failed")
)
