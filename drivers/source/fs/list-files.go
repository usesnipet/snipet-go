package fs

import (
	"io/fs"
	"path/filepath"
)

func listFiles(cfg Config) ([]string, error) {
	var files []string

	matcher := NewMatcher(cfg.Ignore)

	err := filepath.WalkDir(cfg.BasePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(cfg.BasePath, path)
		if err != nil {
			return err
		}

		rel = filepath.ToSlash(rel)

		if matcher.Match(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
