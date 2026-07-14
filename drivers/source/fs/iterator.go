package fs

import (
	"context"
	"time"

	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

type Item struct {
	Path         string
	Name         string
	Metadata     util.JSONMap
	LastModified *time.Time
	Hash         string
}

type Iterator struct {
	files []string

	current *runtime.SourceItem
	hash    string
	err     error
	index   int
}

func NewIterator(files []string) runtime.ISourceIterator {
	return &Iterator{
		files:   files,
		current: nil,
		err:     nil,
		index:   0,
	}
}

func (it *Iterator) Next(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	it.index++

	if it.index >= len(it.files) {
		return false
	}

	path := it.files[it.index]

	sourceItem, hash, err := sourceItemFromFile(path)
	if err != nil {
		it.err = err
		return false
	}
	it.current = sourceItem
	it.hash = hash

	return true
}

func (i *Iterator) Item() *runtime.SourceItem {
	return i.current
}

func (i *Iterator) GetHash() string {
	return i.hash
}

func (i *Iterator) Err() error {
	return i.err
}

func (i *Iterator) Close() error {
	i.current = nil
	i.err = nil
	return nil
}
