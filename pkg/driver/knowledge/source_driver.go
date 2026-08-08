package knowledge

import (
	"context"

	"github.com/usesnipet/snipet/pkg/driver"
	"github.com/usesnipet/snipet/pkg/jsonx"
)

// sourceDriver is the concrete ISourceDriver built by CreateSourceDriver,
// delegating each method to the corresponding configured func.
type sourceDriver struct {
	info driver.Info

	testConnection func(ctx context.Context, config jsonx.JSONMap) error
	iterator       IteratorConstructor
	reader         ReaderConstructor
}

func (d *sourceDriver) Info() driver.Info {
	return d.info
}

// Validate checks Info is well-formed and TestConnection, Iterator, and
// Reader are all configured. It's called by CreateSourceDriver and again by
// R.Register, so a driver missing any of these never enters a registry.
func (d *sourceDriver) Validate() error {
	if err := d.info.Validate(); err != nil {
		return err
	}
	if d.testConnection == nil {
		return ErrTestConnectionNotConfigured
	}
	if d.iterator == nil {
		return ErrIteratorNotConfigured
	}
	if d.reader == nil {
		return ErrReaderNotConfigured
	}
	return nil
}

func (d *sourceDriver) TestConnection(ctx context.Context, config jsonx.JSONMap) error {
	return d.testConnection(ctx, config)
}

func (d *sourceDriver) Iterator(ctx context.Context, config jsonx.JSONMap) (IKnowledgeIterator, error) {
	return d.iterator(ctx, config)
}

func (d *sourceDriver) Reader(ctx context.Context, config jsonx.JSONMap, itemID string) (IKnowledgeReader, error) {
	return d.reader(ctx, config, itemID)
}
