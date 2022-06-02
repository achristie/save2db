package plattsapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func GetToken(Username string, Password string, APIKey string) Token {
	if time.Now().Before(cache.Iat.Add(45 * time.Minute)) {
		return cache
	}
	data := url.Values{}
	data.Set("username", Username)
	data.Set("password", Password)

	req, err := http.NewRequest("POST", TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("appkey", APIKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		log.Fatalf("Not able to fetch a token. Please check your credentials: %s", body)
	}

	j := json.NewDecoder(res.Body)
	var token Token
	if err = j.Decode(&token); err != nil {
		log.Fatal(err)
	}
	token.Iat = time.Now()
	cache = token
	return token
}
