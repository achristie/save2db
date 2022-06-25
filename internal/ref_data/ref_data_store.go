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

const (
	ref_data_table = `CREATE TABLE IF NOT EXISTS ref_data (
		"symbol" TEXT NOT NULL PRIMARY KEY,
		"assessment_frequency" TEXT,
		"commodity" TEXT,
		"contract_type" TEXT,
		"description" TEXT, 
		"publication_frequency_code" TEXT,
		"currency" TEXT,
		"quotation_style" TEXT,
		"delivery_region" TEXT,
		"delivery_region_basis" TEXT,
		"settlement_type" TEXT,
		"active" TEXT,
		"timestamp" TEXT,
		"uom" TEXT,
		"day_of_publication" TEXT,
		"shipping_terms" TEXT,
		"standard_lot_size" INTEGER,
		"commodity_grade" INTEGER,
		"standard_lot_units" TEXT
	);"`
	sym_bate_table = `CREATE TABLE IF NOT EXISTS sym_bate (
			"symbol" TEXT NOT NULL,
			"bate" TEXT NOT NULL	
		);"`
	sym_mdc_table = `CREATE TABLE IF NOT EXISTS sym_mdc (
		"symbol" TEXT NOT NULL,
		"mdc" TEXT NOT NULL,
		"mdc_description" TEXT NOT NULL
	);"`
)

func createTable(db *sql.DB) {
	creates := map[string]string{"ref_data": ref_data_table, "sym_bate": sym_bate_table, "sym_mdc": sym_mdc_table}

	for k, v := range creates {
		query, err := db.Prepare(v)
		if err != nil {
			log.Fatal(err)
		}
		query.Exec()
		log.Printf("db: %s created succesfully", k)

	}

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

	bate_stmt := `INSERT OR REPLACE INTO sym_bate(symbol, bate) VALUES(?, ?)`
	b_query, err := r.database.Prepare(bate_stmt)
	if err != nil {
		return err
	}
	defer b_query.Close()

	del_stmt := `DELETE FROM sym_bate WHERE symbol=?`
	d_query, err := r.database.Prepare(del_stmt)
	if err != nil {
		return err
	}
	defer d_query.Close()

	mdc_stmt := `INSERT OR REPLACE INTO sym_mdc(symbol, mdc, mdc_description) VALUES(?, ?, ?)`
	m_query, err := r.database.Prepare(mdc_stmt)
	if err != nil {
		return err
	}
	defer b_query.Close()

	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range data.Results {
		_, err := tx.Stmt(query).Exec(r.Symbol, r.Description, r.Commodity, r.UOM, r.Active, r.DeliveryRegion)
		_, err2 := tx.Stmt(d_query).Exec(r.Symbol)
		for _, b := range r.Bate {
			_, err := tx.Stmt(b_query).Exec(r.Symbol, b)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		for i, m := range r.MDC {
			_, err := tx.Stmt(m_query).Exec(r.Symbol, m, r.MDCDescription[i])
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		if err != nil || err2 != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
