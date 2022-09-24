package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	platts "github.com/achristie/save2db/pkg/platts"
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

// // Remove deleted records from the DB
// func (m *AssessmentsStore) Remove(records []platts.Assessment) error {
// 	query, err := m.database.Prepare(del)
// 	if err != nil {
// 		return err
// 	}
// 	defer query.Close()

// 	// bulk delete
// 	tx, err := m.database.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	for _, r := range records {
// 		_, err := tx.Stmt(query).Exec(r.Symbol, r.Bate, r.AssessDate)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}
// 	return tx.Commit()
// }

func (s *AssessmentsService) Add(ctx context.Context, tx *Tx, record platts.Assessment) (sql.Result, error) {
	res, err := tx.StmtContext(ctx, s.insert).Exec(record.Symbol, record.Bate, record.Value, record.AssessDate, record.ModDate, record.IsCorrected)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Add Assessments
// func (m *AssessmentsStore) Add(ctx context.Context, records []platts.Assessment) error {
// 	query, err := m.database.Prepare(ins)
// 	if err != nil {
// 		return err
// 	}
// 	defer query.Close()

// 	// bulk insert
// 	tx, err := m.database.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	for _, r := range records {
// 		_, err := tx.Stmt(query).Exec(r.Symbol, r.Bate, r.Value, r.AssessDate, r.ModDate, r.IsCorrected)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 	}
// 	return tx.Commit()
// }
