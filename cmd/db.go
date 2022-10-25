package cmd

import (
	"database/sql"
	"fmt"

	"github.com/achristie/save2db/internal/pg"
	"github.com/achristie/save2db/internal/sqlite"
)

type Database interface {
	Open() error
	BeginTx(*sql.TxOptions) (*sql.Tx, error)
	GetDB() *sql.DB
}

func (app *application) GetTx(cfg Config) (*sql.Tx, error) {
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
	tx, err := db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("db tx: %w", err)
	}

	return tx, nil
}
