package plattsapi

import (
	"net/url"
)

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

func (c *Client) GetHistoryByMDC(mdc []string) (SymbolHistory, error) {
	params := url.Values{}
	params.Add("filter", "mdc IN (\"IF\") AND modDate >= \"2022-5-27\"")
	params.Add("sort", "modDate: asc")
	params.Add("pagesize", "100")
	req, err := c.newRequest("market-data/v3/value/history/mdc", params)

	if err != nil {
		return SymbolHistory{}, err
	}

	var result SymbolHistory
	if _, err = c.do(req, &result); err != nil {
		return SymbolHistory{}, err
	}

	return result, nil

}
