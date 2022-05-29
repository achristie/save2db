package main

import (
	"flag"

	"github.com/achristie/save2db/pkg/plattsapi"
)

func main() {
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")

	flag.Parse()
	c := plattsapi.New(APIKey, Username, Password)
	c.CallApi()

	// t := GetToken(*Username, *Password, *APIKey)

	// log.Print(t.AccessToken)
	// log.Println(t.RefreshToken)
}
