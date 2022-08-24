package platts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
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

// Page through the API response to get all results.
// If the first page fails then abort, otherwise attempt to get all pages.
// Errors must be handled by consumer
func getConcurrently[T Concurrentable](c *Client, req *http.Request, ch chan Result[T], result T) {
	if _, err := c.do(req, &result); err != nil {
		// If first request fails then abort
		log.Fatalf("Could not make first request: %s", err)
	}
	ch <- Result[T]{result, nil}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 1) // semaphore to avoid throttling

	// create a request per page from 2 .. n
	for i := 2; i <= result.GetTotalPages(); i++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}
			var result T

			// copy the request but change the page
			p := req.URL.Path
			q := req.URL.Query()
			q.Set("page", strconv.Itoa(page))

			// generate a request
			req, err := c.newRequest(p, q)
			if err != nil {
				ch <- Result[T]{result, err}
			}

			// make the request
			_, err = c.do(req, &result)
			if err != nil {
				ch <- Result[T]{result, err}
			}
			ch <- Result[T]{result, nil} // send marshalled response to the channel
			<-sem
		}(i)
	}

	wg.Wait()
	close(ch)
}
