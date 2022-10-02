package platts

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (c *Client) GetTradeData(markets []string, StartTime time.Time, PageSize int, ch chan Result[TradeData]) {
	params := url.Values{}
	if len(markets) > 0 {
		params.Add("filter", fmt.Sprintf("update_time >= %q AND market IN (%s)", StartTime.Format("2006-01-02T15:04:05"), "\""+strings.Join(markets, "\",\"")+"\""))
	} else {
		params.Add("filter", fmt.Sprintf("update_time >= %q", StartTime.Format("2006-01-02T15:04:05")))
	}
	params.Add("pagesize", strconv.Itoa(min(1000, PageSize))) // max is 1k
	params.Add("sort", "update_time: asc")

	req, err := c.newRequest("tradedata/v3/ewindowdata", params)
	if err != nil {
		ch <- Result[TradeData]{TradeData{}, err}
	}

	go func() {
		var result TradeData
		getConcurrently(c, req, ch, result)
	}()

}
