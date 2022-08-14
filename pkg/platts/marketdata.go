package platts

type ReferenceData struct {
	Metadata RefMetadata `json:"metadata"`
	Results  []struct {
		Symbol                   string                 `json:"symbol"`
		Description              string                 `json:"description"`
		Commodity                string                 `json:"commodity"`
		UOM                      string                 `json:"uom"`
		Active                   string                 `json:"active"`
		DeliveryRegion           string                 `json:"delivery_region"`
		DeliveryRegionBasis      string                 `json:"delivery_region_basis"`
		ContractType             string                 `json:"contract_type"`
		PublicationFrequencyCode string                 `json:"publication_frequency_code"`
		ShippingTerms            string                 `json:"shipping_terms"`
		DayOfPublication         string                 `json:"day_of_publication"`
		StandardLotSize          float32                `json:"standard_lot_size"`
		StandardLotUnits         string                 `json:"standard_lot_units"`
		QuotationStyle           string                 `json:"quotation_style"`
		Bate                     []string               `json:"bate_code"`
		CommodityGrade           string                 `json:"commodity_grade"`
		Currency                 string                 `json:"currency"`
		AssessmentFrequency      string                 `json:"assessment_frequency"`
		Timestamp                string                 `json:"timestamp"`
		SettlementType           string                 `json:"settlement_type"`
		DecimalPlaces            int                    `json:"decimal_places"`
		Rest                     map[string]interface{} `json:"-"`
	} `json:"results"`
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
	name        string
	description string
}
