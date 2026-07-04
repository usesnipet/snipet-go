package rag

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/teilomillet/raggo"
	"github.com/usesnipet/snipet/internal/runtime"
	"github.com/usesnipet/snipet/internal/util"
)

//go:embed schema.json
var schemaJSON []byte

type Driver struct {
}

func NewDriver() runtime.IIndexDriver {
	return &Driver{}
}

func (d *Driver) GetConfigurationSchema(ctx context.Context) (util.JSONMap, error) {
	var schema util.JSONMap
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		return nil, fmt.Errorf("fs: parse schema: %w", err)
	}
	return schema, nil
}

func (d *Driver) TestConnection(ctx context.Context, configJson util.JSONMap) error {
	cfg, err := util.ParseJSONMap[Config](configJson)
	if err != nil {
		return err
	}

	vectorDB, err := raggo.NewVectorDB(
		raggo.WithAddress(cfg.Milvus.Address),
		raggo.WithDimension(cfg.Milvus.Dimension),
		raggo.WithMaxPoolSize(cfg.Milvus.MaxPoolSize),
	)
	if err != nil {
		return err
	}

	err = vectorDB.Connect(ctx)
	defer vectorDB.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) Reader(config util.JSONMap) runtime.IIndexReader {
}

func (d *Driver) Writer(config util.JSONMap) runtime.IIndexWriter {
}
