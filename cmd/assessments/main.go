package main

import (
	"flag"
	"log"
	"sync"
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
	MDC := flag.String("mdc", "ET", "The MDC (Market Data Category) to fetch data for")
	flag.Parse()

	// create a platts api client
	client := platts.NewClient(APIKey, Username, Password)
	t := platts.GetToken(*Username, *Password, *APIKey)
	log.Print(t)

	// initialize DB and create market_data table if it does not exist
	db := save2db.InitializeDb("database.db")

	// initial parameters
	start, err := time.Parse("2006-01-02T15:04:05", *StartDate)
	if err != nil {
		log.Fatal(err)
	}

	// Update market_data table with records modified since `start`
	// GetAssessments(client, db, *MDC, start, *PageSize)

	var wg sync.WaitGroup
	guard := make(chan struct{}, 2)
	for i := 1; i <= 9; i++ {
		wg.Add(1)

		go func(page int) {
			defer wg.Done()
			guard <- struct{}{}

			sh, err := client.GetHistoryByMDC(*MDC, start, page, *PageSize)

			if err != nil {
				log.Print(err)
			}
			log.Printf("Page [%d] - Fetched up to [%d] of [%d] records in [%s] and added to DB", page, sh.Metadata.PageSize, sh.Metadata.Count, sh.Metadata.QueryTime)
			if err := db.Add(sh); err != nil {
				log.Print("error inserting records: ", err)
			}
			<-guard

		}(i)

	}
	wg.Wait()
}

// Uses the `client` to fetch historical data for given MDC modified since `start`
// Automatically pages through all results
// and stores data into `db`
func GetAssessments(client *platts.Client, db *save2db.MarketDataStore, MDC string, start time.Time, pageSize int) {
	page := 1
	// loop until everything is fetched
	for {
		// call history endpoint
		sh, err := client.GetHistoryByMDC(MDC, start, page, pageSize)

		// if there is an error then log it and break
		if err != nil {
			log.Fatalf("error getting history: %s", err)
			break
		}

		// add data to database
		if err := db.Add(sh); err != nil {
			log.Print("error inserting records: ", err)
		}

		log.Printf("Page [%d] - Fetched up to [%d] of [%d] records in [%s] and added to DB", page, sh.Metadata.PageSize, sh.Metadata.Count, sh.Metadata.QueryTime)

		// exit loop when all data has been fetched
		if sh.Metadata.TotalPages == page || sh.Metadata.TotalPages == 0 {
			break
		}

		// avoid getting throttled by the API
		time.Sleep(1 * time.Second)
		page += 1
	}
}
