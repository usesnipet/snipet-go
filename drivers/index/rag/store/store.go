package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
	"github.com/usesnipet/snipet/drivers/index/rag/chunker"
	"github.com/usesnipet/snipet/drivers/index/rag/embedder"
)

type Hit struct {
	IndexedItemID string
	Chunk         chunker.Chunk
	Embedding     []float32
	Metadata      map[string]any
	Distance      float64
}

type Store struct {
	cfg     Config
	pool    *pgxpool.Pool
	started bool
}

func NewStore(cfg Config) (*Store, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Store{cfg: cfg}, nil
}

func (s *Store) Start(ctx context.Context) error {
	if s.started {
		return nil
	}

	poolCfg, err := pgxpool.ParseConfig(s.cfg.URL)
	if err != nil {
		return fmt.Errorf("store: parse url: %w", err)
	}
	poolCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return fmt.Errorf("store: connect: %w", err)
	}

	s.pool = pool
	if err := s.ensureSchema(ctx); err != nil {
		pool.Close()
		s.pool = nil
		return err
	}

	s.started = true
	return nil
}

func (s *Store) Ping(ctx context.Context) error {
	if err := s.requireStarted(); err != nil {
		return err
	}
	return s.pool.Ping(ctx)
}

func (s *Store) Close() {
	if s.pool != nil {
		s.pool.Close()
		s.pool = nil
	}
	s.started = false
}

