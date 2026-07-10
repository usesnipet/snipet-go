package milvus

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type Config struct {
	Address        string `json:"address"`
	Dimension      int    `json:"dimension"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	DBName         string `json:"dbName"`
	APIKey         string `json:"apiKey"`
	CollectionName string `json:"collectionName"`
}

type Milvus struct {
	cfg     Config
	client  *milvusclient.Client
	started bool
}

func New(cfg Config) (*Milvus, error) {
	return &Milvus{cfg: cfg}, nil
}

func (m *Milvus) Start(ctx context.Context) error {
	if m.cfg.CollectionName == "" {
		return fmt.Errorf("milvus: collection name is required")
	}
	if m.cfg.Dimension <= 0 {
		return fmt.Errorf("milvus: dimension must be greater than zero")
	}

	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address:  m.cfg.Address,
		Username: m.cfg.Username,
		Password: m.cfg.Password,
		DBName:   m.cfg.DBName,
		APIKey:   m.cfg.APIKey,
	})
	if err != nil {
		return fmt.Errorf("milvus: connect: %w", err)
	}

	m.client = client

	if err := m.ensureCollection(ctx); err != nil {
		return err
	}
	if err := m.ensureLoaded(ctx); err != nil {
		return err
	}

	m.started = true
	return nil
}

func (m *Milvus) ensureCollection(ctx context.Context) error {
	has, err := m.client.HasCollection(ctx, milvusclient.NewHasCollectionOption(m.cfg.CollectionName))
	if err != nil {
		return fmt.Errorf("milvus: check collection: %w", err)
	}
	if has {
		return nil
	}

	schema := SourceSchema(m.cfg.CollectionName, m.cfg.Dimension)
	err = m.client.CreateCollection(ctx,
		milvusclient.NewCreateCollectionOption(m.cfg.CollectionName, schema).
			WithIndexOptions(SourceIndexSchema(m.cfg.CollectionName)...),
	)
	if err != nil {
		return fmt.Errorf("milvus: create collection: %w", err)
	}

	return nil
}

func (m *Milvus) ensureLoaded(ctx context.Context) error {
	loadState, err := m.client.GetLoadState(ctx, milvusclient.NewGetLoadStateOption(m.cfg.CollectionName))
	if err != nil {
		return fmt.Errorf("milvus: get load state: %w", err)
	}
	if loadState.State == entity.LoadStateLoaded {
		return nil
	}

	task, err := m.client.LoadCollection(ctx, milvusclient.NewLoadCollectionOption(m.cfg.CollectionName))
	if err != nil {
		return fmt.Errorf("milvus: load collection: %w", err)
	}
	if err := task.Await(ctx); err != nil {
		return fmt.Errorf("milvus: await load collection: %w", err)
	}

	return nil
}

func (m *Milvus) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return m.client.Close(ctx)
}
