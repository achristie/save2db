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

	client := platts.NewClient(APIKey, Username, Password)

	// rd, err := client.GetSubscribedMDC()
	// if err != nil {
	// 	log.Println(err)
	// }
	// for k, v := range rd.Facets.FacetCounts.Mdc {
	// 	log.Print(k, v)
	// }

	MarketDataStore := save2db.InitializeDb("database7.db")
	page := 1
	start := time.Now().AddDate(0, 0, -30)
	for {
		sh, err := client.GetHistoryByMDC("AT", start, page, 10)
		if err != nil {
			log.Println(err)
		}

		MarketDataStore.AddPricingData(sh)
		if sh.Metadata.TotalPages == page {
			break
		}
		time.Sleep(2 * time.Second)
		page += 1
	}
}
