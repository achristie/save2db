package pg

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed migration/*.sql
var migrationFS embed.FS

type DB struct {
	Db     *sql.DB
	Ctx    context.Context
	Cancel func()
	Source string
}

// func NewDB(source string) *DB {
// 	db := &DB{source: source}
// 	db.ctx, db.cancel = context.WithCancel(context.Background())
// 	return db
// }

func (db *DB) Open() (err error) {
	if db.Source == "" {
		return fmt.Errorf("datasource required")
	}
	if db.Db, err = sql.Open("pgx", db.Source); err != nil {
		return err
	}

	if err := db.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func (db *DB) migrate() error {
	if _, err := db.Db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name CHAR(18) PRIMARY KEY);`); err != nil {
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
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = $1`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil
	}

	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES($1)`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	tx, err := db.Db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
