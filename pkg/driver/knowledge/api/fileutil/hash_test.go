package fileutil

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashReturnsSHA256HexDigest(t *testing.T) {
	content := "hello world"
	expected := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	hash, err := Hash(strings.NewReader(content))

	require.NoError(t, err)
	assert.Equal(t, expected, hash)
}

func TestHashIsStableAcrossCalls(t *testing.T) {
	h1, err := Hash(strings.NewReader("same content"))
	require.NoError(t, err)

	h2, err := Hash(strings.NewReader("same content"))
	require.NoError(t, err)

	assert.Equal(t, h1, h2)
}

func TestHashDiffersForDifferentContent(t *testing.T) {
	h1, err := Hash(strings.NewReader("content a"))
	require.NoError(t, err)

	h2, err := Hash(strings.NewReader("content b"))
	require.NoError(t, err)

	assert.NotEqual(t, h1, h2)
}

func TestHashEmptyReader(t *testing.T) {
	expected := fmt.Sprintf("%x", sha256.Sum256(nil))

	hash, err := Hash(strings.NewReader(""))

	require.NoError(t, err)
	assert.Equal(t, expected, hash)
}
