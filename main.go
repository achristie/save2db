package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const TokenEndpoint = "https://api.platts.com/auth/api"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetToken(Username string, Password string, APIKey string) Token {
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
	return token
}

func main() {
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")

	flag.Parse()

	t := GetToken(*Username, *Password, *APIKey)

	log.Print(t.AccessToken)
	log.Println(t.RefreshToken)
}
