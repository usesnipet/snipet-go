package fileutil

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// Hash returns the hex-encoded SHA-256 digest of r's full content.
// Used by source drivers for change detection between sync runs.
func Hash(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
