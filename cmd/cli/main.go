package main

import (
	"flag"

	"github.com/achristie/save2db"
	"github.com/achristie/save2db/pkg/plattsapi"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")

	MarketDataStore := save2db.InitializeDb("database.db")

	flag.Parse()
	client := plattsapi.NewClient(APIKey, Username, Password)
	data := client.CallApi()

	MarketDataStore.AddPricingData(data)

}
