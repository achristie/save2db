package plattsapi

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL  string
	apiKey   string
	username string
	password string
	c        *http.Client
}

func New(apiKey *string, username *string, password *string) *Client {
	return &Client{
		apiKey:   *apiKey,
		baseURL:  "https://api.platts.com/",
		c:        &http.Client{Timeout: time.Minute},
		username: *username,
		password: *password,
	}
}

func (c *Client) CallApi() SymbolHistory {
	u := c.baseURL + "market-data/v3/value/history/mdc?"
	params := url.Values{}
	params.Add("filter", "mdc IN (\"IF\") AND modDate >= \"2022-5-27\"")
	params.Add("sort", "modDate: asc")
	params.Add("pagesize", "5")
	// params.Add("field", "deltaPrice, pValue, deltaPercent, pDate")
	req, err := http.NewRequest(http.MethodGet, u+params.Encode(), nil)
	if err != nil {
		log.Print(err, "Could not make HTTP Request")
	}
	token := GetToken(c.username, c.password, c.apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("appkey", c.apiKey)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	// requestDump, _ := httputil.DumpRequest(req, true)
	// log.Println(string(requestDump))

	res, err := c.c.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	var data SymbolHistory

	switch res.StatusCode {
	case 200:
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			log.Print(err)
		}
		// log.Printf("%+v", data)
	}

	return data

}
