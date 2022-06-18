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
	MDC := flag.String("mdc", "ET", "The MDC (Market Data Category) to fetch data for")
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

	// Update market_data table with records modified since `start`
	GetAssessments(client, db, *MDC, start, *PageSize)

}

// Uses the `client` to fetch historical data for given MDC modified since `start`
// Store results in DB
func GetAssessments(client *platts.Client, db *save2db.MarketDataStore, MDC string, start time.Time, pageSize int) {
	ch := make(chan platts.SymbolHistory)

	go func() {
		err := client.ConcurrentGetHistoryByMDC(MDC, start, pageSize, ch)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for sh := range ch {
		log.Printf("[%d] records received from page [%d]. Adding to DB", len(sh.Results), sh.Metadata.Page)
		if err := db.Add(sh); err != nil {
			log.Printf("Error inserting records: %s", err)
		}
	}
}
