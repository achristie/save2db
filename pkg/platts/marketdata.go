package plattsapi

type Facets struct {
	FacetCounts struct {
		Mdc map[string]string `json:"mdc"`
	} `json:"facet_counts"`
}

type ReferenceData struct {
	Metadata Metadata `json:"metadata"`
	Facets   Facets   `json:"facets"`
}

type Metadata struct {
	Count      int    `json:"count"`
	PageSize   int    `json:"pageSize"`
	Page       int    `json:"page"`
	TotalPages int    `json:"totalPages"`
	QueryTime  string `json:"queryTime"`
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
			Value       float32 `json:"value"`
			AssessDate  string  `json:"assessDate"`
			IsCorrected string  `json:"isCorrected"`
			ModDate     string  `json:"modDate"`
		} `json:"data"`
	} `json:"results"`
}
