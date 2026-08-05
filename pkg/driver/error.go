package driver

import "errors"

var (
	// ErrDriverNotFound is returned when a driver lookup (e.g. by Info.Key)
	// does not match any registered driver.
	ErrDriverNotFound = errors.New("driver not found")
	// ErrDriverConnectionFailed is returned when IDriver.TestConnection
	// fails to reach or authenticate against the underlying service.
	ErrDriverConnectionFailed = errors.New("driver connection failed")
)
