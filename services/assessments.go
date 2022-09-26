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

//go:embed scripts/sqlite/assessments/insert.sql
var insert_sqlite string

type AssessmentsService struct {
	insert *sql.Stmt
	delete *sql.Stmt
}

func getPreparedStmts(s string) string {
	switch s {
	case "pg":
		return insert_pg

	case "sqlite":
		return insert_sqlite
	default:
		return insert_pg
	}

}

func NewAssessmentsService(ctx context.Context, db *sql.DB) (*AssessmentsService, error) {
	s := getPreparedStmts("pg")
	insert, err := db.PrepareContext(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("insert statement: %w", err)
	}

	delete, err := db.PrepareContext(ctx, s)
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
