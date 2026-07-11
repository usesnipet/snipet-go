package fs

import (
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

type Matcher struct {
	patterns []string
}

func NewMatcher(patterns []string) *Matcher {
	return &Matcher{
		patterns: patterns,
	}
}

func (m *Matcher) Match(path string) bool {
	path = filepath.ToSlash(path)

	for _, pattern := range m.patterns {
		ok, err := doublestar.Match(pattern, path)
		if err != nil {
			continue
		}

		if ok {
			return true
		}
	}

	return false
}
