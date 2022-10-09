package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	_ "modernc.org/sqlite"
)

//go:embed migration/*.sql
var migrationFS embed.FS

type DB struct {
	db     *sql.DB
	ctx    context.Context
	cancel func()
	source string
}

func (db *DB) GetDB() *sql.DB {
	return db.db
}

// Simple for now..
func NewDB(source string) *DB {
	db := &DB{
		source: source,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() (err error) {
	if db.source == "" {
		return fmt.Errorf("datasource required")
	}

	if err := os.MkdirAll(filepath.Dir(db.source), 0700); err != nil {
		return err
	}

	if db.db, err = sql.Open("sqlite", db.source); err != nil {
		return err
	}

	if _, err := db.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	if err := db.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func (db *DB) migrate() error {
	if _, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	names, err := fs.Glob(migrationFS, "migration/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(names)

	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q, err=%q", name, err)
		}
	}
	return nil
}

func (db *DB) migrateFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil
	}

	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES(?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