// Replace deletes existing rows for indexedItemID and inserts the new embeddings atomically.
func (s *Store) Replace(ctx context.Context, indexedItemID string, embeddings []*embedder.Embedding) error {
	if err := s.requireStarted(); err != nil {
		return err
	}
	if indexedItemID == "" {
		return fmt.Errorf("store: indexed item id is required")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("store: begin replace: %w", err)
	}
	defer tx.Rollback(ctx)

	table := pgx.Identifier{s.cfg.Table}.Sanitize()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE indexed_item_id = $1`, table), indexedItemID); err != nil {
		return fmt.Errorf("store: delete before replace: %w", err)
	}

	if err := s.copyEmbeddings(ctx, tx, indexedItemID, embeddings); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("store: commit replace: %w", err)
	}
	return nil
}

func (s *Store) Index(ctx context.Context, indexedItemID string, embeddings []*embedder.Embedding) error {
	if err := s.requireStarted(); err != nil {
		return err
	}
	if indexedItemID == "" {
		return fmt.Errorf("store: indexed item id is required")
	}
	return s.copyEmbeddings(ctx, s.pool, indexedItemID, embeddings)
}

func (s *Store) Search(ctx context.Context, query []float32, limit int) ([]Hit, error) {
	if err := s.requireStarted(); err != nil {
		return nil, err
	}
	if len(query) != s.cfg.Length {
		return nil, fmt.Errorf("store: query length %d does not match configured length %d", len(query), s.cfg.Length)
	}
	if limit <= 0 {
		limit = 10
	}

	table := pgx.Identifier{s.cfg.Table}.Sanitize()
	sql := fmt.Sprintf(`
		SELECT indexed_item_id, seq_id, content, start_offset, end_offset, embedding, metadata,
		       embedding <=> $1 AS distance
		FROM %s
		ORDER BY embedding <=> $1
		LIMIT $2
	`, table)

	rows, err := s.pool.Query(ctx, sql, pgvector.NewVector(query), limit)
	if err != nil {
		return nil, fmt.Errorf("store: search: %w", err)
	}
	defer rows.Close()

	hits := make([]Hit, 0, limit)
	for rows.Next() {
		var (
			hit      Hit
			vector   pgvector.Vector
			metadata []byte
		)
		if err := rows.Scan(
			&hit.IndexedItemID,
			&hit.Chunk.SeqID,
			&hit.Chunk.Content,
			&hit.Chunk.Start,
			&hit.Chunk.End,
			&vector,
			&metadata,
			&hit.Distance,
		); err != nil {
			return nil, fmt.Errorf("store: search scan: %w", err)
		}
		hit.Embedding = vector.Slice()
		if len(metadata) > 0 {
			if err := json.Unmarshal(metadata, &hit.Metadata); err != nil {
				return nil, fmt.Errorf("store: unmarshal metadata: %w", err)
			}
		}
		if hit.Metadata == nil {
			hit.Metadata = map[string]any{}
		}
		hit.Chunk.Metadata = hit.Metadata
		hits = append(hits, hit)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: search rows: %w", err)
	}
	return hits, nil
}

func (s *Store) DeleteMany(ctx context.Context, indexedItemIDs []string) error {
	if err := s.requireStarted(); err != nil {
		return err
	}
	if len(indexedItemIDs) == 0 {
		return nil
	}

	table := pgx.Identifier{s.cfg.Table}.Sanitize()
	sql := fmt.Sprintf(`DELETE FROM %s WHERE indexed_item_id = ANY($1)`, table)
	_, err := s.pool.Exec(ctx, sql, indexedItemIDs)
	if err != nil {
		return fmt.Errorf("store: delete many: %w", err)
	}
	return nil
}

type copySource interface {
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

func (s *Store) copyEmbeddings(ctx context.Context, dest copySource, indexedItemID string, embeddings []*embedder.Embedding) error {
	if len(embeddings) == 0 {
		return nil
	}

	rows := make([][]any, 0, len(embeddings))
	for _, emb := range embeddings {
		if emb == nil {
			continue
		}
		if len(emb.Embedding) != s.cfg.Length {
			return fmt.Errorf("store: embedding length %d does not match configured length %d", len(emb.Embedding), s.cfg.Length)
		}
		metadata, err := json.Marshal(emb.Metadata)
		if err != nil {
			return fmt.Errorf("store: marshal metadata: %w", err)
		}
		if metadata == nil {
			metadata = []byte("{}")
		}
		rows = append(rows, []any{
			uuid.NewString(),
			indexedItemID,
			emb.Chunk.SeqID,
			emb.Chunk.Content,
			emb.Chunk.Start,
			emb.Chunk.End,
			pgvector.NewVector(emb.Embedding),
			metadata,
		})
	}
	if len(rows) == 0 {
		return nil
	}

	_, err := dest.CopyFrom(
		ctx,
		pgx.Identifier{s.cfg.Table},
		[]string{
			"id",
			"indexed_item_id",
			"seq_id",
			"content",
			"start_offset",
			"end_offset",
			"embedding",
			"metadata",
		},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("store: index: %w", err)
	}
	return nil
}

func (s *Store) ensureSchema(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector`); err != nil {
		return fmt.Errorf("store: create extension: %w", err)
	}

	table := pgx.Identifier{s.cfg.Table}.Sanitize()
	createTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id UUID PRIMARY KEY,
			indexed_item_id UUID NOT NULL,
			seq_id INT NOT NULL,
			content TEXT NOT NULL,
			start_offset INT NOT NULL,
			end_offset INT NOT NULL,
			embedding vector(%d) NOT NULL,
			metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`, table, s.cfg.Length)
	if _, err := s.pool.Exec(ctx, createTable); err != nil {
		return fmt.Errorf("store: create table: %w", err)
	}

	itemIndex := fmt.Sprintf(
		`CREATE INDEX IF NOT EXISTS %s ON %s (indexed_item_id)`,
		pgx.Identifier{s.cfg.Table + "_indexed_item_id_idx"}.Sanitize(),
		table,
	)
	if _, err := s.pool.Exec(ctx, itemIndex); err != nil {
		return fmt.Errorf("store: create item index: %w", err)
	}

	embeddingIndex := fmt.Sprintf(
		`CREATE INDEX IF NOT EXISTS %s ON %s USING hnsw (embedding vector_cosine_ops)`,
		pgx.Identifier{s.cfg.Table + "_embedding_hnsw_idx"}.Sanitize(),
		table,
	)
	if _, err := s.pool.Exec(ctx, embeddingIndex); err != nil {
		return fmt.Errorf("store: create embedding index: %w", err)
	}

	return nil
}

func (s *Store) requireStarted() error {
	if !s.started || s.pool == nil {
		return fmt.Errorf("store: not started")
	}
	return nil
}
