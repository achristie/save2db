package symbols

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/achristie/save2db/pkg/platts"
)

//go:embed pg/insert.sql
var insert_sym_pg string

//go:embed sqlite/insert.sql
var insert_sym_sqlite string

type SymbolService struct {
	insert *sql.Stmt
}

func getSymbolStmts(s string) string {
	switch s {
	case "PostgreSQL":
		return insert_sym_pg
	case "SQLite":
		return insert_sym_sqlite
	default:
		return insert_sym_sqlite
	}

}

func New(db *sql.DB, dbSelection string) (*SymbolService, error) {
	ins := getSymbolStmts(dbSelection)
	insert, err := db.PrepareContext(context.TODO(), ins)
	if err != nil {
		return nil, fmt.Errorf("insert statement: %w", err)
	}

	ss := SymbolService{
		insert: insert,
	}
	return &ss, nil
}

func (s *SymbolService) Add(tx *sql.Tx, r interface{}) (sql.Result, error) {
	// convert bates to JSON
	record := r.(platts.SymbolResults)
	bates, err := json.Marshal(&record.Bate)
	if err != nil {
		return nil, err
	}
	// convert MDC to json
	mdcs := new(bytes.Buffer)
	enc := json.NewEncoder(mdcs)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(&record.MDC); err != nil {
		return nil, err
	}

	res, err := tx.Stmt(s.insert).Exec(record.Symbol, record.AssessmentFrequency, record.Commodity,
		record.ContractType, record.Description, record.PublicationFrequencyCode, record.Currency,
		record.QuotationStyle, record.DeliveryRegion, record.DeliveryRegionBasis, record.SettlementType,
		record.Active, record.Timestamp, record.UOM, record.DayOfPublication, record.ShippingTerms,
		record.StandardLotSize, record.CommodityGrade, record.StandardLotUnits, record.DecimalPlaces, mdcs.String(), string(bates))

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *SymbolService) Remove(tx *sql.Tx, r interface{}) (sql.Result, error) {
	panic("not implemented")
}
