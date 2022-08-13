package market_data

import (
	"database/sql"
	"log"

	platts "github.com/achristie/save2db/pkg/platts"
)

type AssessmentsStore struct {
	database *sql.DB
}

type dbClass struct {
	symbol       string
	bate         string
	modifiedDate string
	assessedDate string
	price        float64
	isCorrected  string
}

// creating the market_data table
func createAssessmentTables(db *sql.DB) {
	assessments_table := `CREATE TABLE IF NOT EXISTS assessments (
		"symbol" TEXT NOT NULL,
		"bate" TEXT NOT NULL,
		"price" NUM NOT NULL,
		"assessed_date" datetime NOT NULL,
		"modified_date" datetime NOT NULL,
		"is_corrected" string NOT NULL
	);`
	// PRIMARY KEY (symbol, bate, assessed_date) );`
	query, err := db.Prepare(assessments_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	log.Println("db: assessments table created succesfully")
}

// Create our DB (if it does not exist)
// and create `market_data` table
func NewAssessmentsStore(db *sql.DB) *AssessmentsStore {

	createAssessmentTables(db)

	return &AssessmentsStore{database: db}
}

// Remove a Symbol-Bate-Assessed Date from the DB
func (m *AssessmentsStore) Remove(data platts.SymbolCorrection) error {
	var records []dbClass
	for _, v := range data.Results {
		for _, v2 := range v.Data {
			records = append(records, dbClass{
				symbol:       v.Symbol,
				bate:         v2.Bate,
				assessedDate: v2.AssessDate,
			})
		}
	}
	return m.del(records)
}

// del is the internal implementation for Remove
func (m *AssessmentsStore) del(records []dbClass) error {
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
		_, err := tx.Stmt(query).Exec(r.symbol, r.bate, r.assessedDate)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// Add symbol history to the database
func (m *AssessmentsStore) Add(data platts.SymbolHistory) error {
	var records []dbClass
	// change data structure for ease of INSERT
	for _, v := range data.Results {
		for _, v2 := range v.Data {
			records = append(records, dbClass{
				symbol:       v.Symbol,
				bate:         v2.Bate,
				price:        v2.Value,
				modifiedDate: v2.ModDate,
				assessedDate: v2.AssessDate,
				isCorrected:  v2.IsCorrected,
			})
		}
	}
	return m.insert(records)
}

// insert is the internal implementation for Add
func (m *AssessmentsStore) insert(records []dbClass) error {
	ins := `INSERT or REPLACE INTO assessments(symbol, bate, price, assessed_date, modified_date, is_corrected) VALUES(?, ?, ?, ?, ?, ?)`
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
		_, err := tx.Stmt(query).Exec(r.symbol, r.bate, r.price, r.assessedDate, r.modifiedDate, r.isCorrected)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
