package platts

import (
	"fmt"
	"net/url"
)

type Watchlist struct {
	Metadata struct {
		Count     int `json:"count"`
		Page      int `json:"page"`
		PageSize  int `json:"pageSize"`
		QueryTime int `json:"queryTime"`
	} `json:"metadata"`
	Results []struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		ConfigKey struct {
			Name string `json:"name"`
		} `json:"configKey"`
		Payload      []string    `json:"payload"`
		CreatedAtUtc string      `json:"createdAtUtc"`
		UpdatedAtUtc string      `json:"updatedAtUtc"`
		ParentID     interface{} `json:"parentId"`
	} `json:"results"`
}

func (c *Client) GetWatchlistByName(name string) (Watchlist, error) {
	params := url.Values{}
	params.Add("configKey", fmt.Sprintf("{'name': %q}", name))

	req, err := c.newRequest("platts-platform/user-preferences/v1/configurations/listmanagement", params)
	if err != nil {
		return Watchlist{}, err
	}

	var wl Watchlist
	if _, err = c.do(req, &wl); err != nil {
		return Watchlist{}, err
	}

	return wl, nil
}
