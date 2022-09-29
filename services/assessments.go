package services

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	"github.com/achristie/save2db/pkg/platts"
)

//go:embed scripts/pg/assessments/insert.sql
var insert_pg string

//go:embed scripts/pg/assessments/delete.sql
var delete_pg string

//go:embed scripts/sqlite/assessments/insert.sql
var insert_sqlite string

//go:embed scripts/sqlite/assessments/delete.sql
var delete_sqlite string

type AssessmentsService struct {
	insert *sql.Stmt
	delete *sql.Stmt
}

func getPreparedStmts(s string) (string, string) {
	switch s {
	case "PostgreSQL":
		return insert_pg, delete_pg
	case "SQLite":
		return insert_sqlite, delete_sqlite
	default:
		return insert_sqlite, delete_sqlite
	}

}

func NewAssessmentsService(ctx context.Context, db *sql.DB, dbSelection string) (*AssessmentsService, error) {
	ins, del := getPreparedStmts(dbSelection)
	insert, err := db.PrepareContext(ctx, ins)
	if err != nil {
		return nil, fmt.Errorf("insert statement: %w", err)
	}

	delete, err := db.PrepareContext(ctx, del)
	if err != nil {
		return nil, fmt.Errorf("delete statement: %w", err)
	}
	as := AssessmentsService{
		insert: insert,
		delete: delete,
	}
	return &as, nil
}

func (s *AssessmentsService) Remove(ctx context.Context, tx *sql.Tx, record platts.Assessment) (sql.Result, error) {
	res, err := tx.StmtContext(ctx, s.delete).Exec(record.Symbol, record.Bate, record.Value, record.AssessDate, record.ModDate, record.IsCorrected)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AssessmentsService) Add(ctx context.Context, tx *sql.Tx, record platts.Assessment) (sql.Result, error) {
	res, err := tx.StmtContext(ctx, s.insert).Exec(record.Symbol, record.Bate, record.Value, record.AssessDate, record.ModDate, record.IsCorrected)
	if err != nil {
		return nil, err
	}
	return res, nil
}
