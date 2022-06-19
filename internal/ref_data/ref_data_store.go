package ref_data

import (
	"database/sql"
	"log"
	"os"

	"github.com/achristie/save2db/pkg/platts"
)

type RefDataStore struct {
	database *sql.DB
}

func createTable(db *sql.DB) {
	ref_data_table := `CREATE TABLE IF NOT EXISTS ref_data (
		"symbol" TEXT NOT NULL,
		"assessment_frequency" TEXT,
		"commodity" TEXT,
		"contract_type" TEXT,
		"description" TEXT, 
		"publication_frequency_code" TEXT ,
		"quotation_style" TEXT,
		"delivery_region" TEXT,
		"delivery_region_basis" TEXT,
		"settlement_type" TEXT,
		"active" TEXT,
		"timestamp" TEXT,
		"uom" TEXT 
	);"`

	query, err := db.Prepare(ref_data_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
	log.Println("db: ref_data_table created succesfully")
}

// Create our DB (if it does not exist)
// and create `ref_data` table
func InitializeDb(dbFileName string) *RefDataStore {
	file, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("could not open %s %v", dbFileName, err)
	}

	db, err := sql.Open("sqlite3", file.Name())
	if err != nil {
		log.Fatal(err)
	}

	createTable(db)

	return &RefDataStore{database: db}
}

func (r *RefDataStore) Add(data platts.ReferenceData) error {
	stmt := `INSERT or REPLACE INTO ref_data(symbol, description, commodity, uom, active, delivery_region)
	 VALUES(?, ?, ?, ?, ?, ?)`
	query, err := r.database.Prepare(stmt)
	if err != nil {
		return err
	}
	defer query.Close()

	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range data.Results {
		_, err := tx.Stmt(query).Exec(r.Symbol, r.Description, r.Commodity, r.UOM, r.Active, r.DeliveryRegion)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
