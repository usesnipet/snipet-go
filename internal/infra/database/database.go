package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/usesnipet/snipet/config"
	"github.com/usesnipet/snipet/internal/logger"

	_ "ariga.io/atlas-provider-gorm/gormschema"
)

func NewDatabase(cfg *config.Config, logger *logger.Logger) (*gorm.DB, *sql.DB, error) {
	if err := ensureDatabase(cfg, logger); err != nil {
		return nil, nil, fmt.Errorf("ensure database: %w", err)
	}

	sqlDB, err := sql.Open("pgx", cfg.Database.URL)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB}),
		&gorm.Config{TranslateError: true, SkipDefaultTransaction: true},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	if err = runMigrations(cfg, logger); err != nil {
		return nil, nil, fmt.Errorf("run migrations: %w", err)
	}

	return gormDB, sqlDB, nil
}
