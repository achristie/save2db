package save2db

import (
	"database/sql"
	"log"
	"os"

	platts "github.com/achristie/save2db/pkg/platts"
)

type MarketDataStore struct {
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
func createTable(db *sql.DB) {
	market_data_table := `CREATE TABLE IF NOT EXISTS market_data (
		"symbol" TEXT NOT NULL,
		"bate" TEXT NOT NULL,
		"price" NUM NOT NULL,
		"assessed_date" datetime NOT NULL,
		"modified_date" datetime NOT NULL,
		"is_corrected" string NOT NULL
	);`
	// PRIMARY KEY (symbol, bate, assessed_date) );`
	query, err := db.Prepare(market_data_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	log.Println("db: market_data table created succesfully")
}

// Create our DB (if it does not exist)
// and create `market_data` table
func InitializeDb(dbFileName string) *MarketDataStore {
	file, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("could not open %s %v", dbFileName, err)
	}

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		log.Fatal(err)
	}

	createTable(db)

	return &MarketDataStore{database: db}
}

// Remove a Symbol-Bate-Assessed Date from the DB
func (m *MarketDataStore) Remove(data platts.SymbolCorrection) error {
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
func (m *MarketDataStore) del(records []dbClass) error {
	del := `DELETE FROM market_data where symbol = ? and bate = ? and assessed_date = ?`
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
func (m *MarketDataStore) Add(data platts.SymbolHistory) error {
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
func (m *MarketDataStore) insert(records []dbClass) error {
	ins := `INSERT or REPLACE INTO market_data(symbol, bate, price, assessed_date, modified_date, is_corrected) VALUES(?, ?, ?, ?, ?, ?)`
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
