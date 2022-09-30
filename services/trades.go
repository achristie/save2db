package services

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/achristie/save2db/pkg/platts"
)

//go:embed scripts/pg/trades/insert.sql
var insert_trade_pg string

//go:embed scripts/sqlite/trades/insert.sql
var insert_trade_sqlite string

type TradeService struct {
	insert *sql.Stmt
}

func getTradeStmts(s string) string {
	switch s {
	case "PostgreSQL":
		return insert_trade_pg
	case "SQLite":
		return insert_trade_sqlite
	default:
		return insert_trade_sqlite
	}

}

func NewTradeService(ctx context.Context, db *sql.DB, dbSelection string) (*TradeService, error) {
	ins := getTradeStmts(dbSelection)
	insert, err := db.PrepareContext(ctx, ins)
	if err != nil {
		return nil, fmt.Errorf("insert statement: %w", err)
	}

	ts := TradeService{
		insert: insert,
	}
	return &ts, nil
}
func (s *TradeService) Add(ctx context.Context, tx *sql.Tx, r platts.TradeResults) (sql.Result, error) {
	// convert Markets to json
	mkts, err := json.Marshal(&r.Markets)
	if err != nil {
		return nil, err
	}

	res, err := tx.Stmt(s.insert).Exec(r.Buyer, r.BuyerMnemonic, r.BuyerParent, r.C1BasisPeriod,
		r.C1BasisPeriodDetails, r.C1Percentage, r.C1Price, r.C1PriceBasis, r.C2BasisPeriod,
		r.C2BasisPeriodDetails, r.C2Percentage, r.C2Price, r.C2PriceBasis, r.C3BasisPeriod,
		r.C3BasisPeriodDetails, r.C3Percentage, r.C3Price, r.C3PriceBasis, r.Counterparty, r.CounterpartyMnemonic,
		r.CounterpartyParent, r.DealBegin, r.DealEnd, r.DealID, r.DealQuantity, r.DealQuantityMax, r.DealQuantityMin, r.DealTerms,
		r.Hub, r.LegPrices, r.LotSize, r.LotUnit, r.MarketMaker, r.MarketMakerMnemonic, r.MarketMakerParent, r.MarketType,
		r.OcoOrderID, r.OrderBegin, r.OrderCancelled, r.OrderClassification, r.OrderDate, r.OrderDerived, r.OrderEnd, r.OrderID,
		r.OrderPlattsID, r.OrderQuantity, r.OrderQuantityTotal, r.OrderRepeat, r.OrderSequence, r.OrderSpread, r.OrderState, r.OrderStateDetail,
		r.OrderTime, r.OrderType, r.ParentDealID, r.Price, r.PriceUnit, r.Product, r.ReferenceOrderID, r.Seller, r.SellerMnemonic, r.SellerParent,
		r.Strip, r.UpdateTime, r.WindowRegion, r.WindowState, mkts)

	if err != nil {
		return nil, err
	}
	return res, nil

}
