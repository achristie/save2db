package ref_data

import (
	"database/sql"
	"log"
	"os"
)

type RefDataStore struct {
	database *sql.DB
}

func createTable(db *sql.DB) {
	ref_data_table := `CREATE TABLE IF NOT EXISTS ref_data (
		"symbol" TEXT NOT NULL,
		"assessment_frequency" TEXT NOT NULL,
		"commodity" TEXT NOT NULL,
		"contract_type" TEXT NOT NULL,
		"description" TEXT NOT NULL, "publication_frequency_code" TEXT NOT NULL,
		"quotation_style" TEXT NOT NULL,
		"delivery_region" TEXT NOT NULL,
		"delivery_region_basis" TEXT NOT NULL,
		"settlement_type" TEXT NOT NULL,
		"active" TEXT NOT NULL,
		"timestamp" TEXT NOT NULL,
		"uom" TEXT NOT NULL
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
