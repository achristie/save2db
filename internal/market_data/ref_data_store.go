package market_data

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/achristie/save2db/pkg/platts"
)

type SymbolStore struct {
	database *sql.DB
}

const (
	symbol_create = `CREATE TABLE IF NOT EXISTS symbols (
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
		"mdc" json,
		"bates" json
	);"`
	symbol_insert = `INSERT or REPLACE INTO symbols(
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
		decimal_places,
		mdc,
		bates
	)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

func createSymbolTable(db *sql.DB) {
	query, err := db.Prepare(symbol_create)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	query.Exec()
}

// Create DB and table `ref_data`
func NewSymbolStore(db *sql.DB) *SymbolStore {
	createSymbolTable(db)

	return &SymbolStore{database: db}
}

// Add Reference Data to DB
func (r *SymbolStore) Add(data []platts.SymbolResults) error {
	query, err := r.database.Prepare(symbol_insert)
	if err != nil {
		return err
	}
	defer query.Close()

	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range data {

		// convert bates to JSON
		bates, err := json.Marshal(&r.Bate)
		if err != nil {
			return err
		}

		// convert MDC to json
		mdcs := new(bytes.Buffer)
		enc := json.NewEncoder(mdcs)
		enc.SetEscapeHTML(false)

		if err := enc.Encode(&r.MDC); err != nil {
			return err
		}

		_, err = tx.Stmt(query).Exec(r.Symbol, r.AssessmentFrequency, r.Commodity,
			r.ContractType, r.Description, r.PublicationFrequencyCode, r.Currency,
			r.QuotationStyle, r.DeliveryRegion, r.DeliveryRegionBasis, r.SettlementType,
			r.Active, r.Timestamp, r.UOM, r.DayOfPublication, r.ShippingTerms,
			r.StandardLotSize, r.CommodityGrade, r.StandardLotUnits, r.DecimalPlaces, mdcs.String(), string(bates))

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
