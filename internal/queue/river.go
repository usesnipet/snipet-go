package queue

import (
	"database/sql"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
)

func NewRiver(sqlDB *sql.DB, workers *river.Workers) (*river.Client[*sql.Tx], error) {
	riverClient, err := river.NewClient(
		riverdatabasesql.New(sqlDB),
		&river.Config{
			Workers: workers,
		},
	)
	return riverClient, err
}
