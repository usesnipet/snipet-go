package queue

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

type IJobQueue interface {
	Push(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (int64, error)
	JobGet(ctx context.Context, id int64) (*rivertype.JobRow, error)
}
