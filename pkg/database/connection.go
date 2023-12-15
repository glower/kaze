package database

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/glower/kaze/pkg/config"
)

func NewConnection(cfg *config.Config) (*sqlx.DB, error) {
	db, err := setupDB(cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func setupDB(cfg *config.Config) (*sqlx.DB, error) {
	if cfg.DB == "" {
		return nil, fmt.Errorf("please provide database connection string via export APP_DB=...")
	}

	// Connect to the PostgreSQL database
	db, err := sqlx.Connect("postgres", cfg.DB)
	if err != nil {
		slog.Error("can't connect to database", "error", err, "dataSourceName", cfg.DB)
		return nil, err
	}

	// Ping the database to ensure connection is established
	if err := db.Ping(); err != nil {
		slog.Error("can't ping database", "error", err, "dataSourceName", cfg.DB)
		return nil, err
	}

	return db, nil
}
