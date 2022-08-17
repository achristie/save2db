package platts

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// Call Corrections endpoint to get Deletes and Backfills.
// Updates are handled by the History endpoint.
func (c *Client) GetDeletes(StartTime time.Time, Page int, PageSize int) (SymbolCorrection, error) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("correctionType:\"DEL\" AND modDate >= %q", StartTime.Format("2006-01-02T15:04:05")))
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

// Concurrent version of GetDeletes.
// Automatically pages through all results.
// Correction or Error is sent to channel.
// Errors must be handled by consumer.
func (c *Client) GetDeletesConcurrent(StartTime time.Time, PageSize int, ch chan Result[SymbolCorrection]) error {
	sc, err := c.GetDeletes(StartTime, 1, PageSize)
	if err != nil {
		return err
	}
	ch <- Result[SymbolCorrection]{sc, nil}

	sem := make(chan struct{}, 1)
	var wg sync.WaitGroup
	for i := 2; i <= sc.Metadata.TotalPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}
			sc, err := c.GetDeletes(StartTime, page, PageSize)
			if err != nil {
				ch <- Result[SymbolCorrection]{SymbolCorrection{}, err}
			} else {
				ch <- Result[SymbolCorrection]{sc, nil}
			}
			<-sem
		}(i)
	}

	wg.Wait()
	close(ch)
	return nil
}

// Call History endpoint to get historical assessments.
func (c *Client) GetHistoryByMDC(Mdc string, StartTime time.Time, Page int, PageSize int) (SymbolHistory, error) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("mdc IN (%q) AND modDate >= %q", Mdc, StartTime.Format("2006-01-02T15:04:05")))
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
func (c *Client) GetHistoryByMDCConcurrent(Mdc string, StartTime time.Time, PageSize int, ch chan Result[SymbolHistory]) {
	// get first page
	sh, err := c.GetHistoryByMDC(Mdc, StartTime, 1, PageSize)
	if err != nil {
		ch <- Result[SymbolHistory]{SymbolHistory{}, err}
	}
	ch <- Result[SymbolHistory]{sh, nil}

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
				ch <- Result[SymbolHistory]{SymbolHistory{}, err}
			} else {
				ch <- Result[SymbolHistory]{sh, nil}
			}

			<-sem
		}(i)
	}
	wg.Wait()
	close(ch)
}

// Call Search endpoint to get reference data.
func (c *Client) GetRefData(StartTime time.Time, PageSize int, ch chan interface{}) {
	params := url.Values{}
	params.Add("subscribed_only", "true")
	params.Add("pagesize", strconv.Itoa(PageSize))
	params.Add("page", "1")
	params.Add("q", "brent")

	req, err := c.newRequest("market-data/reference-data/v3/search", params)
	if err != nil {
		log.Println(err)
	}

	c.GetConcurrently(req, ch)
	log.Println("here")

}

func (c *Client) GetConcurrently(req *http.Request, ch chan interface{}) {

	// make first request!
	var result ReferenceData
	if _, err := c.do(req, &result); err != nil {
		log.Println(err)
	}
	ch <- result

	var wg sync.WaitGroup
	sem := make(chan struct{}, 2)

	for i := 2; i <= result.Metadata.TotalPages; i++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}

			p := req.URL.Path
			q := req.URL.Query()
			q.Set("page", strconv.Itoa(page))

			req, err := c.newRequest(p, q)
			if err != nil {
				log.Println(err)
			}

			var result ReferenceData
			_, err = c.do(req, &result)

			if err != nil {
				log.Println(err)
			} else {
				ch <- result
			}
			<-sem
		}(i)
	}

	wg.Wait()
	close(ch)
}

// rd, err := c.GetRefData(StartTime, 1, PageSize)
// if err != nil {
// 	ch <- Result[ReferenceData]{ReferenceData{}, err}
// }
// ch <- Result[ReferenceData]{rd, nil}

// var wg sync.WaitGroup
// sem := make(chan struct{}, 2)

// for i := 2; i <= rd.Metadata.TotalPages; i++ {
// 	wg.Add(1)

// 	go func(page int) {
// 		defer wg.Done()
// 		sem <- struct{}{}

// 		rd, err := c.GetRefData(StartTime, page, PageSize)
// 		if err != nil {
// 			ch <- Result[ReferenceData]{ReferenceData{}, err}

// 		} else {
// 			ch <- Result[ReferenceData]{rd, nil}
// 		}
// 		<-sem
// 	}(i)
// }
// wg.Wait()
// close(ch)
