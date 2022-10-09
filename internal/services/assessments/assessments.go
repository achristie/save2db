package assessments

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "embed"

	"github.com/achristie/save2db/pkg/platts"
)

//go:embed pg/insert.sql
var insert_pg string

//go:embed pg/delete.sql
var delete_pg string

//go:embed sqlite/insert.sql
var insert_sqlite string

//go:embed sqlite/delete.sql
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

func New(ctx context.Context, db *sql.DB, dbSelection string) (*AssessmentsService, error) {
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

func (s *AssessmentsService) Remove(ctx context.Context, tx *sql.Tx, r interface{}) (sql.Result, error) {
	record, ok := r.(platts.Assessment)
	if !ok {
		return nil, fmt.Errorf("remove: must use a platts.assessment")
	}
	log.Printf("%+v", record)

	res, err := tx.StmtContext(ctx, s.delete).Exec(record.Symbol, record.Bate, record.AssessDate)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AssessmentsService) Add(ctx context.Context, tx *sql.Tx, r interface{}) (sql.Result, error) {
	record, ok := r.(platts.Assessment)
	if !ok {
		return nil, fmt.Errorf("remove: must use a platts.assessment")
	}

	res, err := tx.StmtContext(ctx, s.insert).Exec(record.Symbol, record.Bate, record.Value, record.AssessDate, record.ModDate, record.IsCorrected)
	if err != nil {
		return nil, err
	}
	return res, nil
}
