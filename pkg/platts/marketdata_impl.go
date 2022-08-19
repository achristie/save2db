package platts

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Call Corrections endpoint to get Deletes
// Only deletes are necessary here because History endpoints will contain correctiosn and backfills
func (c *Client) GetDeletes(StartTime time.Time, PageSize int, ch chan Result[SymbolCorrection]) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("correctionType:\"DEL\" AND modDate >= %q", StartTime.Format("2006-01-02T15:04:05")))
	params.Add("sort", "modDate: asc")                         // important for paging properly
	params.Add("pagesize", strconv.Itoa(min(10000, PageSize))) // max is 10k
	params.Add("page", "1")

	req, err := c.newRequest("market-data/v3/value/correction/modified-date", params)
	if err != nil {
		ch <- Result[SymbolCorrection]{SymbolCorrection{}, err}
	}

	go func() {
		var result SymbolCorrection
		getConcurrently(c, req, ch, result)
	}()
}

func (c *Client) GetHistory(StartTime time.Time, PageSize int, ch chan Result[SymbolHistory]) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("modDate >= %q", StartTime.Format("2006-01-02T15:04:05")))
	params.Add("sort", "modDate: asc")                         // important for paging properly
	params.Add("pagesize", strconv.Itoa(min(10000, PageSize))) // max is 10k

	req, err := c.newRequest("market-data/v3/value/history/", params)
	if err != nil {
		ch <- Result[SymbolHistory]{SymbolHistory{}, err}
	}

	go func() {
		var result SymbolHistory
		getConcurrently(c, req, ch, result)
	}()
}

// Call HistoryByMDC endpoint to get historical price assessments for Symbols by Market Data Category.
// Platts Symbols are grouped into MDCs
func (c *Client) GetHistoryByMDC(Mdc string, StartTime time.Time, PageSize int, ch chan Result[SymbolHistory]) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("mdc IN (%q) AND modDate >= %q", Mdc, StartTime.Format("2006-01-02T15:04:05")))
	params.Add("sort", "modDate: asc")                         // important for paging properly
	params.Add("pagesize", strconv.Itoa(min(10000, PageSize))) // max is 10k

	req, err := c.newRequest("market-data/v3/value/history/mdc", params)
	if err != nil {
		ch <- Result[SymbolHistory]{SymbolHistory{}, err}
	}

	go func() {
		var result SymbolHistory
		getConcurrently(c, req, ch, result)
	}()
}

// Call Search endpoint to get Reference Data for Symbols
// Example fields include Description, Commodity, Geography, etc..
func (c *Client) GetReferenceData(StartTime time.Time, PageSize int, ch chan Result[ReferenceData]) {
	params := url.Values{}
	params.Add("subscribed_only", "true")
	params.Add("pagesize", strconv.Itoa(min(1000, PageSize))) // max is 1k
	params.Add("q", "brent")

	req, err := c.newRequest("market-data/reference-data/v3/search", params)
	if err != nil {
		ch <- Result[ReferenceData]{ReferenceData{}, err}
	}

	go func() {
		var result ReferenceData
		getConcurrently(c, req, ch, result)
	}()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
