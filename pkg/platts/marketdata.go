package plattsapi

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

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

func (c *Client) GetSubscribedMDC() ([]string, error) {
	params := url.Values{}
	params.Add("subscribed_only", "true")
	params.Add("Facet.Field", "mdc")
	params.Add("Field", "symbol")
	params.Add("PageSize", "1")

	req, err := c.newRequest("market-data/reference-data/v3/search", params)

	if err != nil {
		return []string{}, err
	}
	var result ReferenceData

	if _, err = c.do(req, &result); err != nil {
		return []string{}, err
	}

	var s []string
	for k := range result.Facets.FacetCounts.Mdc {
		s = append(s, k)
	}
	return s, nil

}

func (c *Client) GetHistoryByMDC(Mdc string, StartTime time.Time, Page int, PageSize int) (SymbolHistory, error) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("mdc IN (\"%s\") AND modDate >= \"%s\"", Mdc, StartTime.Format("2006-01-02T15:04:05")))
	params.Add("sort", "modDate: asc")
	params.Add("pagesize", strconv.Itoa(PageSize))
	params.Add("page", strconv.Itoa(Page))
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
