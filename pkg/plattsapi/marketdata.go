package plattsapi

type SymbolHistory struct {
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
			Value       float32 `json:"value"`
			AssessDate  string  `json:"assessDate"`
			IsCorrected string  `json:"isCorrected"`
			ModDate     string  `json:"modDate"`
		} `json:"data"`
	} `json:"results"`
}
