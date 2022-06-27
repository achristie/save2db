package platts

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func (c *Client) GetEWindowMarketData(StartTime time.Time, Page int, PageSize int) (EWindowMarketData, error) {
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("update_time >= %q", StartTime.Format("2006-01-02T15:04:05")))
	params.Add("pagesize", strconv.Itoa(PageSize))
	params.Add("page", strconv.Itoa(Page))

	req, err := c.newRequest("tradedata/v3/ewindowdata", params)

	if err != nil {
		return EWindowMarketData{}, err
	}

	var result EWindowMarketData
	if _, err = c.do(req, &result); err != nil {
		return EWindowMarketData{}, err
	}

	return result, nil
}
