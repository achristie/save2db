package plattsapi

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

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

// Only need to get DEL because UPD is handled by correction endpoint
func (c *Client) GetDeletes(StartTime time.Time, Page int, PageSize int) (SymbolCorrection, error) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("correctionType:\"DEL\" AND modDate >= \"%s\"", StartTime.Format("2006-01-02T15:04:05")))
	params.Add("sort", "modDate: asc")
	params.Add("pagesize", strconv.Itoa(PageSize))
	params.Add("page", strconv.Itoa(Page))

	req, err := c.newRequest("market-data/v3/value/correction/modified-date", params)
	if err != nil {
		return SymbolCorrection{}, err
	}

	var result SymbolCorrection
	if _, err = c.do(req, &result); err != nil {
		return SymbolCorrection{}, err
	}

	return result, nil

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
