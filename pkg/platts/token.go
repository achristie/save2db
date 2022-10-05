package platts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const TokenEndpoint = "https://api.platts.com/auth/api"

var cache Token

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Iat          time.Time `json:"-"`
}

func GetToken(Username string, Password string, APIKey string) (*Token, error) {
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
		return nil, fmt.Errorf("unable to fetch a token; please check your credentials: %s", body)
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
