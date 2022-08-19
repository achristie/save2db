package platts

import "encoding/json"

type TradeData struct {
	Metadata struct {
		Count      int    `json:"count"`
		Pagesize   int    `json:"pagesize"`
		Page       int    `json:"page"`
		Maxid      string `json:"maxid"`
		TotalPages int    `json:"total_pages"`
		QueryTime  string `json:"query_time"`
	} `json:"metadata"`
	Results []TradeResults `json:"results"`
}
type TradeResults struct {
	Market               []string `json:"market"`
	Product              string   `json:"product"`
	Hub                  string   `json:"hub"`
	Strip                string   `json:"strip"`
	UpdateTime           string   `json:"update_time"`
	MarketMaker          string   `json:"market_maker"`
	OrderType            string   `json:"order_type"`
	OrderState           string   `json:"order_state"`
	Buyer                string   `json:"buyer,omitempty"`
	Seller               string   `json:"seller,omitempty"`
	Price                float32  `json:"price"`
	PriceUnit            string   `json:"price_unit"`
	OrderQuantity        float32  `json:"order_quantity"`
	LotSize              int      `json:"lot_size"`
	LotUnit              string   `json:"lot_unit"`
	OrderBegin           string   `json:"order_begin"`
	OrderEnd             string   `json:"order_end"`
	OrderDate            string   `json:"order_date"`
	OrderTime            string   `json:"order_time"`
	OrderID              int      `json:"order_id"`
	OrderSequence        int      `json:"order_sequence"`
	DealID               int      `json:"deal_id,omitempty"`
	DealBegin            string   `json:"deal_begin,omitempty"`
	DealEnd              string   `json:"deal_end,omitempty"`
	DealQuantity         float32  `json:"deal_quantity,omitempty"`
	DealQuantityMin      float32  `json:"deal_quantity_min,omitempty"`
	DealQuantityMax      float32  `json:"deal_quantity_max,omitempty"`
	DealTerms            string   `json:"deal_terms,omitempty"`
	CounterpartyParent   string   `json:"counterparty_parent,omitempty"`
	Counterparty         string   `json:"counterparty,omitempty"`
	MarketMakerParent    string   `json:"market_maker_parent"`
	BuyerParent          string   `json:"buyer_parent,omitempty"`
	SellerParent         string   `json:"seller_parent,omitempty"`
	BuyerMnemonic        string   `json:"buyer_mnemonic,omitempty"`
	SellerMnemonic       string   `json:"seller_mnemonic,omitempty"`
	MarketMakerMnemonic  string   `json:"market_maker_mnemonic"`
	CounterpartyMnemonic string   `json:"counterparty_mnemonic,omitempty"`
	WindowRegion         string   `json:"window_region"`
	MarketShortCode      []string `json:"market_short_code"`
	MarketType           string   `json:"market_type"`
	C1PriceBasis         string   `json:"c1_price_basis,omitempty"`
	C1Percentage         int      `json:"c1_percentage"`
	C1Price              float32  `json:"c1_price"`
	C1BasisPeriod        string   `json:"c1_basis_period,omitempty"`
	C1BasisPeriodDetails string   `json:"c1_basis_period_details,omitempty"`
	C2PriceBasis         string   `json:"c2_price_basis,omitempty"`
	C2Percentage         int      `json:"c2_percentage,omitempty"`
	C2Price              int      `json:"c2_price"`
	C2BasisPeriod        string   `json:"c2_basis_period,omitempty"`
	C2BasisPeriodDetails string   `json:"c2_basis_period_details,omitempty"`
	C3PriceBasis         string   `json:"c3_price_basis,omitempty"`
	C3Percentage         int      `json:"c3_percentage"`
	C3Price              int      `json:"c3_price"`
	C3BasisPeriod        string   `json:"c3_basis_period,omitempty"`
	C3BasisPeriodDetails string   `json:"c3_basis_period_details,omitempty"`
	WindowState          string   `json:"window_state"`
	OrderClassification  string   `json:"order_classification,omitempty"`
	OcoOrderID           string   `json:"oco_order_id,omitempty"`
	ReferenceOrderID     int      `json:"reference_order_id,omitempty"`
	OrderPlattsID        int      `json:"order_platts_id"`
	OrderCancelled       string   `json:"order_cancelled"`
	OrderDerived         string   `json:"order_derived"`
	OrderQuantityTotal   float32  `json:"order_quantity_total"`
	OrderRepeat          string   `json:"order_repeat"`
	LegPrices            string   `json:"leg_prices,omitempty"`
	ParentDealID         string   `json:"parent_deal_id,omitempty"`
	OrderSpread          string   `json:"order_spread"`
	OrderStateDetail     string   `json:"order_state_detail"`
	Markets              []Market `json:"-"`
}

// Zip the Market and Market Short Code fields
func (t *TradeResults) UnmarshalJSON(data []byte) error {
	type T TradeResults
	if err := json.Unmarshal(data, (*T)(t)); err != nil {
		return err
	}
	var m []Market
	for i, v := range t.Market {
		m = append(m, Market{Name: v, ShortCode: t.MarketShortCode[i]})
	}

	t.Markets = m

	return nil
}

func (t TradeData) GetTotalPages() int {
	return t.Metadata.TotalPages
}

type Market struct {
	Name      string `json:"name"`
	ShortCode string `json:"short_code"`
}
