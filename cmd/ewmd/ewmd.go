package main

import (
	"flag"
	"log"

	"github.com/achristie/save2db/pkg/platts"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	// StartDate := flag.String("t", time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02T15:04:05"), "Get updates since date. Format 2006-01-02T15:04:05")
	// PageSize := flag.Int("p", 1000, "The page size to use for API Calls. Max is 10,000")
	flag.Parse()

	client := platts.NewClient(APIKey, Username, Password)

	td, err := client.GetEWindowMarketData()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v, %+v", td.Results[0].UpdateTime, td.Results[0].OrderTime)

}
