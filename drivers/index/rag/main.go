package rag

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/usesnipet/snipet/drivers/index/rag/store"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

//go:embed schema.json
var schemaJSON []byte

type Driver struct{}

func NewDriver() runtime.IIndexDriver {
	return &Driver{}
}

func (d *Driver) GetConfigurationSchema(ctx context.Context) (util.JSONMap, error) {
	var schema util.JSONMap
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return nil, fmt.Errorf("rag: parse schema: %w", err)
	}
	return schema, nil
}

func (d *Driver) TestConnection(ctx context.Context, configJson util.JSONMap) error {
	cfg, err := util.ParseJSONMap[Config](configJson)
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	s, err := store.NewStore(cfg.StoreConfig())
	if err != nil {
		return err
	}
	if err := s.Start(ctx); err != nil {
		return err
	}
	defer s.Close()
	return s.Ping(ctx)
}

func (d *Driver) Reader(config util.JSONMap) (runtime.IIndexReader, error) {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return nil, err
	}
	return NewReader(context.Background(), cfg)
}

func (d *Driver) Writer(config util.JSONMap) (runtime.IIndexWriter, error) {
	cfg, err := util.ParseJSONMap[Config](config)
	if err != nil {
		return nil, err
	}
	return NewWriter(context.Background(), cfg)
}
