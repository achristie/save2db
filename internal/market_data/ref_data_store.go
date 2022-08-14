package market_data

import (
	"database/sql"
	"log"

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
		"decimal_places" INTEGER,
		"mdc" json
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
)

func createRefDataTable(db *sql.DB) {

	query, err := db.Prepare(ref_data_table)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	query.Exec()
	log.Println("db: ref_data  table created succesfully")

}

// Create our DB (if it does not exist)
// and create `ref_data` table
func NewRefDataStore(db *sql.DB) *RefDataStore {
	createRefDataTable(db)

	return &RefDataStore{database: db}
}

func (r *RefDataStore) Add(data platts.ReferenceData) error {

	query, err := r.database.Prepare(ref_data_insert)
	if err != nil {
		return err
	}
	defer query.Close()

	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range data.Results {
		_, err := tx.Stmt(query).Exec(r.Symbol, r.AssessmentFrequency, r.Commodity,
			r.ContractType, r.Description, r.PublicationFrequencyCode, r.Currency,
			r.QuotationStyle, r.DeliveryRegion, r.DeliveryRegionBasis, r.SettlementType,
			r.Active, r.Timestamp, r.UOM, r.DayOfPublication, r.ShippingTerms,
			r.StandardLotSize, r.CommodityGrade, r.StandardLotUnits, r.DecimalPlaces)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
