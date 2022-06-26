package platts

import (
	"net/url"
)

func (c *Client) GetEWindowMarketData() (EWindowMarketData, error) {
	params := url.Values{}
	// params.Add("subscribed_only", "true")
	// params.Add("Facet.Field", "mdc")
	// params.Add("Field", "symbol")
	params.Add("PageSize", "2")

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
