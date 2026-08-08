package tool

import "errors"

var (
	// ErrTestConnectionNotConfigured is returned by CreateDriver when the
	// driver was built without API.TestConnection.
	ErrTestConnectionNotConfigured = errors.New("test connection not configured")
	// ErrCallNotConfigured is returned by CreateDriver when the driver was
	// built without API.Call.
	ErrCallNotConfigured = errors.New("call not configured")
)
