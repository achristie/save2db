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
		"standard_lot_units" TEXT,
		"decimal_places" INTEGER
	);"`
	ref_data_insert = `INSERT or REPLACE INTO ref_data(
		symbol,
		assessment_frequency,
		commodity,
		contract_type,
		description,
		publication_frequency_code,
		currency,
		quotation_style,
		delivery_region,
		delivery_region_basis, 
		settlement_type,
		active, 
		timestamp, 
		uom, 
		day_of_publication, 
		shipping_terms, 
		standard_lot_size, 
		commodity_grade, 
		standard_lot_units,
		decimal_places
	)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	sym_bate_table = `CREATE TABLE IF NOT EXISTS sym_bate (
			"symbol" TEXT NOT NULL,
			"bate" TEXT NOT NULL	
		);"`
	bate_insert   = `INSERT OR REPLACE INTO sym_bate(symbol, bate) VALUES(?, ?)`
	bate_delete   = `DELETE FROM sym_bate WHERE symbol=?`
	sym_mdc_table = `CREATE TABLE IF NOT EXISTS sym_mdc (
		"symbol" TEXT NOT NULL,
		"mdc" TEXT NOT NULL,
		"mdc_description" TEXT NOT NULL
	);"`
	mdc_insert = `INSERT OR REPLACE INTO sym_mdc(symbol, mdc, mdc_description) VALUES(?, ?, ?)`
	mdc_delete = `DELETE FROM sym_mdc WHERE symbol=?`
)

func createTable(db *sql.DB) {
	creates := map[string]string{"ref_data": ref_data_table, "sym_bate": sym_bate_table, "sym_mdc": sym_mdc_table}

	for k, v := range creates {
		func() {
			query, err := db.Prepare(v)
			if err != nil {
				log.Fatal(err)
			}
			defer query.Close()

			query.Exec()
			log.Printf("db: %s created succesfully", k)
		}()

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
	stmts := map[string]string{"ref_data_insert": ref_data_insert,
		"bate_insert": bate_insert, "bate_delete": bate_delete,
		"mdc_insert": mdc_insert, "mdc_delete": mdc_delete}
	queries := map[string]*sql.Stmt{}

	for k, v := range stmts {
		query, err := r.database.Prepare(v)
		if err != nil {
			return err
		}
		queries[k] = query
	}
	// defer queries["ref_data_insert"].Close()
	// defer queries["bate_insert"].Close()
	// defer queries["bate_delet"].Close()
	// defer queries["mdc_insert"].Close()

	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range data.Results {
		_, err := tx.Stmt(queries["ref_data_insert"]).Exec(r.Symbol, r.AssessmentFrequency, r.Commodity,
			r.ContractType, r.Description, r.PublicationFrequencyCode, r.Currency,
			r.QuotationStyle, r.DeliveryRegion, r.DeliveryRegionBasis, r.SettlementType,
			r.Active, r.Timestamp, r.UOM, r.DayOfPublication, r.ShippingTerms,
			r.StandardLotSize, r.CommodityGrade, r.StandardLotUnits, r.DecimalPlaces)
		_, err2 := tx.Stmt(queries["bate_delete"]).Exec(r.Symbol)
		_, err3 := tx.Stmt(queries["mdc_delete"]).Exec(r.Symbol)
		for _, b := range r.Bate {
			_, err := tx.Stmt(queries["bate_insert"]).Exec(r.Symbol, b)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		for i, m := range r.MDC {
			_, err := tx.Stmt(queries["mdc_insert"]).Exec(r.Symbol, m, r.MDCDescription[i])
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		if err != nil || err2 != nil || err3 != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
