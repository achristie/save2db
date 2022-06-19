package platts

type Facets struct {
	FacetCounts struct {
		Mdc map[string]string `json:"mdc"`
	} `json:"facet_counts"`
}

type ReferenceData struct {
	Metadata RefMetadata `json:"metadata"`
	Facets   Facets      `json:"facets"`
	Results  []struct {
		Symbol         string `json:"symbol"`
		Description    string `json:"description"`
		Commodity      string `json:"commodity"`
		UOM            string `json:"uom"`
		Active         string `json:"active"`
		DeliveryRegion string `json:"delivery_region"`
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

type MDCCount struct {
	MDC         string
	SymbolCount int
}

type Result struct {
	SH  SymbolHistory
	Err error
}

type DeleteResult struct {
	SC  SymbolCorrection
	Err error
}
