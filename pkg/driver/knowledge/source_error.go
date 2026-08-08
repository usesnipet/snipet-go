package knowledge

import "errors"

var (
	// ErrTestConnectionNotConfigured is returned by TestConnection when the
	// driver was built without WithTestConnection.
	ErrTestConnectionNotConfigured = errors.New("test connection not configured")
	// ErrIteratorNotConfigured is returned by Iterator when the driver was
	// built without WithIterator.
	ErrIteratorNotConfigured = errors.New("iterator not configured")
	// ErrReaderNotConfigured is returned by Reader when the driver was built
	// without WithReader.
	ErrReaderNotConfigured = errors.New("reader not configured")
)
