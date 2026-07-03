package fs

import (
	"io/fs"
	"path/filepath"
)

func listFiles(basePath string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
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
