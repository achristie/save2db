package services

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/achristie/save2db/pkg/platts"
)

//go:embed scripts/pg/symbols/insert.sql
var insert_sym_pg string

//go:embed scripts/sqlite/symbols/insert.sql
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

func NewSymbolService(ctx context.Context, db *sql.DB, dbSelection string) (*SymbolService, error) {
	ins := getSymbolStmts(dbSelection)
	insert, err := db.PrepareContext(ctx, ins)
	if err != nil {
		return nil, fmt.Errorf("insert statement: %w", err)
	}

	ss := SymbolService{
		insert: insert,
	}
	return &ss, nil
}
func (s *SymbolService) Add(ctx context.Context, tx *sql.Tx, record platts.SymbolResults) (sql.Result, error) {
	// convert bates to JSON
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
