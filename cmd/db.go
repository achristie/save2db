package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/achristie/save2db/internal/pg"
	"github.com/achristie/save2db/internal/sqlite"
)

type Database interface {
	Open() error
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	GetDB() *sql.DB
}

func (app *application) GetTx(cfg Config) (*sql.Tx, error) {
	ctx := context.Background()

	switch cfg.Database.Name {
	case "PostgreSQL":
		db = pg.NewDB(cfg.Database.DSN)
	default:
		db = sqlite.NewDB(cfg.Database.DSN)
	}

	if err := db.Open(); err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	// begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("db tx: %w", err)
	}

	return tx, nil
}
