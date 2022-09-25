package pg

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	"github.com/achristie/save2db/pkg/platts"
)

//go:embed scripts/assessments/insert.sql
var insert string

//go:embed scripts/assessments/delete.sql
var delete string

type AssessmentsService struct {
	db     *DB
	insert *sql.Stmt
	delete *sql.Stmt
}

// NewAssessmentsService returns a new instance of AssessmentsService
func NewAssessmentsService(ctx context.Context, db *DB) (*AssessmentsService, error) {
	insert, err := db.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, fmt.Errorf("insert statement: %w", err)
	}

	delete, err := db.db.PrepareContext(ctx, delete)
	if err != nil {
		return nil, fmt.Errorf("delete statement: %w", err)
	}
	as := AssessmentsService{
		db:     db,
		insert: insert,
		delete: delete,
	}
	return &as, nil
}

func (s *AssessmentsService) Remove(ctx context.Context, tx *Tx, record platts.Assessment) (sql.Result, error) {
	res, err := tx.StmtContext(ctx, s.delete).Exec(record.Symbol, record.Bate, record.Value, record.AssessDate, record.ModDate, record.IsCorrected)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AssessmentsService) Add(ctx context.Context, tx *Tx, record platts.Assessment) (sql.Result, error) {
	res, err := tx.StmtContext(ctx, s.insert).Exec(record.Symbol, record.Bate, record.Value, record.AssessDate, record.ModDate, record.IsCorrected)
	if err != nil {
		return nil, err
	}
	return res, nil
}
