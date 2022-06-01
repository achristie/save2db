package main

import (
	"flag"
	"log"

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
	s := []string{"IF"}
	data, err := client.GetHistoryByMDC(s)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", data)
	MarketDataStore.AddPricingData(data)

}
