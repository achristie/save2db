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
	db := save2db.InitializeDb("database.db")

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal(err)
	}
	MDCs, err := client.GetSubscribedMDC()
	if err != nil {
		log.Fatalf("Could not get list of MDCs: %s", err)
	}

	// Update market_data table with records modified since `start`
	UpdateHistory(client, db, MDCs, start, *PageSize)

	// Update market_data table with records marked for deletion since `start`
	UpdateCorrections(client, db, start)

}

// Uses the `client` to fetch historical data for each MDC modified since `start`
// Automatically pages through all results
// and stores data into `db`
func UpdateHistory(client *platts.Client, db *save2db.MarketDataStore, MDCs []string, start time.Time, pageSize int) {
	// loop through every MDC
	page := 1
	for _, v := range MDCs {
		// loop until everything is fetched
		for {
			// call history endpoint
			sh, err := client.GetHistoryByMDC(v, start, page, pageSize)

			// if there is an error then log it and go to the next MDC
			if err != nil {
				log.Printf("error getting history: %s", err)
				break
			}

			// add data to database
			if err := db.Add(sh); err != nil {
				log.Print("error inserting records: ", err)
			}

			log.Printf("Fetched up to [%d] of [%d] records in [%s] and added to DB", sh.Metadata.PageSize, sh.Metadata.Count, sh.Metadata.QueryTime)

			// exit loop when all data has been fetched
			if sh.Metadata.TotalPages == page || sh.Metadata.TotalPages == 0 {
				break
			}

			// avoid getting throttled by the API
			time.Sleep(2 * time.Second)
			page += 1
		}
		// reset page counter, sleep for throttling
		page = 1
		time.Sleep(2 * time.Second)
	}
}

func UpdateCorrections(client *platts.Client, db *save2db.MarketDataStore, start time.Time) {
	sc, err := client.GetDeletes(start, 1, 10000)
	if err != nil {
		log.Printf("error getting corrections: %s", err)
	}
	if err := db.Remove(sc); err != nil {
		log.Printf("error deleting records: %s", err)
	}
	log.Printf("Fetched [%d] records in [%s] and removed from DB", sc.Metadata.PageSize, sc.Metadata.QueryTime)
}
