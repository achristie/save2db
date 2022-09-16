package market_data

import (
	"database/sql"
	"log"

	platts "github.com/achristie/save2db/pkg/platts"
)

type AssessmentsStore struct {
	database *sql.DB
}

// creating the assessments table
func createAssessmentTables(db *sql.DB) {
	assessments_table := `CREATE TABLE IF NOT EXISTS assessments (
		"symbol" TEXT NOT NULL,
		"bate" TEXT NOT NULL,
		"value" NUM NOT NULL,
		"assessed_date" TEXT NOT NULL,
		"modified_date" TEXT NOT NULL,
		"is_corrected" string NOT NULL,
	PRIMARY KEY (symbol, bate, assessed_date) );`

	query, err := db.Prepare(assessments_table)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	query.Exec()
}

// Create DB and `assessments` table
func NewAssessmentsStore(db *sql.DB) *AssessmentsStore {
	createAssessmentTables(db)

	return &AssessmentsStore{database: db}
}

// Remove deleted records from the DB
func (m *AssessmentsStore) Remove(records []platts.Assessment) error {
	del := `DELETE FROM assessments where symbol = ? and bate = ? and assessed_date = ?`
	query, err := m.database.Prepare(del)
	if err != nil {
		return err
	}
	defer query.Close()

	// bulk delete
	tx, err := m.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range records {
		_, err := tx.Stmt(query).Exec(r.Symbol, r.Bate, r.AssessDate)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// Add Assessments
func (m *AssessmentsStore) Add(records []platts.Assessment) error {
	ins := `INSERT or REPLACE INTO assessments(symbol, bate, value, assessed_date, modified_date, is_corrected) VALUES(?, ?, ?, ?, ?, ?)`
	query, err := m.database.Prepare(ins)
	if err != nil {
		return err
	}
	defer query.Close()

	// bulk insert
	tx, err := m.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range records {
		_, err := tx.Stmt(query).Exec(r.Symbol, r.Bate, r.Value, r.AssessDate, r.ModDate, r.IsCorrected)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
