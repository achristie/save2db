package plattsapi

import (
	"encoding/json"
	"net/url"
)

type Facets struct {
	FacetCounts *json.RawMessage `json:"facet_counts"`
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

func (c *Client) GetSubscribedMDC() (ReferenceData, error) {
	params := url.Values{}
	params.Add("subscribed_only", "true")
	params.Add("Facet.Field", "mdc")
	params.Add("PageSize", "1")

	req, err := c.newRequest("market-data/reference-data/v3/search", params)

	if err != nil {
		return ReferenceData{}, err
	}
	var result ReferenceData

	if _, err = c.do(req, &result); err != nil {
		return ReferenceData{}, err
	}

	return result, nil

}

func (c *Client) GetHistoryByMDC(mdc []string) (SymbolHistory, error) {
	params := url.Values{}
	params.Add("filter", "mdc IN (\"IF\") AND modDate >= \"2022-5-27\"")
	params.Add("sort", "modDate: asc")
	params.Add("pagesize", "5")
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
