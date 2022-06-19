package platts

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"
)

func (c *Client) GetSubscribedMDC() ([]MDCCount, error) {
	params := url.Values{}
	params.Add("subscribed_only", "true")
	params.Add("Facet.Field", "mdc")
	params.Add("Field", "symbol")
	params.Add("PageSize", "1")

	req, err := c.newRequest("market-data/reference-data/v3/search", params)

	if err != nil {
		return []MDCCount{}, err
	}
	var result ReferenceData

	if _, err = c.do(req, &result); err != nil {
		return []MDCCount{}, err
	}

	var s []MDCCount
	for k, v := range result.Facets.FacetCounts.Mdc {
		count, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("platts: Could not convert count to int for MDC: [%s], %s", k, err)
		}
		s = append(s, MDCCount{MDC: k, SymbolCount: count})
	}
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].SymbolCount < s[j].SymbolCount
	})
	return s, nil

}

// Only need to get DEL because UPD is handled by modified date endpoint
func (c *Client) GetDeletes(StartTime time.Time, Page int, PageSize int) (SymbolCorrection, error) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("correctionType:\"DEL\" AND modDate >= \"%s\"", StartTime.Format("2006-01-02T15:04:05")))
	params.Add("sort", "modDate: asc") // important for paging properly
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
	params.Add("sort", "modDate: asc") // important for paging properly
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

// Concurrently get history.
// SymbolHistory or Error is sent to channel.
// Consumer must decide how to handle err.
func (c *Client) GetHistoryByMDCConcurrent(Mdc string, StartTime time.Time, PageSize int, ch chan Result) error {
	// get first page
	sh, err := c.GetHistoryByMDC(Mdc, StartTime, 1, PageSize)
	if err != nil {
		return err
	}
	ch <- Result{sh, nil}

	var wg sync.WaitGroup
	// semaphore to prevent throttling
	sem := make(chan struct{}, 2)

	// loop through remaining pages and fetch concurrently
	for i := 2; i <= sh.Metadata.TotalPages; i++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}

			sh, err := c.GetHistoryByMDC(Mdc, StartTime, page, PageSize)
			if err != nil {
				ch <- Result{SymbolHistory{}, err}
			} else {
				ch <- Result{sh, nil}
			}

			<-sem
		}(i)
	}
	wg.Wait()
	close(ch)
	return nil
}

func (c *Client) GetRefData(Page int, PageSize int) (ReferenceData, error) {
	params := url.Values{}
	params.Add("subscribed_only", "true")
	params.Add("pagesize", strconv.Itoa(PageSize))
	params.Add("page", strconv.Itoa(Page))

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
