package trades

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/achristie/save2db/pkg/platts"
)

type TradeDataStore struct {
	database *sql.DB
}

const (
	trade_insert = `INSERT INTO trades(
			buyer, buyer_mnemonic, buyer_parent, c1_basis_period, c1_basis_period_details, c1_percentage, c1_price, c1_price_basis,
			c2_basis_period, c2_basis_period_details , c2_percentage , c2_price , c2_price_basis , c3_basis_period ,
			c3_basis_period_details , c3_percentage , c3_price , c3_price_basis , counterparty , counterparty_mnemonic , counterparty_parent ,
			deal_begin , deal_end , deal_id , deal_quantity ,
			deal_quantity_max , deal_quantity_min , deal_terms , hub , leg_prices , lot_size , lot_unit, 
			market_maker , market_maker_mnemonic , market_maker_parent, market_type , oco_order_id , order_begin , order_cancelled ,
			order_classification , order_date , order_derived , order_end , order_id , order_platts_id , order_quantity , order_quantity_total , order_repeat ,
			order_sequence , order_spread , order_state , order_state_detail , order_time , order_type , parent_deal_id , price ,
			price_unit , product , reference_order_id , seller , seller_mnemonic , seller_parent , strip , update_time ,
			window_region , window_state, markets
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?)
`

	trade_create = `CREATE TABLE IF NOT EXISTS trades
	(
		markets json,
		product TEXT,
		strip TEXT,
		hub TEXT,
		update_time TEXT,
		market_maker TEXT,
		order_type TEXT,
		order_state TEXT,
		buyer TEXT,
		seller TEXT,
		price REAL,
		price_unit TEXT,
		order_quantity REAL,
		lot_size INTEGER,
		lot_unit TEXT,
		order_begin TEXT,
		order_end TEXT,
		order_date TEXT,
		order_time TEXT,
		order_id INTEGER,
		order_sequence INTEGER,
		deal_id INTEGER,
		deal_begin TEXT,
		deal_end TEXT,
		deal_quantity REAL,
		deal_quantity_max REAL,
		deal_quantity_min REAL,
		deal_terms TEXT,
		counterparty_parent TEXT,
		counterparty TEXT,
		market_maker_parent text,
		buyer_parent TEXT,
		seller_parent TEXT,
		buyer_mnemonic TEXT,
		seller_mnemonic TEXT,
		market_maker_mnemonic TEXT,
		counterparty_mnemonic TEXT,
		window_region TEXT,
		market_type TEXT,
		c1_price_basis TEXT,
		c1_percentage REAL,
		c1_price REAL,
		c1_basis_period TEXT,
		c1_basis_period_details TEXT,
		c2_price_basis TEXT,
		c2_percentage REAL,
		c2_price REAL,
		c2_basis_period TEXT,
		c2_basis_period_details TEXT,
		c3_price_basis TEXT,
		c3_percentage REAL,
		c3_price REAL,
		c3_basis_period TEXT,
		c3_basis_period_details TEXT,
		window_state TEXT,
		order_classification TEXT,
		oco_order_id TEXT,
		reference_order_id INTEGER,
		order_platts_id INTEGER,
		order_cancelled TEXT,
		order_derived TEXT,
		order_quantity_total REAL,
		order_repeat TEXT,
		leg_prices TEXT,
		parent_deal_id TEXT,
		order_spread TEXT,
		order_state_detail TEXT
	);`
)

func createTradeDataTable(db *sql.DB) {
	query, err := db.Prepare(trade_create)
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	_, err = query.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

// Create DB and table `trades`
func NewTradeDataStore(db *sql.DB) *TradeDataStore {
	createTradeDataTable(db)

	return &TradeDataStore{database: db}
}

// Add Reference Data to DB
func (r *TradeDataStore) Add(data []platts.TradeResults) error {

	query, err := r.database.Prepare(trade_insert)
	if err != nil {
		return err
	}
	defer query.Close()

	tx, err := r.database.Begin()
	if err != nil {
		return err
	}

	for _, r := range data {

		// convert Markets to json
		mkts, err := json.Marshal(&r.Markets)
		if err != nil {
			return err
		}

		_, err = tx.Stmt(query).Exec(r.Buyer, r.BuyerMnemonic, r.BuyerParent, r.C1BasisPeriod,
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
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
