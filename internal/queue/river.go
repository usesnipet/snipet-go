package queue

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/riverqueue/river/rivertype"
	"github.com/usesnipet/snipet/internal/repository"
)

type RiverQueue struct {
	riverClient *river.Client[*sql.Tx]
}

func NewRiver(sqlDB *sql.DB, workers *river.Workers) (*RiverQueue, error) {
	// migrate
	migrator, err := rivermigrate.New(riverdatabasesql.New(sqlDB), nil)
	if err != nil {
		return nil, err
	}
	migrator.Migrate(context.Background(), rivermigrate.DirectionUp, nil)

	riverClient, err := river.NewClient(
		riverdatabasesql.New(sqlDB),
		&river.Config{
			Workers: workers,
			Queues: map[string]river.QueueConfig{
				river.QueueDefault: {
					MaxWorkers: 10,
				},
			},
		},
	)
	return &RiverQueue{riverClient: riverClient}, err
}

func (q *RiverQueue) Push(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (int64, error) {
	var result *rivertype.JobInsertResult
	var err error
	if tx, hasTx := repository.GetTx(ctx); hasTx {
		sqlTx := tx.Statement.ConnPool.(*sql.Tx)
		result, err = q.riverClient.InsertTx(ctx, sqlTx, args, opts)
	} else {
		result, err = q.riverClient.Insert(ctx, args, opts)
	}
	return result.Job.ID, err
}

func (q *RiverQueue) JobGet(ctx context.Context, id int64) (*rivertype.JobRow, error) {
	return q.riverClient.JobGet(ctx, id)
}

func (q *RiverQueue) Start(ctx context.Context) error {
	return q.riverClient.Start(ctx)
}

func (q *RiverQueue) Stop(ctx context.Context) error {
	return q.riverClient.Stop(ctx)
}
