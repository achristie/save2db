package platts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func NewClient(apiKey *string, username *string, password *string) *Client {
	return &Client{
		apiKey:   *apiKey,
		baseURL:  "https://api.platts.com/",
		c:        &http.Client{Timeout: time.Minute},
		username: *username,
		password: *password,
	}
}

func (c *Client) newRequest(path string, query url.Values) (*http.Request, error) {
	url := &c.baseURL
	req, _ := http.NewRequest(http.MethodGet, *url+path+"?"+query.Encode(), nil)

	token := GetToken(c.username, c.password, c.apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("appkey", c.apiKey)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	return req, nil

}

func (c *Client) do(req *http.Request, target interface{}) (*http.Response, error) {
	req.Close = true
	u, _ := url.QueryUnescape(req.URL.String())
	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("[%s] %s %s\n %s", req.Method, res.Status, u, body)
	}

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return nil, fmt.Errorf("response error [%s] %s: %s", req.Method, u, err)
	}
	return res, nil
}
