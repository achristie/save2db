package token

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const tokenEndpoint = "https://api.platts.com/auth/api"

var (
	cache          Token
	ErrInvalidCred = errors.New("token: invalid credentials or API key")
	ErrServerIssue = errors.New("token: unable to reach the server")
	ErrRateLimited = errors.New("token: rate limit exceeded")
)

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Iat          time.Time `json:"-"`
}

func (t *Token) String() string {
	return t.AccessToken
}

type TokenClient struct {
	TokenEndpoint string
	username      string
	password      string
	apikey        string
}

func NewTokenClient(username, password, apikey string) *TokenClient {
	return &TokenClient{
		TokenEndpoint: tokenEndpoint,
		username:      username,
		password:      password,
		apikey:        apikey,
	}
}

func (tc *TokenClient) GetToken() (*Token, error) {
	if time.Now().Before(cache.Iat.Add(50 * time.Minute)) {
		return &cache, nil
	}

	data := url.Values{}
	data.Set("username", tc.username)
	data.Set("password", tc.password)

	req, err := http.NewRequest("POST", tc.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("appkey", tc.apikey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		switch s := res.StatusCode; s {
		case http.StatusForbidden, http.StatusUnauthorized, http.StatusBadRequest:
			return nil, ErrInvalidCred
		case http.StatusTooManyRequests:
			return nil, ErrRateLimited
		default:
			return nil, ErrServerIssue
		}
	}

	var token Token
	j := json.NewDecoder(res.Body)
	if err = j.Decode(&token); err != nil {
		return nil, err
	}

	token.Iat = time.Now()
	cache = token

	return &token, nil
}
