package platts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/achristie/save2db/pkg/platts/token"
)

const (
	baseURL = "https://api.platts.com/"
)

type Client struct {
	baseURL  string
	apiKey   string
	username string
	password string
	c        *http.Client
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewClient(apiKey string, username string, password string) *Client {
	return &Client{
		apiKey:   apiKey,
		baseURL:  baseURL,
		c:        &http.Client{Timeout: time.Minute},
		username: username,
		password: password,
	}
}

func (c *Client) newRequest(path string, query url.Values) (*http.Request, error) {
	url := &c.baseURL
	req, err := http.NewRequest(http.MethodGet, *url+path+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	tc := token.NewTokenClient(c.username, c.password, c.apiKey)
	token, err := tc.GetToken()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("appkey", c.apiKey)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}

func (c *Client) do(req *http.Request, target interface{}) (*http.Response, error) {
	// req.Close = true
	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	u, _ := url.QueryUnescape(req.URL.String())
	// c.infoLog.Printf("platts: [%d] %s", res.StatusCode, u)

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		// c.errorLog.Printf("platts: [%d] %s", res.StatusCode, body)
		return nil, fmt.Errorf("platts: [%s] %s %s\n %s", req.Method, res.Status, u, body)
	}

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Page through the API response to get all results.
// If the first page fails then abort, otherwise attempt to get all pages.
// Errors must be handled by consumer
func getConcurrently[T Concurrentable](c *Client, req *http.Request, ch chan Result[T], result T) {
	// make an initial request
	if _, err := c.do(req, &result); err != nil {
		ch <- Result[T]{nil, err}
		return
	}
	ch <- Result[T]{&result, nil}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 3) // semaphore to avoid throttling

	// create a request per page from 2 .. n
	for i := 2; i <= result.GetTotalPages(); i++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()
			sem <- struct{}{}
			var result T

			// copy the request but change the page
			p := req.URL.Path[1:len(req.URL.Path)]
			q := req.URL.Query()
			q.Set("page", strconv.Itoa(page))

			// generate a new request
			req, err := c.newRequest(p, q)
			if err != nil {
				ch <- Result[T]{nil, err}
			}

			// make the request
			_, err = c.do(req, &result)
			if err != nil {
				ch <- Result[T]{nil, err}
			}
			ch <- Result[T]{&result, nil} // send marshalled response to the channel
			<-sem
		}(i)
	}

	wg.Wait()
	close(ch)
}
