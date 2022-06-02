package main

import (
	"flag"
	"log"
	"time"

	"github.com/achristie/save2db"
	platts "github.com/achristie/save2db/pkg/platts"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	flag.Parse()

	// A Platts API Client
	client := platts.NewClient(APIKey, Username, Password)

	// rd, err := client.GetSubscribedMDC()
	// if err != nil {
	// 	log.Println(err)
	// }
	// for k, v := range rd.Facets.FacetCounts.Mdc {
	// 	log.Print(k, v)
	// }

	// Initialize DB and create market_data table if it does not exist
	MarketDataStore := save2db.InitializeDb("database9.db")

	// Initial parameters
	page := 1
	log.Println(MarketDataStore.GetLatestModifiedDate())
	start := time.Now().AddDate(0, 0, -30)
	pageSize := 20

	// Page through response
	for {
		// Call History API
		sh, err := client.GetHistoryByMDC("AT", start, page, pageSize)
		if err != nil {
			log.Fatal(err)
		}

		// Add Response to database
		MarketDataStore.AddPricingData(sh)

		// Exit loop when all records have been fetched
		if sh.Metadata.TotalPages == page {
			break
		}

		// Avoid getting throttled by the API
		time.Sleep(2 * time.Second)
		page += 1
	}
}
