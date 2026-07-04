package queue

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/usesnipet/snipet/internal/repository"
)

type RiverQueue struct {
	riverClient *river.Client[*sql.Tx]
}

func NewRiver(sqlDB *sql.DB, workers *river.Workers) (*RiverQueue, error) {
	riverClient, err := river.NewClient(
		riverdatabasesql.New(sqlDB),
		&river.Config{
			Workers: workers,
		},
	)
	return &RiverQueue{riverClient: riverClient}, err
}

func (q *RiverQueue) Push(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) error {
	var err error
	if tx, hasTx := repository.GetTx(ctx); hasTx {
		sqlTx := tx.Statement.ConnPool.(*sql.Tx)
		_, err = q.riverClient.InsertTx(ctx, sqlTx, args, opts)
	} else {
		_, err = q.riverClient.Insert(ctx, args, opts)
	}
	return err
}

func (q *RiverQueue) Start(ctx context.Context) error {
	return q.riverClient.Start(ctx)
}

func (q *RiverQueue) Stop(ctx context.Context) error {
	return q.riverClient.Stop(ctx)
}
