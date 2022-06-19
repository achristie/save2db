package main

import (
	"flag"
	"log"
	"time"

	save2db "github.com/achristie/save2db/internal"
	platts "github.com/achristie/save2db/pkg/platts"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	StartDate := flag.String("t", time.Now().UTC().AddDate(0, 0, -3).Format("2006-01-02T15:04:05"), "Get updates since date. Format 2006-01-02T15:04:05")
	flag.Parse()

	// create a platts api client
	client := platts.NewClient(APIKey, Username, Password)

	// initialize DB and create market_data table if it does not exist
	db := save2db.InitializeDb("database.db")

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal(err)
	}

	// Update market_data table with records marked for deletion since `start`
	GetCorrections(client, db, start)

}

func GetCorrections(client *platts.Client, db *save2db.MarketDataStore, start time.Time) {
	page := 1
	for {
		sc, err := client.GetDeletes(start, page, 100)
		if err != nil {
			log.Printf("error getting corrections: %s", err)
		}
		if err := db.Remove(sc); err != nil {
			log.Printf("error deleting records: %s", err)
		}
		log.Printf("Page [%d] - Fetched up to [%d] records in [%s] and removed from DB", page, sc.Metadata.PageSize, sc.Metadata.QueryTime)

		if sc.Metadata.TotalPages == page || sc.Metadata.TotalPages == 0 {
			break
		}

		// avoid getting throttled by the API
		time.Sleep(1 * time.Second)
		page += 1
	}
}
