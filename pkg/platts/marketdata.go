package platts

import (
	"bytes"
	"encoding/json"
)

type ReferenceData struct {
	Metadata RefMetadata  `json:"metadata"`
	Results  []RefResults `json:"results"`
}

type RefResults struct {
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
	BateJson                 string   `json:"-"`
	CommodityGrade           string   `json:"commodity_grade"`
	Currency                 string   `json:"currency"`
	AssessmentFrequency      string   `json:"assessment_frequency"`
	Timestamp                string   `json:"timestamp"`
	SettlementType           string   `json:"settlement_type"`
	DecimalPlaces            int      `json:"decimal_places"`
	MDCNames                 []string `json:"mdc"`
	MDCDescriptions          []string `json:"mdc_description"`
	MDCJson                  string   `json:"-"`
}

// Extend unmarshalling to zip the MDC fields
// And create *Json fields for ease of saving in DB
func (r *RefResults) UnmarshalJSON(data []byte) error {
	type R RefResults
	if err := json.Unmarshal(data, (*R)(r)); err != nil {
		return err
	}
	var m []MDC
	for i, v := range r.MDCNames {
		m = append(m, MDC{Name: v, Description: r.MDCDescriptions[i]})
	}

	j := new(bytes.Buffer)
	e := json.NewEncoder(j)
	e.SetEscapeHTML(false)

	if err := e.Encode(&m); err != nil {
		return err
	}
	r.MDCJson = j.String()

	b, err := json.Marshal(&r.Bate)
	if err != nil {
		return err
	}
	r.BateJson = string(b)
	return nil
}

type Metadata struct {
	Count      int    `json:"count"`
	PageSize   int    `json:"pageSize"`
	Page       int    `json:"page"`
	TotalPages int    `json:"totalPages"`
	QueryTime  string `json:"queryTime"`
}

type RefMetadata struct {
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

type Result struct {
	SH  SymbolHistory
	Err error
}

type DeleteResult struct {
	SC  SymbolCorrection
	Err error
}

type MDC struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
