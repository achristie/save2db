package platts

import (
	"encoding/json"
)

type SymbolData struct {
	Metadata SymbolMetadata  `json:"metadata"`
	Results  []SymbolResults `json:"results"`
}

type SymbolResults struct {
	Symbol                   string   `json:"symbol"`
	Description              string   `json:"description"`
	Commodity                string   `json:"commodity"`
	UOM                      string   `json:"uom"`
	Active                   string   `json:"active"`
	DeliveryRegion           string   `json:"delivery_region"`
	DeliveryRegionBasis      string   `json:"delivery_region_basis"`
	ContractType             string   `json:"contract_type"`
	PublicationFrequencyCode string   `json:"publication_frequency_code"`
	ShippingTerms            string   `json:"shipping_terms"`
	DayOfPublication         string   `json:"day_of_publication"`
	StandardLotSize          float32  `json:"standard_lot_size"`
	StandardLotUnits         string   `json:"standard_lot_units"`
	QuotationStyle           string   `json:"quotation_style"`
	Bate                     []string `json:"bate_code"`
	CommodityGrade           string   `json:"commodity_grade"`
	Currency                 string   `json:"currency"`
	AssessmentFrequency      string   `json:"assessment_frequency"`
	Timestamp                string   `json:"timestamp"`
	SettlementType           string   `json:"settlement_type"`
	DecimalPlaces            int      `json:"decimal_places"`
	MDCNames                 []string `json:"mdc"`
	MDCDescriptions          []string `json:"mdc_description"`
	MDC                      []MDC    `json:"-"`
}

// Extend unmarshalling to zip the MDC fields
func (r *SymbolResults) UnmarshalJSON(data []byte) error {
	type R SymbolResults
	if err := json.Unmarshal(data, (*R)(r)); err != nil {
		return err
	}
	var m []MDC
	for i, v := range r.MDCNames {
		m = append(m, MDC{Name: v, Description: r.MDCDescriptions[i]})
	}

	r.MDC = m

	return nil
}

type Metadata struct {
	Count      int    `json:"count"`
	PageSize   int    `json:"pageSize"`
	Page       int    `json:"page"`
	TotalPages int    `json:"totalPages"`
	QueryTime  string `json:"queryTime"`
}

type SymbolMetadata struct {
	Count      int    `json:"count"`
	PageSize   int    `json:"page_size"`
	Page       int    `json:"page"`
	TotalPages int    `json:"total_pages"`
	QueryTime  string `json:"query_time"`
}

type SymbolCorrection struct {
	Metadata Metadata `json:"metadata"`
	Results  []struct {
		Symbol string `json:"symbol"`
		Data   []struct {
			Bate           string  `json:"bate"`
			Value          float32 `json:"value"`
			AssessDate     string  `json:"assessDate"`
			CorrectionType string  `json:"correctionType"`
		} `json:"data"`
	} `json:"results"`
}

type SymbolHistory struct {
	Metadata Metadata `json:"metadata"`
	Results  []struct {
		Symbol string `json:"symbol"`
		Data   []struct {
			Bate        string  `json:"bate"`
			Value       float64 `json:"value"`
			AssessDate  string  `json:"assessDate"`
			IsCorrected string  `json:"isCorrected"`
			ModDate     string  `json:"modDate"`
		} `json:"data"`
	} `json:"results"`
}

type Assessment struct {
	Symbol      string
	Bate        string
	Value       float64
	AssessDate  string
	IsCorrected string
	ModDate     string
}

func (s SymbolData) GetTotalPages() int {
	return s.Metadata.TotalPages
}
func (s SymbolData) GetResults() []interface{} {
	i := make([]interface{}, len(s.Results))
	for idx, r := range s.Results {
		i[idx] = r
	}
	return i
}

func (s SymbolCorrection) GetTotalPages() int {
	return s.Metadata.TotalPages
}
func (s SymbolCorrection) GetResults() []interface{} {
	a := s.Flatten()
	i := make([]interface{}, len(a))
	for idx, r := range a {
		i[idx] = r
	}
	return i
}

func (s SymbolHistory) GetTotalPages() int {
	return s.Metadata.TotalPages
}
func (s SymbolHistory) GetResults() []interface{} {
	a := s.Flatten()
	i := make([]interface{}, len(a))
	for idx, r := range a {
		i[idx] = r
	}
	return i
}

type MDC struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
