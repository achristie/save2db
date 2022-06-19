package main

import (
	"flag"
	"log"
	"time"

	MD "github.com/achristie/save2db/internal/market_data"
	platts "github.com/achristie/save2db/pkg/platts"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	StartDate := flag.String("t", time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02T15:04:05"), "Get updates since date. Format 2006-01-02T15:04:05")
	PageSize := flag.Int("p", 1000, "The page size to use for API Calls. Max is 10,000")
	flag.Parse()

	// create a platts api client
	client := platts.NewClient(APIKey, Username, Password)

	// initialize DB and create market_data table if it does not exist
	db := MD.InitializeDb("database.db")

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal(err)
	}

	// Update market_data table with records marked for deletion since `start`
	GetCorrections(client, db, start, *PageSize)

}

func GetCorrections(client *platts.Client, db *MD.MarketDataStore, start time.Time, pageSize int) {
	ch := make(chan platts.DeleteResult)

	go func() {
		log.Printf("Fetching corrections since %s", start.String())
		err := client.GetDeletesConcurrent(start, pageSize, ch)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for result := range ch {
		if result.Err != nil {
			log.Printf("Error retrieving data: %s", result.Err)
		} else {
			log.Printf("%d records received from page [%d] in [%s]. Removing from DB", len(result.SC.Results), result.SC.Metadata.Page, result.SC.Metadata.QueryTime)
			if err := db.Remove(result.SC); err != nil {
				log.Printf("Error removing records: %s", err)
			}
		}
	}
}
