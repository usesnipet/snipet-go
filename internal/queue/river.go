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
	return insertJob(ctx, q.riverClient, args, opts)
}

// PushFromContext enqueues a job using the River client available in a worker's context.
func PushFromContext(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (int64, error) {
	client := river.ClientFromContext[*sql.Tx](ctx)
	return insertJob(ctx, client, args, opts)
}

func insertJob(ctx context.Context, client *river.Client[*sql.Tx], args river.JobArgs, opts *river.InsertOpts) (int64, error) {
	var result *rivertype.JobInsertResult
	var err error
	if tx, hasTx := repository.GetTx(ctx); hasTx {
		sqlTx := tx.Statement.ConnPool.(*sql.Tx)
		result, err = client.InsertTx(ctx, sqlTx, args, opts)
	} else {
		result, err = client.Insert(ctx, args, opts)
	}
	if err != nil {
		return 0, err
	}
	return result.Job.ID, nil
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
