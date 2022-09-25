package database

import (
	"context"

	"github.com/achristie/save2db/pg"
	"github.com/achristie/save2db/sqlite"
)

func NewSqliteDB(source string) *sqlite.DB {
	db := &sqlite.DB{
		Source: source,
	}
	db.Ctx, db.Cancel = context.WithCancel(context.Background())
	return db
}

func NewPgDB(source string) *pg.DB {
	db := &pg.DB{Source: source}
	db.Ctx, db.Cancel = context.WithCancel(context.Background())
	return db
}
