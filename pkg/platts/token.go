package platts

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const TokenEndpoint = "https://api.platts.com/auth/api"

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

func GetToken(Username string, Password string, APIKey string, errorLog, infoLog *log.Logger) (*Token, error) {
	if time.Now().Before(cache.Iat.Add(50 * time.Minute)) {
		return &cache, nil
	}
	data := url.Values{}
	data.Set("username", Username)
	data.Set("password", Password)

	req, err := http.NewRequest("POST", TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("appkey", APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		errorLog.Printf("[%d] %s", res.StatusCode, body)
		switch s := res.StatusCode; s {
		case http.StatusForbidden, http.StatusUnauthorized, http.StatusBadRequest:
			return nil, ErrInvalidCred
		case http.StatusTooManyRequests:
			return nil, ErrRateLimited
		default:
			return nil, ErrServerIssue
		}
	}

	j := json.NewDecoder(res.Body)
	var token Token
	if err = j.Decode(&token); err != nil {
		return nil, err
	}
	token.Iat = time.Now()
	cache = token

	return &token, nil
}
