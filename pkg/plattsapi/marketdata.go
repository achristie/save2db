package plattsapi

type CurrentSymbol struct {
	Metadata struct {
		Count      int    `json:"count"`
		PageSize   int    `json:"pageSize"`
		Page       int    `json:"page"`
		TotalPages int    `json:"totalPages"`
		QueryTime  string `json:"queryTime"`
	} `json:"metadata"`
	Results []struct {
		Symbol string `json:"symbol"`
		Data   []struct {
			Bate        string  `json:"bate"`
			AssessDate  string  `json:"assessDate"`
			Value       float64 `json:"value"`
			IsCorrected string  `json:"isCorrected"`
			ModDate     string  `json:"modDate"`
			Change      struct {
				PValue       float64 `json:"pValue"`
				PDate        string  `json:"pDate"`
				DeltaPrice   float64 `json:"deltaPrice"`
				DeltaPercent float64 `json:"deltaPercent"`
			} `json:"change"`
		} `json:"data"`
	} `json:"results"`
}
