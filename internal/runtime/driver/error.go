package driver

import "errors"

var (
	ErrDriverNotFound         = errors.New("driver not found")
	ErrDriverConnectionFailed = errors.New("driver connection failed")
)
