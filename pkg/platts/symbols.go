package platts

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type SymbolResponse struct {
	Metadata SymbolResponseMetadata
	Results  []SymbolResponseResults
}

type SymbolResponseMetadata struct {
	Count      int    `json:"count"`
	PageSize   int    `json:"page_size"`
	Page       int    `json:"page"`
	TotalPages int    `json:"total_pages"`
	QueryTime  string `json:"query_time"`
}

type SymbolResponseResults struct {
	Symbol                   string   `json:"symbol"`
	Description              string   `json:"description"`
	Commodity                string   `json:"commodity"`
	UOM                      string   `json:"uom"`
	Active                   string   `json:"active"`
	DeliveryRegion           string   `json:"delivery_region"`
	DeliveryRegionBasis      string   `json:"delivery_region_basis"`
	ContractType             string   `json:"contract_type"`
	PublicationFrequencyCode string   `json:"publication_frequency_code"`
	ShippingTerms            string   `json:"shipping_terms"`
	DayOfPublication         string   `json:"day_of_publication"`
	StandardLotSize          float32  `json:"standard_lot_size"`
	StandardLotUnits         string   `json:"standard_lot_units"`
	QuotationStyle           string   `json:"quotation_style"`
	Bate                     []string `json:"bate_code"`
	CommodityGrade           string   `json:"commodity_grade"`
	Currency                 string   `json:"currency"`
	AssessmentFrequency      string   `json:"assessment_frequency"`
	Timestamp                string   `json:"timestamp"`
	SettlementType           string   `json:"settlement_type"`
	DecimalPlaces            int      `json:"decimal_places"`
	MDCNames                 []string `json:"mdc"`
	MDCDescriptions          []string `json:"mdc_description"`
}

type PagingDetails struct {
	Url        *url.URL
	TotalPages int
	PageKey    string
}

type PageOptions struct {
	PageSize int
	Page     int
}

// func (s SymbolResponse) GetTotalPages() int {
// 	return s.Metadata.TotalPages
// }

func (c *Client) GetSymbols(q string, lastModified time.Time, po PageOptions) (SymbolResponse, PagingDetails, error) {
	params := url.Values{}
	params.Add("subscribed_only", "true")
	params.Add("pagesize", strconv.Itoa(1000)) // max is 1k
	params.Add("q", fmt.Sprintf("%q", q))

	req, err := c.newRequest("market-data/reference-data/v3/search", params)
	if err != nil {
		return SymbolResponse{}, PagingDetails{}, err
	}

	var sr SymbolResponse
	_, err = c.do(req, &sr)
	if err != nil {
		return SymbolResponse{}, PagingDetails{}, err
	}

	p := PagingDetails{Url: req.URL, TotalPages: sr.Metadata.TotalPages, PageKey: "page"}

	return sr, p, nil
}

// FetchAll modifies the `p.Url` to fetch results for pages 2 .. `p.TotalPages` with a concurrency of `limit`
// Results and Errors are returned as they arrive on the returned channel
// Channel is closed after all results are fetched
func FetchAll[T any](c *Client, p PagingDetails, limit int) <-chan Result[T] {
	out := make(chan Result[T])
	go func() {
		sem := make(chan struct{}, limit) // semaphore
		var wg sync.WaitGroup

		for i := 2; i <= p.TotalPages; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				sem <- struct{}{}

				q := p.Url.Query()
				q.Set(p.PageKey, strconv.Itoa(i))

				req, err := c.newRequest(p.Url.Path, q)
				if err != nil {
					out <- Result[T]{Message: nil, Error: err}
				}

				var target T
				_, err = c.do(req, &target)
				if err != nil {
					out <- Result[T]{Message: nil, Error: err}
				}

				out <- Result[T]{Message: &target, Error: nil}

				<-sem
			}(i)

		}
		wg.Wait()
		close(out)
	}()

	return out
}
